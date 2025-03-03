package app

import (
	"backend/shared/base"
	"context"
	"errors"
	"time"
)

// Errors
var (
	ErrTokenInvalid    = errors.New("token is invalid")
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenGeneration = errors.New("failed to generate token")
)

// TokenProvider defines the interface for Keycloak token operations
type TokenProvider interface {
	GenerateWebSocketToken(ctx context.Context, tenantID, userID string) (string, time.Time, error)
	ValidateToken(ctx context.Context, token string) (valid bool, err error)
	RevokeToken(ctx context.Context, token string) error
}

// TokenService handles token operations using Keycloak
type TokenService struct {
	base     *base.BaseService
	provider TokenProvider
}

// New creates a new TokenService
func New(base *base.BaseService, provider TokenProvider) *TokenService {
	return &TokenService{
		base:     base,
		provider: provider,
	}
}

// GenerateWebSocketToken generates a short-lived token for WebSocket authentication
func (s *TokenService) GenerateWebSocketToken(ctx context.Context, tenantID, userID string) (string, error) {
	// Generate token using Keycloak
	tokenStr, _, err := s.provider.GenerateWebSocketToken(ctx, tenantID, userID)
	if err != nil {
		s.base.Logger.Error(ctx, err, "Failed to generate WebSocket token",
			map[string]interface{}{
				"tenantID": tenantID,
				"userID":   userID,
			})
		return "", ErrTokenGeneration
	}

	return tokenStr, nil
}

// ValidateToken validates a token
func (s *TokenService) ValidateToken(ctx context.Context, tokenStr string) error {
	// Validate token with Keycloak
	valid, err := s.provider.ValidateToken(ctx, tokenStr)
	if err != nil {
		s.base.Logger.Error(ctx, err, "Error validating token", nil)
		return ErrTokenInvalid
	}

	if !valid {
		s.base.Logger.Info(ctx, "Token validation failed", nil)
		return ErrTokenInvalid
	}

	return nil
}
