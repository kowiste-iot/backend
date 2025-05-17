package keycloak

import (
	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
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

func (rk RoleKeycloak) CreateRole(ctx context.Context, input *command.CreateRoleInput) (role *domain.Role, err error) {

	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	keycloakRole := gocloak.Role{
		Name:        &input.Name,
		Description: &input.Description,
	}
	err = rk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	roleID, err := rk.Client.CreateClientRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*input.ClientID,
		keycloakRole,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant role: %w", err)
	}
	role, err = rk.GetRole(ctx, &command.RoleIDInput{
		BaseInput: input.BaseInput,
		RoleID:    roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get  role: %w", err)
	}

	policy, err := rk.createPolicy(ctx, input.TenantDomain, *input.ClientID, policy{
		Name:             fmt.Sprintf("%s-policy", input.Name),
		Description:      fmt.Sprintf("Policy for %s ", util.CapitalizeFirst(input.Name)),
		Type:             domain.TypeRole,
		Roles:            []string{role.ID},
		Logic:            domain.LogicPositive,
		DecisionStrategy: domain.DecisionAffirmative,
	})
	if err != nil {
		return
	}
	return role.WithPolicy(policy.ID), nil
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
		pol, err := rk.getPolicyByName(ctx, input, domain.PolicyName(*rol.Name))
		if err != nil {
			return nil, err
		}
		r := domain.NewRole(*rol.ID, *rol.Name)
		r.WithDescription(rol.Description)
		authRoles = append(authRoles, *r.WithPolicy(pol.ID))

	}
	return authRoles, nil
}
