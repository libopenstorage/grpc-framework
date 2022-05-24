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
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	grpcclient "github.com/libopenstorage/grpc-framework/pkg/grpc/client"
	appapi "github.com/libopenstorage/grpc-framework/test/app/api"
	appserver "github.com/libopenstorage/grpc-framework/test/app/pkg/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	grpcSocket = "/tmp/grpc-framework-testServer.sock"
)

type testServer struct {
	config *ServerConfig
	server *Server
	conn   *grpc.ClientConn
}

func (s *testServer) Stop() {
	s.server.Stop()
	s.conn.Close()
}

func (s *testServer) Address() string {
	return s.server.Address()
}

func (s *testServer) Conn() *grpc.ClientConn {
	return s.conn
}

func newDefaultConfig(t *testing.T) *ServerConfig {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
		Socket:  grpcSocket,
	}
	config.WithDefaultRestServer("9001").
		RegisterGrpcServers(func(gs *grpc.Server) {
			appapi.RegisterHelloGreeterServer(gs, &appserver.HelloGreeter{})
		})

	return config
}

func newTestServer(t *testing.T, config *ServerConfig) *testServer {

	if config.Socket != "" {
		os.Remove(config.Socket)
	}

	s, err := New(config)
	assert.NoError(t, err)

	err = s.Start()
	assert.NoError(t, err)

	// Setup connection to server
	conn, err := grpcclient.Connect(s.Address(), []grpc.DialOption{grpc.WithInsecure()})
	assert.NoError(t, err)

	return &testServer{
		config: config,
		server: s,
		conn:   conn,
	}
}

func newDefaultTestServer(t *testing.T) *testServer {
	return newTestServer(t, newDefaultConfig(t))
}

func TestSimpleServer(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()
}

func TestServerWithoutRest(t *testing.T) {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
		Socket:  grpcSocket,
	}

	s := newTestServer(t, config)
	assert.Nil(t, s.server.restGateway)
	assert.NotNil(t, s.server.udsServer)
	assert.NotNil(t, s.server.netServer)
	defer s.Stop()
}

func TestServerWithoutUdsAndRest(t *testing.T) {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
	}

	s := newTestServer(t, config)
	assert.Nil(t, s.server.restGateway)
	assert.Nil(t, s.server.udsServer)
	assert.NotNil(t, s.server.netServer)
	defer s.Stop()
}

func TestServerErrorWhenRestWithoutUds(t *testing.T) {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
	}
	config.WithDefaultRestServer("9001").
		RegisterGrpcServers(func(gs *grpc.Server) {
			appapi.RegisterHelloGreeterServer(gs, &appserver.HelloGreeter{})
		})
	_, err := New(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must provide unix domain socket for REST")
}

func TestSimpleServerLockTest(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()

	value := 0
	err := s.server.Transaction(func() error {
		value = 1
		return nil
	})
	assert.Equal(t, value, 1)
	assert.NoError(t, err)

	err = s.server.Transaction(func() error {
		value = 2
		return fmt.Errorf("ERROR")
	})
	assert.Equal(t, value, 2)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ERROR")
}

func rateLimiterShowsDenial(t *testing.T, s *testServer) bool {
	denied := false
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			ctx := context.Background()
			conn, err := grpcclient.Connect(s.Address(), []grpc.DialOption{grpc.WithInsecure()})
			g := appapi.NewHelloGreeterClient(conn)

			for {
				_, err = g.SayHello(ctx, &appapi.HelloGreeterSayHelloRequest{})
				if err != nil {
					serverError, ok := status.FromError(err)
					assert.True(t, ok)
					assert.Equal(t, serverError.Code(), codes.ResourceExhausted)

					if serverError.Code() == codes.ResourceExhausted {
						time.Sleep(time.Millisecond * time.Duration(i))
						denied = true
					}
					continue
				}
				break
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	return denied
}

func TestServerRateLimiter(t *testing.T) {
	c := newDefaultConfig(t)
	c.WithRateLimiter(rate.NewLimiter(2, 2))
	s := newTestServer(t, c)
	defer s.Stop()

	assert.True(t, rateLimiterShowsDenial(t, s))
}

func TestServerNoRateLimiter(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()

	assert.False(t, rateLimiterShowsDenial(t, s))
}
