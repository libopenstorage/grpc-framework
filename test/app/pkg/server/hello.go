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

	"github.com/libopenstorage/grpc-framework/test/app/api"
	"github.com/sirupsen/logrus"
)

type HelloGreeter struct {
	api.HelloGreeterServer
	api.HelloIdentityServer
}

func (h *HelloGreeter) SayHello(
	ctx context.Context,
	req *api.HelloGreeterSayHelloRequest,
) (*api.HelloGreeterSayHelloResponse, error) {
	logrus.Info("Received a request in SayHello()")

	return &api.HelloGreeterSayHelloResponse{
		Message: fmt.Sprintf("Hello, %s", req.GetName()),
	}, nil
}

func (h *HelloGreeter) Version(
	ctx context.Context,
	req *api.HelloIdentityVersionRequest,
) (*api.HelloIdentityVersionResponse, error) {
	logrus.Info("Received request for version")
	return &api.HelloIdentityVersionResponse{
		HelloVersion: &api.HelloVersion{
			Major: int32(api.HelloVersion_MAJOR),
			Minor: int32(api.HelloVersion_MINOR),
			Patch: int32(api.HelloVersion_PATCH),
			Version: fmt.Sprintf("%d.%d.%d",
				api.HelloVersion_MAJOR,
				api.HelloVersion_MINOR,
				api.HelloVersion_PATCH),
		},
	}, nil
}
