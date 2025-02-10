package keycloak

import (
	"backend/internal/features/role/domain"
	"backend/internal/features/role/domain/command"
	baseCmd "backend/shared/base/command"
	"backend/shared/util"
	"context"
	"fmt"

	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

type RoleKeycloak struct {
	*keycloak.Keycloak
}

func New(core *keycloak.Keycloak) *RoleKeycloak {
	return &RoleKeycloak{
		Keycloak: core,
	}
}

func (rk RoleKeycloak) CreateRole(ctx context.Context, input *command.CreateRoleInput) (string, error) {

	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	keycloakRole := gocloak.Role{
		Name:        &input.Name,
		Description: &input.Description,
	}
	err = rk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return "", fmt.Errorf("error getting client: %w", err)
	}
	roleID, err := rk.Client.CreateClientRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		keycloakRole,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant role: %w", err)
	}
	role, err := rk.GetRole(ctx, &command.RoleIDInput{
		BaseInput: input.BaseInput,
		RoleID:    roleID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get  role: %w", err)
	}

	rk.createPolicy(ctx, input.TenantDomain, *input.ClientID, policy{
		Name:             fmt.Sprintf("%s-policy", input.Name),
		Description:      fmt.Sprintf("Policy for %s ", util.CapitalizeFirst(input.Name)),
		Type:             TypeRole,
		Roles:            []string{role.ID},
		Logic:            LogicPositive,
		DecisionStrategy: DecisionAffirmative,
	})
	return role.ID, nil
}
func (rk RoleKeycloak) DeleteRole(ctx context.Context, input *command.RoleIDInput) error {
	//dont need delete policy, keycloak automatically delete it
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	err = rk.Client.DeleteClientRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		input.RoleID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete tenant role: %w", err)
	}

	return nil

}
func (rk RoleKeycloak) GetRole(ctx context.Context, input *command.RoleIDInput) (*domain.Role, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	rol, err := rk.Client.GetClientRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		input.RoleID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get client role: %w", err)
	}

	return &domain.Role{
		ID:          *rol.ID,
		Name:        *rol.Name,
		Description: *rol.Description,
	}, nil

}
func (rk RoleKeycloak) GetRoles(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Role, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	roles, err := rk.Client.GetClientRoles(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		gocloak.GetRoleParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant roles: %w", err)
	}

	var authRoles []domain.Role
	for _, rol := range roles {
		if *rol.Name == domain.RoleUma {
			continue //dont show uma role
		}
		r := domain.New(*rol.ID, *rol.Name)
		r.WithDescription(rol.Description)
		authRoles = append(authRoles, *r)
	}

	return authRoles, nil
}
func (rk RoleKeycloak) AssignRoleToResource(ctx context.Context, input *command.ResourceAssignRoleInput) (err error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	fmt.Println(token)
	_, err = rk.getPolicyByName(ctx, &input.BaseInput, policyName(input.RoleName))
	if err != nil {
		return
	}
	scopes := make([]string, 0)
	for i := range input.Scopes {
		scopes = append(scopes, input.Scopes[i].Name)
	}
	// perm := permission.Permission{
	// 	Name:             permission.NameNonAdmin(input.RoleName, input.ResourceName),
	// 	Description:      fmt.Sprintf("Permission for %s resource with %s role", input.ResourceName, input.RoleName),
	// 	Type:             permission.TypeScope,
	// 	Resources:        []string{input.ResourceID},
	// 	Scopes:           scopes,
	// 	Policies:         []string{p.ID},
	// 	DecisionStrategy: permission.DecisionAffirmative,
	// }

	// _, err = rk.CreatePermission(ctx, input.TenantDomain, *input.ClientID, perm)
	// if err != nil {
	// 	return fmt.Errorf("failed to create permission for %s %s: %w", input.ResourceName, input.RoleName, err)
	// }

	return
}
func (rk RoleKeycloak) RemoveRolesFromResource(ctx context.Context, input *command.ResourceAssignRoleInput) (err error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	permissions, err := rk.Client.GetPermissions(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, gocloak.GetPermissionParams{
		Resource: &input.ResourceID,
	})
	for i := range permissions {
		err = rk.Client.DeletePermission(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, *permissions[i].ID)
		if err != nil {
			return
		}
	}
	return

}
