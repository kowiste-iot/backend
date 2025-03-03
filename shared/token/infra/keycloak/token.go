package keycloak

import (
	"backend/shared/keycloak"
	"backend/shared/token/domain"
	"context"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v13"
)

// TokenKeycloak implements token operations using Keycloak
type TokenKeycloak struct {
	*keycloak.Keycloak
	tokenConfig *domain.TokenConfiguration
}

const (
	clientID     = "opel"
	clientSecret = "gi93Z9DzGLgJodlF9adPUJ4kCIHx26l3"
	user         = "ae55f59b-5f9c-4ed7-81d8-69439952ef2f"
	tenant       = "opel"
)

// New creates a new TokenKeycloak provider
func New(cfg *domain.TokenConfiguration, core *keycloak.Keycloak) *TokenKeycloak {
	return &TokenKeycloak{
		Keycloak:    core,
		tokenConfig: cfg,
	}
}

// GenerateWebSocketToken generates a short-lived token specifically for WebSocket connections
func (tk *TokenKeycloak) GenerateWebSocketToken(ctx context.Context, tenantID, userID string) (string, time.Time, error) {
	// Get admin token
	adminToken, err := tk.GetValidToken(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get admin token: %w", err)
	}
	tokenParams := gocloak.TokenOptions{
		GrantType:          gocloak.StringP("urn:ietf:params:oauth:grant-type:token-exchange"),
		SubjectToken:       gocloak.StringP(adminToken.AccessToken),
		RequestedSubject:   gocloak.StringP(user),
		Audience:           gocloak.StringP(tenant),
	}
	t, err := tk.Client.GetToken(ctx, tenant, tokenParams)
	if err != nil {
		return "", time.Time{}, err
	}
	println("t", t)
	// Use token exchange to get a token for the specific user with the websocket client as audience
	tokenSet, err := tk.Client.LoginClientTokenExchange(
		ctx,
		tk.Config.ClientID,     // Client ID (admin client)
		adminToken.AccessToken, // Token to exchange (admin token)
		tk.Config.ClientSecret, // Client secret (admin client)
		tenantID,               // Realm
		clientID,               // Target client (audience)
		userID,                 // User ID
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Set expiration time based on token lifetime config
	expiresAt := time.Now().Add(time.Duration(tk.tokenConfig.TokenLifetime) * time.Second)

	return tokenSet.AccessToken, expiresAt, nil
}

// ValidateToken validates a token
func (tk *TokenKeycloak) ValidateToken(ctx context.Context, tokenStr string) (bool, error) {
	// Use Keycloak's introspection endpoint to validate the token
	rptResult, err := tk.Client.RetrospectToken(
		ctx,
		tokenStr,
		clientID,
		clientSecret,
		tenant,
	)
	if err != nil {
		return false, fmt.Errorf("failed to introspect token: %w", err)
	}
	err = tk.RevokeToken(ctx, tokenStr)
	if err != nil {
		return false, fmt.Errorf("failed to revoke token: %w", err)
	}
	// Check if the token is active (not expired, not revoked)
	if !*rptResult.Active {
		return false, nil
	}

	return true, nil
}

// RevokeToken revokes a token in Keycloak
// Note: This is a best-effort operation as Keycloak doesn't provide a direct way to revoke
// a specific access token. For our short-lived tokens, expiration will handle most cases.
func (tk *TokenKeycloak) RevokeToken(ctx context.Context, tokenStr string) error {
	return tk.Client.RevokeToken(ctx, tenant, clientID, clientSecret, tokenStr)
}
