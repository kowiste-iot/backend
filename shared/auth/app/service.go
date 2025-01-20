package app

import (
	"ddd/pkg/config"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/policy"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
	"ddd/shared/base"
)

type Service struct {
	tenantProvider     auth.TenantProvider
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
		clientProvider:     clientProvider,
		resourceProvider:   resourceProvider,
		policyProvider:     policyProvider,
		permissionProvider: permissionProvider,
		scopeProvider:      scopeProvider,
	}
}
