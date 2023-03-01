package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewMultiAuthenticator_EmptyAuthenticators(t *testing.T) {
	m, err := NewMultiAuthenticator(map[string][]Authenticator{
		"issuer-1": []Authenticator{},
	})
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestNewMultiAuthenticator_ListIssuers(t *testing.T) {
	// Given.
	jwtAuth, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)
	m, err := NewMultiAuthenticator(map[string][]Authenticator{
		"issuer-1": []Authenticator{jwtAuth},
		// TODO: MultiAuthenticator does not check for duplicate authenticators across issuers.
		"issuer-2": []Authenticator{jwtAuth},
	})
	assert.NoError(t, err)

	// When.
	issuers := m.ListIssuers()

	// Then.
	assert.Equal(t, 2, len(issuers))
}
func TestMultiAuthenticator_GetAuthenticators_Ok(t *testing.T) {
	// Given.
	jwtAuth1, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret1"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	jwtAuth2, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret2"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	m, err := NewMultiAuthenticator(map[string][]Authenticator{
		"issuer-1": []Authenticator{jwtAuth1, jwtAuth2},
	})
	assert.NoError(t, err)

	// When.
	authenticators := m.GetAuthenticators("issuer-1")

	// Then.
	assert.Equal(t, 2, len(authenticators))
}

func TestMultiAuthenticator_GetAuthenticator_FailNoIssuer(t *testing.T) {
	// Given.
	m, err := NewMultiAuthenticator(map[string][]Authenticator{})
	assert.NoError(t, err)

	// When.
	authenticators := m.GetAuthenticators("issuer-1")

	// Then.
	assert.Empty(t, authenticators)
}

func TestMultiAuthenticator_AuthenticateToken_Ok(t *testing.T) {
	// Given.
	jwtAuth1, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret-1"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	jwtAuth2, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret-2"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	m, err := NewMultiAuthenticator(map[string][]Authenticator{
		"issuer-1": []Authenticator{jwtAuth1, jwtAuth2},
	})
	assert.NoError(t, err)

	rawToken1, err := Token(&Claims{
		Audience: "1",
		Issuer:   "issuer-1",
		Email:    "my@email.com",
		Name:     "my-name",
		Subject:  "my-sub",
		Roles:    []string{"tester"},
	}, &Signature{
		Type: jwt.SigningMethodHS256,
		Key:  []byte("my-secret-2"),
	}, &Options{
		Expiration: time.Now().Add(time.Minute * 10).Unix(),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, rawToken1)

	// When.
	authenticateClaims, err := m.AuthenticateToken(context.Background(), rawToken1)

	// Then.
	assert.NoError(t, err)
	assert.Equal(t, "issuer-1", authenticateClaims.Issuer)
	assert.Equal(t, "my@email.com", authenticateClaims.Email)
	assert.Equal(t, "my-name", authenticateClaims.Name)
	assert.Equal(t, []string{"tester"}, authenticateClaims.Roles)

	rawToken2, err := Token(&Claims{
		Audience: "1",
		Issuer:   "issuer-1",
		Email:    "my@email.com",
		Name:     "my-name",
		Subject:  "my-sub",
		Roles:    []string{"tester"},
	}, &Signature{
		Type: jwt.SigningMethodHS256,
		Key:  []byte("my-secret-2"),
	}, &Options{
		Expiration: time.Now().Add(time.Minute * 10).Unix(),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, rawToken2)

	// When.
	authenticateClaims, err = m.AuthenticateToken(context.Background(), rawToken2)

	// Then.
	assert.NoError(t, err)
	assert.Equal(t, "issuer-1", authenticateClaims.Issuer)
	assert.Equal(t, "my@email.com", authenticateClaims.Email)
	assert.Equal(t, "my-name", authenticateClaims.Name)
	assert.Equal(t, []string{"tester"}, authenticateClaims.Roles)
}

func TestMultiAuthenticator_AuthenticateToken_Fail(t *testing.T) {
	// Given.
	jwtAuth1, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret-1"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	jwtAuth2, err := NewJwtAuthenticator(&JwtAuthConfig{
		SharedSecret:  []byte("my-secret-2"),
		UsernameClaim: UsernameClaimTypeSubject,
	})
	assert.NoError(t, err)

	m, err := NewMultiAuthenticator(map[string][]Authenticator{
		"issuer-1": []Authenticator{jwtAuth1, jwtAuth2},
	})
	assert.NoError(t, err)

	rawToken1, err := Token(&Claims{
		Audience: "1",
		Issuer:   "issuer-1",
		Email:    "my@email.com",
		Name:     "my-name",
		Subject:  "my-sub",
		Roles:    []string{"tester"},
	}, &Signature{
		Type: jwt.SigningMethodHS256,
		Key:  []byte("my-secret-3"),
	}, &Options{
		Expiration: time.Now().Add(time.Minute * 10).Unix(),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, rawToken1)

	// When.
	authenticateClaims, err := m.AuthenticateToken(context.Background(), rawToken1)

	// Then.
	assert.Error(t, err)
	assert.Nil(t, authenticateClaims)
}
