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
package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoAuth(t *testing.T) {
	// Check for interface implementation
	var na TokenGenerator = &noauth{}

	// Get no auth Instance
	na = NoAuth()
	assert.NotNil(t, na)

	assert.Equal(t, na.Issuer(), "")
	authctr, err := na.GetAuthenticator()
	assert.Error(t, err)
	assert.Nil(t, authctr)
	token, err := na.GetToken(&Options{})
	assert.NoError(t, err)
	assert.Equal(t, token, "")
}
