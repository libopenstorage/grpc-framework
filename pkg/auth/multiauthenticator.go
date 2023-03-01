package auth

import (
	"context"
	"fmt"
)

type IssuerWithClientID struct {
	// Issuer of the token. The issuer cannot be empty.
	Issuer string
	// ClientID of the token. The clientID cannot be empty.
	ClientID string
}

type multiAuthenticatorWithClientIDImpl struct {
	authenticators map[IssuerWithClientID]Authenticator
}

func NewMultiAuthenticatorWithClientID(
	authenticators map[IssuerWithClientID]Authenticator) (MultiAuthenticatorWithClientID, error) {
	authMap := make(map[IssuerWithClientID]Authenticator)
	for issuerWithClientID, authenticator := range authenticators {
		if issuerWithClientID.ClientID == "" {
			return nil, fmt.Errorf("clientID cannot be empty")
		}
		if issuerWithClientID.Issuer == "" {
			return nil, fmt.Errorf("issuer cannot be empty")
		}
		authMap[issuerWithClientID] = authenticator
	}
	return &multiAuthenticatorWithClientIDImpl{
		authenticators: authMap,
	}, nil
}

func (m *multiAuthenticatorWithClientIDImpl) ListIssuersWithClientID() []IssuerWithClientID {
	var issuerWithClientIDs []IssuerWithClientID
	for issuerWithClientID, _ := range m.authenticators {
		issuerWithClientIDs = append(issuerWithClientIDs, issuerWithClientID)
	}
	return issuerWithClientIDs
}

func (m *multiAuthenticatorWithClientIDImpl) GetAuthenticators(issuer string) []Authenticator {
	var authenticators []Authenticator
	for issuerWithClientID, authenticator := range m.authenticators {
		if issuerWithClientID.Issuer == issuer {
			authenticators = append(authenticators, authenticator)
		}
	}
	return authenticators
}

func (m *multiAuthenticatorWithClientIDImpl) AuthenticateToken(ctx context.Context, idToken string) (*Claims, error) {
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

func (m *multiAuthenticatorWithClientIDImpl) Username(claims *Claims) string {
	var username string
	// TODO: This code does not handle the case where there are mulitple authenticators
	// registered with the same issuer but different UsernameClaimTypes.
	for _, authenticator := range m.GetAuthenticators(claims.Issuer) {
		if username = authenticator.Username(claims); username != "" {
			return username
		}
	}
	return username
}
