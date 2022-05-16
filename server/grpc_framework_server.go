/*
Copyright 2022 Pure Storage

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package server

import (
	"fmt"
	"io"
	"sync"

	"github.com/libopenstorage/grpc-framework/pkg/correlation"
	grpcserver "github.com/libopenstorage/grpc-framework/pkg/grpc/server"
	"github.com/libopenstorage/grpc-framework/pkg/role"

	"github.com/sirupsen/logrus"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcFrameworkServer struct {
	*grpcserver.GrpcServer

	lock   sync.RWMutex
	name   string
	config ServerConfig

	// Loggers
	log             *logrus.Entry
	auditLogOutput  io.Writer
	accessLogOutput io.Writer

	roleServer role.RoleManager
}

// New creates a new gRPC server for the gRPC framework
func NewGrpcFrameworkServer(config *ServerConfig) (*GrpcFrameworkServer, error) {
	if nil == config {
		return nil, fmt.Errorf("Configuration must be provided")
	}

	// Default to tcp
	if len(config.Net) == 0 {
		config.Net = "tcp"
	}

	// Create a log object for this server
	name := "grpc-framework-" + config.Net
	log := logrus.WithFields(logrus.Fields{
		"name": name,
	})

	// Setup authentication
	for issuer, _ := range config.Security.Authenticators {
		log.Infof("Authentication enabled for issuer: %s", issuer)

		// Check the necessary security config options are set
		if config.Security.Role == nil {
			return nil, fmt.Errorf("Must supply role manager when authentication enabled")
		}
	}

	// Create gRPC server
	gServer, err := grpcserver.New(&grpcserver.GrpcServerConfig{
		Name:    name,
		Net:     config.Net,
		Address: config.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to setup %s server: %v", name, err)
	}

	s := &GrpcFrameworkServer{
		GrpcServer:      gServer,
		accessLogOutput: config.AccessOutput,
		auditLogOutput:  config.AuditOutput,
		roleServer:      config.Security.Role,
		config:          *config,
		name:            name,
		log:             log,
	}
	return s, nil
}

// Start is used to start the server.
// It will return an error if the server is already running.
func (s *GrpcFrameworkServer) Start() error {

	// Setup https if certs have been provided
	opts := make([]grpc.ServerOption, 0)
	if s.config.Net != "unix" && s.config.Security.Tls != nil {
		creds, err := credentials.NewServerTLSFromFile(
			s.config.Security.Tls.CertFile,
			s.config.Security.Tls.KeyFile)
		if err != nil {
			return fmt.Errorf("Failed to create credentials from cert files: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
		s.log.Info("TLS enabled")
	} else {
		s.log.Info("TLS disabled")
	}

	// Add correlation interceptor
	correlationInterceptor := correlation.ContextInterceptor{
		Origin: correlation.ComponentSDK,
	}

	// Setup authentication and authorization using interceptors if auth is enabled
	if len(s.config.Security.Authenticators) != 0 {
		opts = append(opts, grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				s.rwlockUnaryIntercepter,
				correlationInterceptor.ContextUnaryServerInterceptor,
				grpc_auth.UnaryServerInterceptor(s.auth),
				s.authorizationServerUnaryInterceptor,
				s.loggerServerUnaryInterceptor,
				grpc_prometheus.UnaryServerInterceptor,
			)))
		opts = append(opts, grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				s.rwlockStreamIntercepter,
				grpc_auth.StreamServerInterceptor(s.auth),
				s.authorizationServerStreamInterceptor,
				s.loggerServerStreamInterceptor,
				grpc_prometheus.StreamServerInterceptor,
			)))
	} else {
		opts = append(opts, grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				s.rwlockUnaryIntercepter,
				correlationInterceptor.ContextUnaryServerInterceptor,
				s.loggerServerUnaryInterceptor,
				grpc_prometheus.UnaryServerInterceptor,
			)))
		opts = append(opts, grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				s.rwlockStreamIntercepter,
				s.loggerServerStreamInterceptor,
				grpc_prometheus.StreamServerInterceptor,
			)))
	}

	// Start the gRPC Server
	err := s.GrpcServer.StartWithServer(func() *grpc.Server {
		grpcServer := grpc.NewServer(opts...)

		// Register gRPC Handlers
		for _, ext := range s.config.GrpcServerExtensions {
			ext(grpcServer)
		}

		// Register stats for all the services
		s.registerPrometheusMetrics(grpcServer)

		return grpcServer
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *GrpcFrameworkServer) registerPrometheusMetrics(grpcServer *grpc.Server) {
	// Register the gRPCs and enable latency historgram
	grpc_prometheus.Register(grpcServer)
	grpc_prometheus.EnableHandlingTimeHistogram()

	// Initialize the metrics
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	grpcMetrics.InitializeMetrics(grpcServer)
}
