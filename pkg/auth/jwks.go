/*
Copyright 2022 Portworx

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

	oidc "github.com/coreos/go-oidc/v3/oidc"
)

// From the oidc package:
// NewRemoteKeySet returns a KeySet that can validate JSON web tokens by using HTTP
// GETs to fetch JSON web token sets hosted at a remote URL. This is automatically
// used by NewProvider using the URLs returned by OpenID Connect discovery, but is
// exposed for providers that don't support discovery or to prevent round trips to the
// discovery URL.
//
// The returned KeySet is a long lived verifier that caches keys based on any
// keys change. Reuse a common remote key set instead of creating new ones as needed.

// JWKSAuthConfig configures an JWKS connection
type JWKSAuthConfig struct {
	// Issuer of the tokens
	// e.g. https://accounts.google.com
	Issuer string
	// JWKSUrl is the actual URL to the public key in jwks format
	// e.g. https://www.googleapis.com/oauth2/v3/certs
	JWKSUrl string
	// UsernameClaim has the location of the unique id for the user.
	// If empty, "sub" will be used for the user name unique id.
	UsernameClaim UsernameClaimType
	// Namespace sets the namespace for all custom claims. For example
	// if the claims had the key: "https://mynamespace/roles", then
	// the namespace would be "https://mynamespace/".
	Namespace string
}

// JWKSAuthenticator is used to validate tokens with an JWKS
type JWKSAuthenticator struct {
	OIDCAuthenticator

	jwksUrl string
	keyset  *oidc.RemoteKeySet
}

// NewJWKS returns a new JWKS authenticator
// Example:
//
//      c := &JWKSAuthConfig{
//          Issuer:  "http://localhost",
//          JWKSUrl: "https://localhost:3080/.well-known/jwks.json",
//      }
//      a, err := NewJWKSAuthenticator(c)
//
func NewJWKSAuthenticator(config *JWKSAuthConfig) (*JWKSAuthenticator, error) {
	keyset := oidc.NewRemoteKeySet(context.Background(), config.JWKSUrl)
	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}
	verifier := oidc.NewVerifier(config.Issuer, keyset, oidcConfig)
	return &JWKSAuthenticator{
		OIDCAuthenticator: OIDCAuthenticator{
			url:           config.Issuer,
			verifier:      verifier,
			usernameClaim: config.UsernameClaim,
			namespace:     config.Namespace,
		},
		jwksUrl: config.JWKSUrl,
		keyset:  keyset,
	}, nil
}
