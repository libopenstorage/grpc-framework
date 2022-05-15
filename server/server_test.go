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
	"testing"

	grpcclient "github.com/libopenstorage/grpc-framework/pkg/grpc/client"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

func newDefaultTestServer(t *testing.T) *testServer {

	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
		Socket:  "/tmp/grpc-framework-testServer.sock",
	}
	config.WithDefaultRestServer("9001")

	return newTestServer(t, config)
}

func newTestServer(t *testing.T, config *ServerConfig) *testServer {

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

func TestSimpleServer(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()
}
