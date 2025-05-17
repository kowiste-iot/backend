package services

import (
	"backend/internal/features/permission/app"
	permissionKeycloak "backend/internal/features/permission/infra/keycloak"
	"errors"
)

func (c *Container) initializePermissionService(s *Services) error {
	if s.ScopeService == nil || s.RoleService == nil {
		return errors.New("scope and role services must be initialized first")
	}

	permissionKC := permissionKeycloak.New(c.auth)
	s.PermissionService = app.NewService(c.base, &app.ServiceDependencies{
		Repo:  permissionKC,
		Scope: s.ScopeService,
		Role:  s.RoleService,
		Config: &app.Config{
			DefaultRoles: c.tenantConfig.Authorization.Roles,
		},
	})
	return nil
}
