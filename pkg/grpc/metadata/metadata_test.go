/*
Package server is a generic gRPC server manager
Copyright 2020 Pure Storage

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

package metadata

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stretchr/testify/assert"
)

func TestContextMetadata(t *testing.T) {

	// setup context
	ctx := AddMetadataToContext(context.Background(), "hello", "world")
	ctx = AddMetadataToContext(ctx, "jay", "kay")
	ctx = AddMetadataToContext(ctx, "one", "two")

	// TODO: Replace this manual conversion to an actual grpc call
	outgoingMd := metautils.ExtractOutgoing(ctx)
	incomingCtx := outgoingMd.ToIncoming(context.Background())

	assert.Equal(t, GetMetadataValueFromKey(incomingCtx, "hello"), "world")
	assert.Equal(t, GetMetadataValueFromKey(incomingCtx, "jay"), "kay")
	assert.Equal(t, GetMetadataValueFromKey(incomingCtx, "one"), "two")
}
