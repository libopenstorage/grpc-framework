/*
Copyright 2019 Portworx

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
	"github.com/stretchr/testify/require"
)

func TestClaims(t *testing.T) {
	email, name, subject := "a@b.com", "hello", "123"
	claims := &Claims{
		Email:   email,
		Name:    name,
		Subject: subject,
	}

	// Check getUsername for correctness

	// default claim type is sub
	un, err := claims.GetUsername()
	assert.NoError(t, err)
	assert.Equal(t, un, subject)

	// claim type = email
	claims.UsernameClaim = UsernameClaimTypeEmail
	un, err = claims.GetUsername()
	assert.NoError(t, err)
	assert.Equal(t, un, email)

	// claim type = name
	claims.UsernameClaim = UsernameClaimTypeName
	un, err = claims.GetUsername()
	assert.NoError(t, err)
	assert.Equal(t, un, name)

	// claim type = sub
	claims.UsernameClaim = UsernameClaimTypeSubject
	un, err = claims.GetUsername()
	assert.NoError(t, err)
	assert.Equal(t, un, subject)
}

func TestValidateUsername(t *testing.T) {
	email, name, subject := "a@b.com", "hello", "123"
	goodClaims := &Claims{
		Email:   email,
		Name:    name,
		Subject: subject,
	}
	badClaims := &Claims{
		Email:   "",
		Name:    "",
		Subject: "",
	}

	typesToTest := []UsernameClaimType{
		UsernameClaimTypeEmail,
		UsernameClaimTypeName,
		UsernameClaimTypeSubject,
		UsernameClaimTypeDefault,
	}
	for _, unType := range typesToTest {
		goodClaims.UsernameClaim = unType
		err := goodClaims.ValidateUsername()
		assert.NoError(t, err)

		badClaims.UsernameClaim = unType
		err = badClaims.ValidateUsername()
		assert.Error(t, err)
	}
}

func TestGetAudience(t *testing.T) {
	claims := &Claims{}

	// no audience
	actual, err := claims.GetAudience()
	require.NoError(t, err)
	require.Empty(t, actual)

	// string
	claims.Audience = "1"
	actual, err = claims.GetAudience()
	require.NoError(t, err)
	require.Equal(t, []string{"1"}, actual)

	// []string single value
	actual, err = claims.GetAudience()
	require.NoError(t, err)
	require.Equal(t, []string{"1"}, actual)

	// []string with multiple values
	claims.Audience = []string{"1", "2"}
	actual, err = claims.GetAudience()
	require.NoError(t, err)
	require.Equal(t, []string{"1", "2"}, actual)

	// []interface{}
	claims.Audience = []interface{}{"1", "2", "3"}
	actual, err = claims.GetAudience()
	require.NoError(t, err)
	require.Equal(t, []string{"1", "2", "3"}, actual)

	// wrong type: int
	claims.Audience = 1
	actual, err = claims.GetAudience()
	require.Error(t, err)
	require.Nil(t, actual)

	// wrong type: []interface{} with int
	claims.Audience = []interface{}{1}
	actual, err = claims.GetAudience()
	require.Error(t, err)
	require.Nil(t, actual)
}
