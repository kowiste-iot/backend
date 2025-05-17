package services

import (
	appScope "backend/internal/features/scope/app"
	scopeKeycloak "backend/internal/features/scope/infra/keycloak"
)

func (c *Container) initializeScopeService(s *Services) error {
	scopeKC := scopeKeycloak.New(c.auth)
	s.ScopeService = appScope.NewService(c.base, &appScope.ServiceDependencies{
		Repo: scopeKC,
	})
	return nil
}
