package inmemmo

import (
	"backend/shared/base"
	"backend/shared/http/httputil"
	"backend/shared/token/domain"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// TokenStore defines the interface for token storage operations
type TokenStore interface {
	// Store saves a token with its metadata
	Store(token string, tenant, branch, userID string, expiry time.Time) error
	// Get retrieves token metadata if exists
	Get(token string) (tenant, branch, userID string, expiry time.Time, exists bool)
	// Delete removes a token from storage
	Delete(token string) error
}

// InMemoryTokenStore implements TokenStore using an in-memory map
type InMemoryTokenStore struct {
	mu    sync.RWMutex
	store map[string]tokenData
}

type tokenData struct {
	TenantID string
	BranchID string
	UserID   string
	Expiry   time.Time
}

// NewInMemoryTokenStore creates a new in-memory token store
func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		store: make(map[string]tokenData),
	}
}

// Store saves a token with its metadata
func (s *InMemoryTokenStore) Store(token string, tenant, branch, userID string, expiry time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[token] = tokenData{
		TenantID: tenant,
		BranchID: branch,
		UserID:   userID,
		Expiry:   expiry,
	}
	return nil
}

// Get retrieves token metadata if exists
func (s *InMemoryTokenStore) Get(token string) (tenantID, branchID, userID string, expiry time.Time, exists bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.store[token]
	if !ok {
		return "", "", "", time.Time{}, false
	}
	return data.TenantID, data.BranchID, data.UserID, data.Expiry, true
}

// Delete removes a token from storage
func (s *InMemoryTokenStore) Delete(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, token)
	return nil
}

// TokenProvider implements the TokenProvider interface
type TokenProvider struct {
	base       *base.BaseService
	config     domain.TokenConfiguration
	tokenStore TokenStore
}

// NewTokenProvider creates a new Keycloak token provider
func NewTokenProvider(base *base.BaseService, config domain.TokenConfiguration, tokenStore TokenStore) *TokenProvider {
	return &TokenProvider{
		base:       base,
		config:     config,
		tokenStore: tokenStore,
	}
}

// GenerateWebSocketToken generates a token for WebSocket authentication
func (p *TokenProvider) GenerateWebSocketToken(ctx context.Context, tenantID, branchID, userID string) (string, time.Time, error) {
	// Generate a secure random token
	tokenBytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", time.Time{}, err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Set expiry time based on configuration
	expiry := time.Now().Add(time.Duration(p.config.TokenLifetime) * time.Second)

	// Store the token
	err = p.tokenStore.Store(token, tenantID, branchID, userID, expiry)
	if err != nil {
		return "", time.Time{}, errors.New("failed to store token")
	}

	p.base.Logger.Info(ctx, "Generated WebSocket token", map[string]interface{}{
		"tenantID": tenantID,
		"branchID": branchID,
		"userID":   userID,
		"expiry":   expiry,
	})

	return token, expiry, nil
}

// ValidateToken checks if a token is valid
func (p *TokenProvider) ValidateToken(ctx context.Context, token string) (bool, string, error) {
	tenantID, branchID, userID, expiry, exists := p.tokenStore.Get(token)

	// Check if token exists
	if !exists {
		return false, "", nil
	}
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		return false, "", err
	}
	if tenant.Domain() != tenantID || branch != branchID {
		return false, "", errors.New("unmatch tenat/branch")
	}
	// Check if token has expired
	if time.Now().After(expiry) {
		p.base.Logger.Info(ctx, "Token expired", map[string]interface{}{
			"tenantID": tenantID,
			"branchID": branchID,
			"userID":   userID,
		})
		// Clean up expired token
		_ = p.tokenStore.Delete(token)
		return false, "", nil
	}

	return true, userID, nil
}

// RevokeToken invalidates a token
func (p *TokenProvider) RevokeToken(ctx context.Context, token string) error {
	tenantID, branchID, userID, _, exists := p.tokenStore.Get(token)

	if !exists {
		// Token doesn't exist, nothing to revoke
		return nil
	}

	err := p.tokenStore.Delete(token)
	if err != nil {
		return err
	}

	p.base.Logger.Info(ctx, "Token revoked", map[string]interface{}{
		"tenantID": tenantID,
		"branchID": branchID,
		"userID":   userID,
	})

	return nil
}

// StartCleanupTask starts a background task to clean up expired tokens
func (p *TokenProvider) StartCleanupTask(ctx context.Context, cleanupInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.cleanupExpiredTokens(ctx)
			}
		}
	}()
}

// cleanupExpiredTokens is a helper method for the cleanup task
func (p *TokenProvider) cleanupExpiredTokens(ctx context.Context) {

	p.base.Logger.Info(ctx, "Token cleanup attempted", nil)

	if inMemoryStore, ok := p.tokenStore.(*InMemoryTokenStore); ok {
		inMemoryStore.mu.Lock()
		defer inMemoryStore.mu.Unlock()

		now := time.Now()
		for token, data := range inMemoryStore.store {
			if now.After(data.Expiry) {
				delete(inMemoryStore.store, token)
			}
		}
	}
}
