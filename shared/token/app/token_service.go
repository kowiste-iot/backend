package app

import (
	appTenant "backend/internal/features/tenant/app"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/http/httputil"

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
	GenerateWebSocketToken(ctx context.Context, tenantID, branchID, userID string) (string, time.Time, error)
	ValidateToken(ctx context.Context, token string) (valid bool, userID string, err error)
	RevokeToken(ctx context.Context, token string) error
}

// TokenService handles token operations using Keycloak
type TokenService struct {
	base     *base.BaseService
	provider TokenProvider
	tenant   appTenant.TenantService
}

type ServiceDependencies struct {
	Provider     TokenProvider
	TenatService appTenant.TenantService
}

// New creates a new TokenService
func New(base *base.BaseService, deps *ServiceDependencies) *TokenService {
	return &TokenService{
		base:     base,
		provider: deps.Provider,
		tenant:   deps.TenatService,
	}
}

// GenerateWebSocketToken generates a short-lived token for WebSocket authentication
func (s *TokenService) GenerateWebSocketToken(ctx context.Context, tenantID, branchID, userID string) (string, error) {
	// Generate token using Keycloak
	tokenStr, _, err := s.provider.GenerateWebSocketToken(ctx, tenantID, branchID, userID)
	if err != nil {
		s.base.Logger.Error(ctx, err, "Failed to generate WebSocket token",
			map[string]interface{}{
				"tenantID": tenantID,
				"branchID": branchID,
				"userID":   userID,
			})
		return "", ErrTokenGeneration
	}

	return tokenStr, nil
}

// ValidateToken validates a token
func (s *TokenService) ValidateToken(ctx context.Context, input baseCmd.BaseInput, tokenStr string) (userID string, err error) {
	ctx = httputil.SetBranch(ctx, input.BranchName)
	tenant, err := s.tenant.GetTenant(ctx, input.TenantDomain)
	if err != nil {
		s.base.Logger.Error(ctx, err, "Error getting tenant", nil)
		return "", err
	}
	ctx = httputil.SetTenant(ctx, tenant)
	// Validate token with Keycloak
	valid, userID, err := s.provider.ValidateToken(ctx, tokenStr)
	if err != nil {
		s.base.Logger.Error(ctx, err, "Error validating token", nil)
		return "", ErrTokenInvalid
	}

	if !valid {
		s.base.Logger.Info(ctx, "Token validation failed", nil)
		return "", ErrTokenInvalid
	}
	//remove token after use
	err = s.provider.RevokeToken(ctx, tokenStr)
	return
}
