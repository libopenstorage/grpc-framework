package auth

import (
	"context"
	"fmt"
	"strings"
)

type IssuerWithClientID struct {
	Issuer   string // Issuer of the token.
	ClientID string // ClientID of the token, can be empty.
}

type multiAuthenticatorWithClientIDImpl struct {
	authenticators map[IssuerWithClientID]Authenticator
}

func NewMultiAuthenticatorWithClientID(
	authenticators map[IssuerWithClientID]Authenticator) MultiAuthenticatorWithClientID {
	return &multiAuthenticatorWithClientIDImpl{
		authenticators: authenticators,
	}
}

// GetAuthenticator returns nil if there is no authenticator for the given issuer and client ID.
func (m *multiAuthenticatorWithClientIDImpl) GetAuthenticator(issuer, audience string) Authenticator {
	for issuerWithClientID, authenticator := range m.authenticators {
		// Audience field in a claim can contain more than just the clientID.
		// The strings.Contains check here is same as what is being used by the OIDC library
		// which verifies the token.
		if issuer == issuerWithClientID.Issuer && strings.Contains(audience, issuerWithClientID.ClientID) {
			return authenticator
		}
	}
	return nil
}

func (m *multiAuthenticatorWithClientIDImpl) ListIssuersWithClientID() []IssuerWithClientID {
	var issuerWithClientIDs []IssuerWithClientID
	for issuerWithClientID, _ := range m.authenticators {
		issuerWithClientIDs = append(issuerWithClientIDs, issuerWithClientID)
	}
	return issuerWithClientIDs
}

func (m *multiAuthenticatorWithClientIDImpl) AuthenticateToken(ctx context.Context, accessToken string) (*Claims, error) {
	claims, err := TokenClaims(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims from token: %w", err)
	}

	authenticator := m.GetAuthenticator(claims.Issuer, claims.Audience)
	if authenticator == nil {
		return nil, fmt.Errorf("no authenticator found for issuer %s & audience %s", claims.Issuer, claims.Audience)
	}

	return authenticator.AuthenticateToken(ctx, accessToken)
}

func (m *multiAuthenticatorWithClientIDImpl) Username(claims *Claims) string {
	authenticator := m.GetAuthenticator(claims.Issuer, claims.Audience)
	if authenticator == nil {
		return ""
	}

	return authenticator.Username(claims)
}
