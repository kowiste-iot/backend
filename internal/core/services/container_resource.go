package services

import (
	"errors"

	"backend/internal/features/resource/app"
	resourceKeycloak "backend/internal/features/resource/infra/keycloak"
)

func (c *Container) initializeResourceService(s *Services) error {
	if s.PermissionService == nil || s.ScopeService == nil || s.RoleService == nil {
		return errors.New("permission, scope and role services must be initialized first")
	}

	resourceKC := resourceKeycloak.New(c.auth)
	s.ResourceService = app.NewService(c.base, &app.ServiceDependencies{
		Repo:       resourceKC,
		Roles:      s.RoleService,
		Permission: s.PermissionService,
		Scopes:     s.ScopeService,
		Config: &app.Config{
			DefaultRoles: c.tenantConfig.Authorization.Roles,
		},
	})
	return nil
}
