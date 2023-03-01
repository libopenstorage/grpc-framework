package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewMultiAuthenticatorWithClientID_EmptyClientID(t *testing.T) {
	// Given.
	jwtAuth, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)
	m, err := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{IssuerWithClientID{"issuer-1", ""}: jwtAuth})
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestNewMultiAuthenticatorWithClientID_EmptyIssuer(t *testing.T) {
	// Given.
	jwtAuth, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)
	m, err := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{IssuerWithClientID{"", "1"}: jwtAuth})
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestMultiAuthenticatorWithClientID_GetAuthenticators_Ok(t *testing.T) {
	// Given.
	jwtAuth, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)
	m, err := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{
		IssuerWithClientID{"issuer-1", "1"}: jwtAuth,
		IssuerWithClientID{"issuer-1", "2"}: jwtAuth,
	})
	assert.NoError(t, err)

	// When.
	authenticators := m.GetAuthenticators("issuer-1")

	// Then.
	assert.Equal(t, 2, len(authenticators))
}

func TestMultiAuthenticatorWithClientID_GetAuthenticator_FailNoIssuer(t *testing.T) {
	// Given.
	m, err := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{})
	assert.NoError(t, err)

	// When.
	authenticators := m.GetAuthenticators("issuer-1")

	// Then.
	assert.Empty(t, authenticators)
}

func TestMultiAuthenticatorWithClientID_AuthenticateToken_Ok(t *testing.T) {
	// Given.
	jwtAuth1, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret-1"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	m, err := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{
		IssuerWithClientID{"issuer-1", "1"}: jwtAuth1,
		IssuerWithClientID{"issuer-1", "2"}: jwtAuth1,
	})
	assert.NoError(t, err)

	rawToken, err := Token(&Claims{
		Audience: "1",
		Issuer:   "issuer-1",
		Email:    "my@email.com",
		Name:     "my-name",
		Subject:  "my-sub",
		Roles:    []string{"tester"},
	}, &Signature{
		Type: jwt.SigningMethodHS256,
		Key:  []byte("my-secret-1"),
	}, &Options{
		Expiration: time.Now().Add(time.Minute * 10).Unix(),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, rawToken)

	// When.
	authenticateClaims, err := m.AuthenticateToken(context.Background(), rawToken)

	// Then.
	assert.NoError(t, err)
	assert.Equal(t, "issuer-1", authenticateClaims.Issuer)
	assert.Equal(t, "my@email.com", authenticateClaims.Email)
	assert.Equal(t, "my-name", authenticateClaims.Name)
	assert.Equal(t, []string{"tester"}, authenticateClaims.Roles)
}
