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
package main

import (
	"fmt"
	"os"

	"github.com/libopenstorage/grpc-framework/pkg/util"
	"github.com/libopenstorage/grpc-framework/server"
	"github.com/libopenstorage/grpc-framework/test/app/api"
	helloserver "github.com/libopenstorage/grpc-framework/test/app/pkg/server"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	helloSocket = "/tmp/hello-server.sock"
)

func main() {
	hello := &helloserver.HelloGreeter{}
	config := &server.ServerConfig{
		Name:     "hello",
		Address:  "127.0.0.1:9009",
		RestPort: "9010",
		Socket:   helloSocket,
	}

	config.RegisterGrpcServers(func(gs *grpc.Server) {
		api.RegisterHelloGreeterServer(gs, hello)
		api.RegisterHelloIdentityServer(gs, hello)
	})

	config.RegisterRestHandlers(
		api.RegisterHelloGreeterHandler,
		api.RegisterHelloIdentityHandler,
	)

	// Create grpc framework server
	os.Remove(helloSocket)
	s, err := server.New(config)
	if err != nil {
		fmt.Printf("Unable to create server: %v", err)
		os.Exit(1)
	}

	// Setup a signal handler
	signal_handler := util.NewSigIntManager(func() {
		s.Stop()
		os.Remove(helloSocket)
		os.Exit(0)
	})
	signal_handler.Start()

	// Start server
	err = s.Start()
	if err != nil {
		fmt.Printf("Unable to start server: %v", err)
		os.Exit(1)
	}

	// Wait. The signal handler will exit cleanly
	logrus.Info("Hello server running")
	select {}
}
