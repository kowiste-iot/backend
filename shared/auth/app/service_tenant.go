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
	"ddd/shared/util"
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
	_, err = s.clientProvider.CreateClient(ctx, input.TenantDomain, auth.Client{
		ClientID:            s.tenantConfig.WebClient.ClientID,
		Name:                s.tenantConfig.WebClient.Name,
		RedirectURIs:        s.tenantConfig.WebClient.RedirectURIs,
		WebOrigins:          s.tenantConfig.WebClient.Origins,
		StandardFlowEnabled: true,
		PublicClient:        true,
	})
	if err != nil {
		return
	}
	//backend client
	upperBranch := util.CapitalizeFirst(input.BranchName)
	client, err := s.clientProvider.CreateClient(ctx, input.TenantDomain, auth.Client{
		ClientID:                command.ClientName(input.BranchName),
		Name:                    upperBranch + s.tenantConfig.BackendClient.Name,
		Description:             upperBranch + s.tenantConfig.BackendClient.Description,
		RootURL:                 *s.tenantConfig.BackendClient.RootURL,
		AdminURL:                *s.tenantConfig.BackendClient.AdminURL,
		ClientAuthenticatorType: "client-secret",
		RedirectURIs:            s.tenantConfig.BackendClient.RedirectURIs,
		WebOrigins:              s.tenantConfig.BackendClient.Origins,
		ServiceAccountEnabled:   true,
		AuthorizationEnabled:    true,
		Authorization:           true,
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
	for _, role := range auth.AllRoles(s.tenantConfig.Authorization.Roles) {
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

	for _, scName := range scope.AllScopes() {
		_, err := s.scopeProvider.CreateScope(ctx, input.TenantDomain, *client.ID, scope.Scope{
			Name:        scName,
			DisplayName: util.CapitalizeFirst(scName),
		})
		if err != nil {
			return err
		}

	}

	//Create policy 1 for each role
	policies := make(map[string]*policy.Policy)
	for _, role := range auth.AllRoles(s.tenantConfig.Authorization.Roles) {
		r, err := s.tenantProvider.GetRole(ctx, &command.RoleIDInput{
			BaseInput: baseCmd.NewInput(input.TenantDomain, input.BranchName),
			RoleID:    role.Name,
		})
		if err != nil {
			return fmt.Errorf("failed to get role: %w", err)
		}
		pol := policy.Policy{
			Name:             fmt.Sprintf("%s-policy", role.Name),
			Description:      fmt.Sprintf("Policy for %s ", util.CapitalizeFirst(role.Name)),
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
	resources := resource.EndpointsResources(s.tenantConfig.Authorization.Resources)

	//Create resources and non admin permission
	for _, res := range resources {

		createdResource, err := s.resourceProvider.CreateResource(ctx, input.TenantDomain, *client.ID, resource.Resource{
			Name:        command.ResourceName(res.Name),
			DisplayName: res.Name,
			Type:        res.Type,
			Scopes:      res.Scopes,
		})
		if err != nil {
			return fmt.Errorf("failed to create resource %s: %w", res.Name, err)
		}
		resourceConfig, exist := s.tenantConfig.Authorization.Resources[command.ResourceName(res.Name)]
		if !exist {
			return fmt.Errorf("error fetching resource %s", res.Name)
		}

		for roleName, scopes := range resourceConfig.Permissions {

			p, exist := policies[roleName]
			if !exist {
				return fmt.Errorf("error fetching policy for  %s", roleName)
			}

			perm := permission.Permission{
				Name:             fmt.Sprintf("%s-%s-permission", roleName, res.Name),
				Description:      fmt.Sprintf("Permission for %s resource with %s role", res.Name, roleName),
				Type:             permission.TypeScope,
				Resources:        []string{createdResource.ID},
				Scopes:           scopes,
				Policies:         []string{p.ID},
				DecisionStrategy: permission.DecisionAffirmative,
			}

			_, err = s.permissionProvider.CreatePermission(ctx, input.TenantDomain, *client.ID, perm)
			if err != nil {
				return fmt.Errorf("failed to create permission for %s %s: %w", res.Name, roleName, err)
			}
		}

	}
	//resource permission for admin
	sc := scope.AllScopes()

	p := policies[auth.RoleAdmin]
	perm := permission.Permission{
		Name:             fmt.Sprintf("%s-permission", auth.RoleAdmin),
		Description:      fmt.Sprintf("Permission for %s resource with %s role", auth.RoleAdmin, auth.RoleAdmin),
		Type:             permission.TypeResource,
		ResourceType:     s.tenantConfig.Authorization.AdminGroup,
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
