/*
Package server is the gRPC implementation of the SDK gRPC server
Copyright 2018 Portworx

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
	"context"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/libopenstorage/grpc-framework/pkg/auth"
	"github.com/libopenstorage/grpc-framework/pkg/role"
	"google.golang.org/grpc"
)

// TLSConfig points to the cert files needed for HTTPS
type TLSConfig struct {
	// CertFile is the path to the cert file
	CertFile string
	// KeyFile is the path to the key file
	KeyFile string
}

// SecurityConfig provides configuration for SDK auth
type SecurityConfig struct {
	// Role implementation
	Role role.RoleManager
	// Tls configuration
	Tls *TLSConfig
	// Authenticators per issuer. You can register multple authenticators
	// based on the "iss" string in the string. For example:
	// map[string]auth.Authenticator {
	//     "https://accounts.google.com": googleOidc,
	//     "openstorage-sdk-auth: selfSigned,
	// }
	Authenticators map[string]auth.Authenticator
}

// ServerConfig provides the configuration to the SDK server
type ServerConfig struct {
	// Name of the server
	Name string
	// Net is the transport for gRPC: unix, tcp, etc.
	// For the gRPC Server. This value goes together with `Address`.
	Net string
	// Address is the port number or the unix domain socket path.
	// For the gRPC Server. This value goes together with `Net`.
	Address string
	// RestAdress is the port number. Example: 9110
	// For the gRPC REST Gateway.
	RestPort string
	// Unix domain socket for local communication. This socket
	// will be used by the REST Gateway to communicate with the gRPC server.
	// Only set for testing. Having a '%s' can be supported to use the
	// name of the driver as the driver name.
	Socket string
	// (optional) Location for audit log.
	// If not provided, it will go to /var/log/openstorage-audit.log
	AuditOutput io.Writer
	// (optional) Location of access log.
	// This is useful when authorization is not running.
	// If not provided, it will go to /var/log/grpc-framework-access.log
	AccessOutput io.Writer
	// (optional) The OpenStorage driver to use
	DriverName string
	// Security configuration
	Security *SecurityConfig
	// ServerExtensions allows you to extend the SDK gRPC server
	// with callback functions that are sequentially executed
	// at the end of Server.Start()
	//
	// To add your own service to the SDK gRPC server,
	// just append a function callback that registers it:
	//
	// s.config.ServerExtensions = append(s.config.ServerExtensions,
	// 		func(gs *grpc.Server) {
	//			api.RegisterCustomService(gs, customHandler)
	//		})
	GrpcServerExtensions []func(grpcServer *grpc.Server)

	// RestServerExtensions allows for extensions to be added
	// to the SDK Rest Gateway server.
	//
	// To add your own service to the SDK REST Server, simply add your handlers
	// to the RestSererExtensions slice. These handlers will be registered on the
	// REST Gateway http server.
	RestServerExtensions []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}
