package app

import (
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
}

func NewAuthService(
	identityProvider auth.IdentityProvider,
	tenantProvider auth.TenantProvider,
	clientProvider auth.ClientProvider,
	resourceProvider resource.ResourceProvider,
	policyProvider policy.PolicyProvider,
	permissionProvider permission.PermissionProvider,
	scopeProvider scope.ScopeProvider,
) *Service {
	return &Service{
		tenantProvider:     tenantProvider,
		clientProvider:     clientProvider,
		resourceProvider:   resourceProvider,
		policyProvider:     policyProvider,
		permissionProvider: permissionProvider,
		scopeProvider:      scopeProvider,
	}
}
