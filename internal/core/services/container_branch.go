package services

import (
	"errors"

	"backend/internal/features/tenant/app"
	repoTenant "backend/internal/features/tenant/infra/gorm"
	tenantKeycloak "backend/internal/features/tenant/infra/keycloak"
)

func (c *Container) initializeBranchService(s *Services) error {
	if s.ResourceService == nil || s.PermissionService == nil ||
		s.ScopeService == nil || s.RoleService == nil {
		return errors.New("resource, permission, scope and role services must be initialized first")
	}

	branchRepo := repoTenant.NewBranchRepository(c.base.DB)
	branchKC := tenantKeycloak.NewBranch(c.tenantConfig, c.auth)

	s.BranchService = app.NewBranchService(c.base, &app.BranchDependencies{
		Branch:     branchKC,
		Repo:       branchRepo,
		Role:       s.RoleService,
		Scope:      s.ScopeService,
		Resource:   s.ResourceService,
		Permission: s.PermissionService,
		Config:     c.tenantConfig,
	})
	return nil
}
