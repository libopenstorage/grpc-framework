/*
Package auth can be used for authentication and authorization
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
package auth

import (
	"context"

	"github.com/libopenstorage/grpc-framework/pkg/correlation"
	"github.com/sirupsen/logrus"
)

func init() {
	correlation.RegisterComponent(correlation.ComponentGrpcAuth)
}

const (
	systemGuestRoleName = "system.guest"
)

var (
	systemTokenInst TokenGenerator = &noauth{}

	// Inst returns the instance of system token manager.
	// This function can be overridden for testing purposes
	InitSystemTokenManager = func(tg TokenGenerator) {
		systemTokenInst = tg
	}

	// SystemTokenManagerInst returns the systemTokenManager instance
	SystemTokenManagerInst = func() TokenGenerator {
		return systemTokenInst
	}
)

// AuthenticatorKey is the hash key used to get a specific authenticator
type AuthenticatorKey struct {
	// Issuer is the issuer of the token to match for a speific authenticator.
	// This value is *required*
	Issuer string
	// Audience is the value to match to the token 'aud' which maps to a
	// specific authenticator.
	// This value is (optional)
	Audience string
}

// AuthenticatorManager interface groups token validators to help authenticate tokens
type AuthenticatorManager interface {
	// AuthenticateToken authenticates the token and modifies the context with
	// the user information
	AuthenticateToken(context.Context, string) (context.Context, error)

	// AddAuthenticator adds an authenticator to the manager
	// for a specific issuer and audience
	AddAuthenticator(issuer string, audience string, a Authenticator) error

	// AddAuthenticatorWithKey adds an authenticator to the manager using a Key
	AddAuthenticatorWithKey(AuthenticatorKey, Authenticator) error

	// AddAuthenticatorForIssuer adds an authenticator to the manager
	// for a specific issuer only and does not include an audidence
	AddAuthenticatorForIssuer(string, Authenticator) error

	// Len returns the number of authenticators being managed
	Len() int

	// Log prints the authenticators to the log
	Log(*logrus.Entry)
}

// Authenticator interface validates and extracts the claims from a raw token
type Authenticator interface {
	// AuthenticateToken validates the token and returns the claims
	AuthenticateToken(context.Context, string) (*Claims, error)

	// Username returns the unique id according to the configuration. Default
	// it will return the value for "sub" in the token claims, but it can be
	// configured to return the email or name as the unique id.
	Username(*Claims) string
}

// Enabled returns whether or not auth is enabled.
func Enabled() bool {
	return len(systemTokenInst.Issuer()) != 0
}
