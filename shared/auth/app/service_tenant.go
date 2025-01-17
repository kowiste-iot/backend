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
)

func (s *Service) CreateTenant(ctx context.Context, input *command.CreateTenantInput) (tenant *auth.Tenant, err error) {
	defer func() {
		if err != nil {
			s.tenantProvider.DeleteTenant(ctx, tenant.ID)
		}
	}()
	tenant = &auth.Tenant{
		Name:     input.Name,
		Domain:   input.Domain,
		Settings: auth.TenantSettings{},
		Theme: &auth.TenantTheme{
			Login: "custom",
		},
	}

	tenant.Domain, err = s.tenantProvider.CreateTenant(ctx, tenant)
	if err != nil {
		return
	}
	tenant, err = s.tenantProvider.GetTenant(ctx, tenant.Domain)
	if err != nil {
		return
	}

	i := baseCmd.NewInput(tenant.Domain, input.DefaultBranch)
	err = s.createTenantClients(ctx, &i)
	if err != nil {
		return
	}

	return
}

func (s *Service) createTenantClients(ctx context.Context, input *baseCmd.BaseInput) (err error) {
	//TODO: move struct inside keycloak implementation here still dont need to know about kc implementation
	_, err = s.clientProvider.CreateClient(ctx, input.TenantDomain, auth.Client{
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
	client, err := s.clientProvider.CreateClient(ctx, input.TenantDomain, auth.Client{
		ClientID:                  command.ClientName(input.BranchName),
		Name:                      input.BranchName + " Service",
		Description:               input.BranchName + " Backend API Service",
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
		Authorization: true,
	})
	if err != nil {
		return
	}
	err = s.createRoles(ctx, input)
	if err != nil {
		return
	}

	err = s.createClientPermissions(ctx, input, client)
	return
}
func (s Service) createRoles(ctx context.Context, input *baseCmd.BaseInput) (err error) {
	// Create default roles
	for _, role := range auth.AllRoles() {
		input := command.CreateRoleInput{
			BaseInput:   *input,
			Name:        role.Name,
			Description: role.Description,
		}
		_, err = s.tenantProvider.CreateRole(ctx, &input)
		if err != nil {
			return
		}
	}
	return
}

func (s *Service) createClientPermissions(ctx context.Context, input *baseCmd.BaseInput, client *auth.Client) (err error) {

	for _, scName := range []string{scope.View, scope.Create, scope.Update, scope.Delete} {
		_, err := s.scopeProvider.CreateScope(ctx, input.TenantDomain, *client.ID, scope.Scope{
			Name:        scName,
			DisplayName: scName,
		})
		if err != nil {
			return err
		}

	}

	//Create policy 1 for each role
	policies := make(map[string]*policy.Policy)
	for _, role := range auth.AllRoles() {
		r, err := s.tenantProvider.GetRole(ctx, &command.RoleIDInput{
			BaseInput: baseCmd.NewInput(input.TenantDomain, input.BranchName),
			RoleID:    role.Name,
		})
		if err != nil {
			return fmt.Errorf("failed to get role: %w", err)
		}
		pol := policy.Policy{
			Name:             fmt.Sprintf("%s-policy", role.Name),
			Description:      fmt.Sprintf("Policy for %s ", role.Name),
			Type:             policy.TypeRole,
			Roles:            []string{r.ID},
			Logic:            permission.LogicPositive,
			DecisionStrategy: permission.DecisionAffirmative,
		}

		createdPolicy, err := s.policyProvider.CreatePolicy(ctx, input.TenantDomain, *client.ID, pol)
		if err != nil {
			return fmt.Errorf("failed to create policy for %s: %w", role.Name, err)
		}
		policies[role.Name] = createdPolicy
	}
	resources := resource.EndpointsResources()

	//Create reosources and non admin permission
	for _, res := range resources {

		createdResource, err := s.resourceProvider.CreateResource(ctx, input.TenantDomain, *client.ID, resource.Resource{
			Name:        res.Name,
			DisplayName: res.Name,
			Type:        res.Type,
			Scopes:      res.Scopes,
		})
		if err != nil {
			return fmt.Errorf("failed to create resource %s: %w", res.Name, err)
		}
		for _, role := range auth.NonAdminRoles() {
			sc := []string{scope.View}
			p := policies[role.Name]

			perm := permission.Permission{
				Name:             fmt.Sprintf("%s-%s-permission", role.Name, res.Name),
				Description:      fmt.Sprintf("Permission for %s resource with %s role", res.Name, role.Name),
				Type:             permission.TypeScope,
				Resources:        []string{createdResource.ID},
				Scopes:           sc,
				Policies:         []string{p.ID},
				DecisionStrategy: permission.DecisionAffirmative,
			}

			_, err = s.permissionProvider.CreatePermission(ctx, input.TenantDomain, *client.ID, perm)
			if err != nil {
				return fmt.Errorf("failed to create permission for %s %s: %w", res.Name, role.Name, err)
			}
		}

	}
	//resource permission for admin
	sc := []string{scope.View, scope.Create, scope.Update, scope.Delete}

	p := policies[auth.RoleAdmin]
	perm := permission.Permission{
		Name:             fmt.Sprintf("%s-permission", auth.RoleAdmin),
		Description:      fmt.Sprintf("Permission for %s resource with %s role", auth.RoleAdmin, auth.RoleAdmin),
		Type:             permission.TypeResource,
		ResourceType:     resource.TypeBase,
		Scopes:           sc,
		Policies:         []string{p.ID},
		DecisionStrategy: permission.DecisionAffirmative,
	}

	_, err = s.permissionProvider.CreatePermission(ctx, input.TenantDomain, *client.ID, perm)
	if err != nil {
		return fmt.Errorf("failed to create permission for %s: %w", auth.RoleAdmin, err)
	}
	return nil
}
