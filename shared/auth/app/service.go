package app

import (
	"backend/pkg/config"
	auth "backend/shared/auth/domain"
	"backend/shared/auth/domain/permission"
	"backend/shared/auth/domain/policy"
	"backend/shared/auth/domain/resource"
	"backend/shared/auth/domain/role"
	"backend/shared/auth/domain/scope"
	"backend/shared/base"
)

type Service struct {
	tenantProvider     auth.TenantProvider
	roleProvider       role.RoleProvider
	clientProvider     auth.ClientProvider
	resourceProvider   resource.ResourceProvider
	policyProvider     policy.PolicyProvider
	permissionProvider permission.PermissionProvider
	scopeProvider      scope.ScopeProvider
	tenantConfig       *config.TenantConfiguration
	*base.BaseService
}

func NewAuthService(
	tenantConfig *config.TenantConfiguration,
	base *base.BaseService,
	identityProvider auth.IdentityProvider,
	tenantProvider auth.TenantProvider,
	roleProvider role.RoleProvider,
	clientProvider auth.ClientProvider,
	resourceProvider resource.ResourceProvider,
	policyProvider policy.PolicyProvider,
	permissionProvider permission.PermissionProvider,
	scopeProvider scope.ScopeProvider,
) *Service {

	return &Service{
		tenantConfig:       tenantConfig,
		BaseService:        base,
		tenantProvider:     tenantProvider,
		roleProvider:       roleProvider,
		clientProvider:     clientProvider,
		resourceProvider:   resourceProvider,
		policyProvider:     policyProvider,
		permissionProvider: permissionProvider,
		scopeProvider:      scopeProvider,
	}
}
