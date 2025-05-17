package services

import (
	"backend/shared/token/app"
	"backend/shared/token/domain"
	"backend/shared/token/infra"
	"context"
)

func (c *Container) initializeTokenService(s *Services) (err error) {
	// Create token configuration
	tokenConfig := &domain.TokenConfiguration{
		WebSocketAudience: "websocket", // Set appropriate value
		TokenLifetime:     3600,        // Set appropriate value (seconds)
	}

	// Create token factory
	tokenFactory := infra.NewTokenFactory(c.base, *tokenConfig)

	// Create token provider with in-memory storage
	tokenProvider := tokenFactory.CreateInMemoryProvider(context.Background())

	// Create token service using the provider
	s.TokenService = app.New(c.base, &app.ServiceDependencies{
		Provider:     tokenProvider,
		TenatService: s.TenantService,
    })

	return
}
