package app

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/command"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/policy"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
	baseCmd "ddd/shared/base/command"
	"fmt"
	"strings"
)
//TODO:maybe this can be in the tenant feature and remove this auth layer tenant
func (s *Service) CreateTenant(ctx context.Context, input *command.CreateTenantInput) (id string, err error) {
	defer func() {
		if err != nil {
			s.tenantProvider.DeleteTenant(ctx, id)
		}
	}()
	tenant := &auth.Tenant{
		Name:     input.Name,
		Domain:   input.Domain,
		Settings: auth.TenantSettings{},
		Theme: &auth.TenantTheme{
			Login: "custom",
		},
	}

	id, err = s.tenantProvider.CreateTenant(ctx, tenant)
	if err != nil {
		return
	}
	tenant, err = s.tenantProvider.GetTenant(ctx, id)
	if err != nil {
		return
	}
	// Create default roles
	for _, role := range auth.DefaultRoles {
		input := command.CreateRoleInput{
			BaseInput:   baseCmd.NewInput(tenant.Domain, ""),
			Name:        role.Name,
			Description: role.Description,
		}
		_, err = s.tenantProvider.CreateRole(ctx, &input)
		if err != nil {
			return
		}
	}

	err = s.createTenantClients(ctx, id)
	if err != nil {
		return
	}

	return
}

func (s *Service) createTenantClients(ctx context.Context, tenanatDomain string) (err error) {
	//TODO: move struct inside keycloak implementation here still dont need to know about kc implementation
	_, err = s.clientProvider.CreateClient(ctx, tenanatDomain, auth.Client{
		ClientID:                  "vue-client",
		ClientAuthenticatorType:   "client-secret",
		RedirectURIs:              []string{"http://localhost:5173/*"},
		WebOrigins:                []string{"*"},
		StandardFlowEnabled:       true,
		DirectAccessGrantsEnabled: true,
		PublicClient:              true,
		Protocol:                  "openid-connect",
		Attributes: map[string]string{
			"realm_client":                              "false",
			"oidc.ciba.grant.enabled":                   "false",
			"backchannel.logout.session.required":       "true",
			"post.logout.redirect.uris":                 "http://localhost:5173",
			"display.on.consent.screen":                 "false",
			"oauth2.device.authorization.grant.enabled": "false",
			"backchannel.logout.revoke.offline.tokens":  "false",
		},
		DefaultClientScopes: []string{
			"web-origins",
			"acr",
			"roles",
			"profile",
			"basic",
			"email",
		},
		OptionalClientScopes: []string{
			"address",
			"phone",
			"offline_access",
			"organization",
			"microprofile-jwt",
		},
	})
	if err != nil {
		return
	}
	//backend client
	client, err := s.clientProvider.CreateClient(ctx, tenanatDomain, auth.Client{
		ClientID:                  auth.BackendClient,
		Name:                      "Backend Service",
		Description:               "Backend API Service",
		RootURL:                   "http://localhost:8080",
		AdminURL:                  "http://localhost:8080",
		ClientAuthenticatorType:   "client-secret",
		RedirectURIs:              []string{"http://localhost:8080/*"},
		WebOrigins:                []string{"http://localhost:8080"},
		NotBefore:                 0,
		DirectAccessGrantsEnabled: true,
		ServiceAccountsEnabled:    true,
		ServiceAccountEnabled:     true,
		AuthorizationEnabled:      true,
		Protocol:                  "openid-connect",
		Attributes: map[string]string{
			"realm_client":                             "false",
			"backchannel.logout.session.required":      "true",
			"backchannel.logout.revoke.offline.tokens": "false",
		},
		NodeReRegistrationTimeout: -1,
		DefaultClientScopes: []string{
			"web-origins",
			"roles",
			"profile",
			"email",
		},
		OptionalClientScopes: []string{
			"address",
			"phone",
			"offline_access",
			"microprofile-jwt",
		},
	})
	if err != nil {
		return
	}
	err = s.createClientPermissions(ctx, tenanatDomain, client)
	return
}

func (s *Service) createClientPermissions(ctx context.Context, tenanatDomain string, client *auth.Client) (err error) {

	for _, scName := range []string{scope.View, scope.Create, scope.Update, scope.Delete} {
		_, err := s.scopeProvider.CreateScope(ctx, tenanatDomain, *client.ID, scope.Scope{
			Name:        scName,
			DisplayName: scName,
		})
		if err != nil {
			return err
		}

	}

	resources := []struct {
		name   string
		scopes []string
	}{
		{
			name:   resource.Asset,
			scopes: []string{scope.View, scope.Create, scope.Update, scope.Delete},
		},
		{
			name:   resource.User,
			scopes: []string{scope.View, scope.Create, scope.Update, scope.Delete},
		},
	}

	for _, res := range resources {
		createdResource, err := s.resourceProvider.CreateResource(ctx, tenanatDomain, *client.ID, resource.Resource{
			Name:        res.name,
			DisplayName: res.name,
			Scopes:      res.scopes,
		})
		if err != nil {
			return fmt.Errorf("failed to create resource %s: %w", res.name, err)
		}

		adminRoles := []string{auth.RoleAdmin, auth.RoleSuperAdmin}
		for _, scope := range res.scopes {
			perm := permission.Permission{
				Name:             fmt.Sprintf("%s-%s-permission", res.name, scope),
				Description:      fmt.Sprintf("Permission to %s %s", scope, res.name),
				Type:             permission.TypeScope,
				Resources:        []string{createdResource.ID},
				Scopes:           []string{scope},
				DecisionStrategy: permission.DecisionAffirmative,
			}

			createdPerm, err := s.permissionProvider.CreatePermission(ctx, tenanatDomain, *client.ID, perm)
			if err != nil {
				return fmt.Errorf("failed to create permission for %s-%s: %w", res.name, scope, err)
			}

			pol := policy.Policy{
				Name:             fmt.Sprintf("%s-%s-policy", res.name, scope),
				Description:      fmt.Sprintf("Policy for %s to %s %s", strings.Join(adminRoles, " and "), scope, res.name),
				Type:             permission.TypeRole,
				Roles:            adminRoles,
				Logic:            permission.LogicPositive,
				DecisionStrategy: permission.DecisionAffirmative,
			}

			createdPolicy, err := s.policyProvider.CreatePolicy(ctx, tenanatDomain, *client.ID, pol)
			if err != nil {
				return fmt.Errorf("failed to create policy for %s-%s: %w", res.name, scope, err)
			}

			createdPerm.Policies = []string{createdPolicy.ID}
			if err = s.permissionProvider.UpdatePermission(ctx, tenanatDomain, *client.ID, *createdPerm); err != nil {
				return fmt.Errorf("failed to update permission with policy for %s-%s: %w", res.name, scope, err)
			}
		}
	}
	return nil
}


