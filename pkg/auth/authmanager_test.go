/*
Copyright 2023 Portworx

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
	"context"
	"fmt"
	"testing"

	"github.com/libopenstorage/grpc-framework/pkg/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

type testAuthenticator struct {
	iss             string
	aud             string
	username        string
	verficationFail bool
}

func (a *testAuthenticator) AuthenticateToken(ctx context.Context, token string) (*Claims, error) {
	if a.verficationFail {
		return nil, fmt.Errorf("bad signature")
	}

	return &Claims{
		Issuer:   a.iss,
		Audience: a.aud,
	}, nil
}
func (a *testAuthenticator) Username(claims *Claims) string {
	if a.username == "" {
		return "test"
	}
	return a.username
}

func TestAddAuthenticatorFailures(t *testing.T) {
	a := NewAuthenticatorManagerDefault()

	err := a.AddAuthenticator("", "aud", &testAuthenticator{})
	assert.Error(t, err)

	err = a.AddAuthenticatorForIssuer("", &testAuthenticator{})
	assert.Error(t, err)

	err = a.AddAuthenticatorWithKey(AuthenticatorKey{}, &testAuthenticator{})
	assert.Error(t, err)
}

func TestAuthenticatorManagerDefault(t *testing.T) {
	tests := []struct {
		iss      string
		aud      string
		username string
		fail     bool
	}{
		{
			iss:      "one",
			aud:      "two",
			username: "me",
		},
		{
			iss:      "one",
			username: "you",
		},
		{
			iss:      "two",
			username: "them",
		},
		{
			iss:      "three",
			aud:      "two",
			username: "us",
		},
		{
			iss:  "three",
			fail: true,
		},
		{
			iss:  "unknown",
			fail: true,
		},
		{
			iss:  "one",
			aud:  "notforme",
			fail: true,
		},
		{
			iss:  "two",
			aud:  "notformeeither",
			fail: true,
		},
		{
			iss:  "who",
			aud:  "notme",
			fail: true,
		},
	}

	a := NewAuthenticatorManagerDefault()
	assert.NotNil(t, a)

	inserted_authenticators := 0
	for _, test := range tests {
		if test.fail {
			continue
		}
		inserted_authenticators++

		if test.aud != "" {
			a.AddAuthenticator(test.iss, test.aud, &testAuthenticator{
				iss:      test.iss,
				aud:      test.aud,
				username: test.username,
			})
		} else {
			a.AddAuthenticatorForIssuer(test.iss, &testAuthenticator{
				iss:      test.iss,
				username: test.username,
			})
		}
	}

	assert.Equal(t, a.Len(), inserted_authenticators)

	sign, err := NewSignatureSharedSecret("test")
	assert.NoError(t, err)

	// Create a token
	for _, test := range tests {
		claims := &Claims{
			Issuer:   test.iss,
			Audience: test.aud,
			Subject:  "test",
			Email:    "test@test.com",
			Name:     "test",
		}

		token, err := Token(claims, sign, &Options{})
		assert.NoError(t, err)

		ctx, err := a.AuthenticateToken(context.Background(), token)
		if test.fail {
			assert.Error(t, err)
			s := util.FromError(err)
			assert.Equal(t, codes.Unauthenticated, s.Code())
			if test.aud != "" {
				assert.Contains(t, s.Message(), "for audience")
			} else {
				assert.Contains(t, s.Message(), "is not a trusted issuer")
			}
			continue
		}
		assert.NoError(t, err)
		u, ok := NewUserInfoFromContext(ctx)
		assert.True(t, ok)

		assert.Equal(t, test.iss, u.Claims.Issuer)
		assert.Equal(t, test.aud, u.Claims.Audience)
		assert.Equal(t, test.username, u.Username)
	}
}
