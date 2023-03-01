package auth

import (
	"context"
	"fmt"
)

type multiAuthenticatorImpl struct {
	authenticators map[string][]Authenticator
}

// NewMultiAuthenticator maintains a list of authenticators for a given issuer.
// The input argument is a map of issuers to a list of authenticators for that issuer.
// NOTE: This interface does not check if there are duplicate authenticators.
func NewMultiAuthenticator(
	authenticators map[string][]Authenticator,
) (MultiAuthenticatorWithClientID, error) {
	authMap := make(map[string][]Authenticator)
	for issuer, authenticatorsList := range authenticators {
		if len(authenticatorsList) == 0 {
			return nil, fmt.Errorf("empty authenticators list for issuer %v", issuer)
		}
		authMap[issuer] = authenticatorsList
	}
	return &multiAuthenticatorImpl{
		authenticators: authMap,
	}, nil
}

func (m *multiAuthenticatorImpl) GetAuthenticators(issuer string) []Authenticator {
	return m.authenticators[issuer]
}

func (m *multiAuthenticatorImpl) ListIssuers() []string {
	var issuers []string
	for issuer, _ := range m.authenticators {
		issuers = append(issuers, issuer)
	}
	return issuers
}

func (m *multiAuthenticatorImpl) AuthenticateToken(ctx context.Context, idToken string) (*Claims, error) {
	tokenClaims, err := TokenClaims(idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims from token: %w", err)
	}

	for _, authenticator := range m.GetAuthenticators(tokenClaims.Issuer) {
		claims, err := authenticator.AuthenticateToken(ctx, idToken)
		if err == nil {
			return claims, nil
		}
	}

	return nil, fmt.Errorf("failed to authenticate token for issuer %v and audience %v",
		tokenClaims.Issuer, tokenClaims.Audience)
}

func (m *multiAuthenticatorImpl) Username(claims *Claims) string {
	var username string
	// TODO: This code does not handle the case where there are multiple authenticators
	// registered with the same issuer but different UsernameClaimTypes.
	for _, authenticator := range m.GetAuthenticators(claims.Issuer) {
		if username = authenticator.Username(claims); username != "" {
			return username
		}
	}
	return username
}
