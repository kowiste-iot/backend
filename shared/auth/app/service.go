package app

import (
	"ddd/pkg/config"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/policy"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
)

type Service struct {
	tenantProvider     auth.TenantProvider
	clientProvider     auth.ClientProvider
	resourceProvider   resource.ResourceProvider
	policyProvider     policy.PolicyProvider
	permissionProvider permission.PermissionProvider
	scopeProvider      scope.ScopeProvider
	tenantConfig       *config.TenantConfiguration
}

func NewAuthService(
	tenantConfig *config.TenantConfiguration,
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
		tenantProvider:     tenantProvider,
		clientProvider:     clientProvider,
		resourceProvider:   resourceProvider,
		policyProvider:     policyProvider,
		permissionProvider: permissionProvider,
		scopeProvider:      scopeProvider,
	}
}
