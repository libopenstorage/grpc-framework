/*
Package auth can be used for authentication and authorization
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
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthenticatorManagerDefault manages authenticators for specific values
// in the token like 'iss' and 'aud'
type AuthenticatorManagerDefault struct {
	lock           sync.Mutex
	authenticators map[AuthenticatorKey]Authenticator
}

// NewAuthenticatorManagerDefault returns a default authenticator manager that maps
// token information like issuer and audience to specific authenticators
func NewAuthenticatorManagerDefault() *AuthenticatorManagerDefault {
	return &AuthenticatorManagerDefault{
		authenticators: make(map[AuthenticatorKey]Authenticator),
	}
}

func (a *AuthenticatorManagerDefault) AuthenticateToken(ctx context.Context, token string) (context.Context, error) {

	// Get token issuer
	issuer, err := TokenIssuer(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Create a key
	key := AuthenticatorKey{
		Issuer: issuer,
	}

	// Get audience if any
	audience, err := TokenAudience(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if audience != "" {
		key.Audience = audience
	}

	// Get Authenticator
	if authenticator, ok := a.authenticators[key]; ok {
		claims, err := authenticator.AuthenticateToken(ctx, token)
		if err == nil {
			// Token was validated and claims are available.
			// Add authorization information back into the context so that other
			// functions can get access to this information.
			// If this is in the context is how functions will know that security is enabled.
			ctx = ContextSaveUserInfo(ctx, &UserInfo{
				Username: authenticator.Username(claims),
				Claims:   *claims,
			})
			return ctx, nil
		}

		// Error validating token
		return nil, status.Errorf(codes.PermissionDenied, "Unable to authenticate token")
	}

	// Unable to find the authenticator
	if key.Audience != "" {
		return nil, status.Errorf(codes.Unauthenticated, "Issuer %s for audience %s is unaccepted",
			key.Issuer,
			key.Audience)
	}
	return nil, status.Errorf(codes.Unauthenticated, "%s is not a trusted issuer", issuer)
}

func (a *AuthenticatorManagerDefault) AddAuthenticatorWithKey(key AuthenticatorKey, obj Authenticator) error {
	if key.Issuer == "" {
		return fmt.Errorf("must supply an issuer")
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	a.authenticators[key] = obj

	return nil
}

func (a *AuthenticatorManagerDefault) AddAuthenticatorForIssuer(issuer string, obj Authenticator) error {
	if issuer == "" {
		return fmt.Errorf("must supply an issuer")
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	a.authenticators[AuthenticatorKey{Issuer: issuer}] = obj

	return nil
}

func (a *AuthenticatorManagerDefault) AddAuthenticator(issuer, audience string, obj Authenticator) error {
	if issuer == "" {
		return fmt.Errorf("must supply an issuer")
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	a.authenticators[AuthenticatorKey{
		Issuer:   issuer,
		Audience: audience,
	}] = obj

	return nil
}

func (a *AuthenticatorManagerDefault) Len() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return len(a.authenticators)
}

func (a *AuthenticatorManagerDefault) Log(l *logrus.Entry) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for key := range a.authenticators {
		if key.Audience != "" {
			l.Infof("Authentication enabled for issuer %s for audience %s",
				key.Issuer, key.Audience)
		} else {
			l.Infof("Authentication enabled for issuer: %s", key.Issuer)
		}
	}
}
