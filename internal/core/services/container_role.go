package services

import (
	appRole "backend/internal/features/user/app"
	roleKeycloak "backend/internal/features/role/infra/keycloak"
)

func (c *Container) initializeRoleService(s *Services) error {
	roleKC := roleKeycloak.New(c.auth)
	s.RoleService = appRole.NewRoleService(c.base, roleKC, appRole.Config{
		DefaultRoles: c.tenantConfig.Authorization.Roles,
	})
	return nil
}
