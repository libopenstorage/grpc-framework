package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMultiAuthenticator_EmptyAuthenticators(t *testing.T) {
	m, err := NewMultiAuthenticatorByClientID("issuer-1", map[string]Authenticator{})
	assert.Error(t, err)
	assert.Nil(t, m)

	m, err = NewIteratingMultiAuthenticator("issuer-1", []Authenticator{})
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestIteratingMultiAuthenticator_AuthenticateToken_Ok(t *testing.T) {
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

	m, err := NewIteratingMultiAuthenticator("issuer-1", []Authenticator{jwtAuth1, jwtAuth2})
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

func TestIteratingMultiAuthenticator_AuthenticateToken_Fail(t *testing.T) {
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

	m, err := NewIteratingMultiAuthenticator("issuer-1", []Authenticator{jwtAuth1, jwtAuth2})
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

func TestClientIDMultiAuthenticator_AuthenticateToken_Ok(t *testing.T) {
	testClientID := "1"
	testIssuer := "issuer-1"

	ctrl := gomock.NewController(t)
	a := NewMockAuthenticator(ctrl)
	a.EXPECT().AuthenticateToken(gomock.Any(), gomock.Any()).Return(&Claims{}, nil)

	ma, err := NewMultiAuthenticatorByClientID(testIssuer, map[string]Authenticator{testClientID: a})
	require.NoError(t, err)

	// Given.
	rawToken1, err := Token(&Claims{
		Audience: testClientID,
		Issuer:   testIssuer,
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
	require.NoError(t, err)
	require.NotEmpty(t, rawToken1)

	// When.
	_, err = ma.AuthenticateToken(context.TODO(), rawToken1)

	// Then.
	require.NoError(t, err)
}

func TestClientIDMultiAuthenticator_AuthenticateToken_WrongIssuer(t *testing.T) {
	testClientID := "1"
	testIssuer := "issuer-1"

	ctrl := gomock.NewController(t)
	a := NewMockAuthenticator(ctrl)

	ma, err := NewMultiAuthenticatorByClientID(testIssuer, map[string]Authenticator{testClientID: a})
	require.NoError(t, err)

	// Given.
	rawToken1, err := Token(&Claims{
		Audience: testClientID,
		Issuer:   "issuer-2",
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
	require.NoError(t, err)
	require.NotEmpty(t, rawToken1)

	// When.
	_, err = ma.AuthenticateToken(context.TODO(), rawToken1)

	// Then.
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not match issuer")
}

func TestClientIDMultiAuthenticator_AuthenticateToken_Multiple_Aud(t *testing.T) {
	testClientID := "1"
	testIssuer := "issuer-1"

	ctrl := gomock.NewController(t)
	a := NewMockAuthenticator(ctrl)

	ma, err := NewMultiAuthenticatorByClientID(testIssuer, map[string]Authenticator{testClientID: a})
	require.NoError(t, err)

	// Given.
	rawToken1, err := Token(&Claims{
		Audience: []string{testClientID, "extraAud"},
		Issuer:   testIssuer,
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
	require.NoError(t, err)
	require.NotEmpty(t, rawToken1)

	// When.
	_, err = ma.AuthenticateToken(context.TODO(), rawToken1)

	// Then.
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to determine client ID")
}
