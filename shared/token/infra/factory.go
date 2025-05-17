package infra

import (
	"backend/shared/base"
	"backend/shared/token/domain"
	inmemmo "backend/shared/token/infra/in_memory"
	"context"
	"time"
)

// Factory creates and configures token-related components
type TokenFactory struct {
	Base   *base.BaseService
	Config domain.TokenConfiguration
}

// NewTokenFactory creates a new token factory
func NewTokenFactory(base *base.BaseService, config domain.TokenConfiguration) *TokenFactory {
	return &TokenFactory{
		Base:   base,
		Config: config,
	}
}

// CreateInMemoryProvider creates a token provider with in-memory storage
func (f *TokenFactory) CreateInMemoryProvider(ctx context.Context) *inmemmo.TokenProvider {
	store := inmemmo.NewInMemoryTokenStore()
	provider := inmemmo.NewTokenProvider(f.Base, f.Config, store)

	// Start a cleanup task that runs every hour
	provider.StartCleanupTask(ctx, 1*time.Hour)

	f.Base.Logger.Info(ctx, "Created token provider with in-memory storage", nil)
	return provider
}

// CreateWithCustomStore creates a token provider with a custom storage implementation
func (f *TokenFactory) CreateWithCustomStore(ctx context.Context, store inmemmo.TokenStore) *inmemmo.TokenProvider {
	provider := inmemmo.NewTokenProvider(f.Base, f.Config, store)

	// Start a cleanup task that runs every hour
	provider.StartCleanupTask(ctx, 1*time.Hour)

	f.Base.Logger.Info(ctx, "Created token provider with custom storage", nil)
	return provider
}
