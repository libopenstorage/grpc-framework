/*
Package defaultcontext manage the default context and timeouts
Copyright 2021 Portworx

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
package defaultcontext

import (
	"testing"
	"time"

	grpcgw "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
)

func TestDefaultContextManager(t *testing.T) {
	d := Inst()
	assert.NotNil(t, d)
	assert.Equal(t, d.GetDefaultTimeout(), defaultDuration)
	assert.Equal(t, grpcgw.DefaultContextTimeout, defaultDuration)

	timeout := 100 * time.Hour
	err := d.SetDefaultTimeout(timeout)
	defer d.SetDefaultTimeout(defaultDuration)
	assert.NoError(t, err)
	assert.Equal(t, d.GetDefaultTimeout(), timeout)
	assert.Equal(t, grpcgw.DefaultContextTimeout, timeout)
}
