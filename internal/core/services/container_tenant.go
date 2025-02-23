package services

import (
	"errors"

	"backend/internal/features/tenant/app"
	repoTenant "backend/internal/features/tenant/infra/gorm"
	tenantKeycloak "backend/internal/features/tenant/infra/keycloak"
)

func (c *Container) initializeTenantService(s *Services) error {
	if s.BranchService == nil || s.UserService == nil {
		return errors.New("branch and user services must be initialized first")
	}

	tenantRepo := repoTenant.NewTenantRepository(c.base.DB)
	tenantKC := tenantKeycloak.New(c.tenantConfig, c.auth)

	s.TenantService = app.NewTenantService(c.base, &app.ServiceDependencies{
		Branch: s.BranchService,
		Tenant: tenantKC,
		Repo:   tenantRepo,
		User:   s.UserService,
	})
	return nil
}
