package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestMultiAuthenticatorWithClientID_GetAuthenticator_Ok(t *testing.T) {
	// Given.
	jwtAuth, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)
	m := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{IssuerWithClientID{"issuer-1", "1"}: jwtAuth})

	// When.
	found := m.GetAuthenticator("issuer-1", "1")

	// Then.
	assert.NotNil(t, found)
}

func TestMultiAuthenticatorWithClientID_GetAuthenticator_FailNoClientId(t *testing.T) {
	// Given.
	jwtAuth, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)
	m := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{
		IssuerWithClientID{"issuer-1", "1"}: jwtAuth,
	})

	// When.
	found := m.GetAuthenticator("issuer-1", "2")

	// Then.
	assert.Nil(t, found)
}

func TestMultiAuthenticatorWithClientID_GetAuthenticator_FailNoIssuer(t *testing.T) {
	// Given.
	m := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{})

	// When.
	found := m.GetAuthenticator("issuer-1", "1")

	// Then.
	assert.Nil(t, found)
}

func TestMultiAuthenticatorWithClientID_AuthenticateToken_Ok(t *testing.T) {
	// Given.
	jwtAuth1, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret-1"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	m := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{
		IssuerWithClientID{"issuer-1", "1"}: jwtAuth1,
		IssuerWithClientID{"issuer-1", "2"}: jwtAuth1,
	})

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

func TestMultiAuthenticatorWithClientID_AuthenticateTokenMultiAud_Ok(t *testing.T) {
	// Given.
	jwtAuth1, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret-1"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	m := NewMultiAuthenticatorWithClientID(map[IssuerWithClientID]Authenticator{
		IssuerWithClientID{"issuer-1", ""}:          jwtAuth1,
		IssuerWithClientID{"issuer-1", "clientID2"}: jwtAuth1,
	})

	rawToken, err := Token(&Claims{
		Issuer:   "issuer-1",
		Email:    "my@email.com",
		Name:     "my-name",
		Subject:  "my-sub",
		Roles:    []string{"tester"},
		Audience: "extra info & clientID2",
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
