package keycloak

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/command"
	baseCmd "ddd/shared/base/command"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

// CreatetRole creates a new role in the specified tenant in the client branchName-service
func (ks *KeycloakService) CreateRole(ctx context.Context, input *command.CreateRoleInput) (string, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	keycloakRole := gocloak.Role{
		Name:        &input.Name,
		Description: &input.Description,
	}
	client, err := ks.GetClientByClientID(ctx, input.TenantDomain, command.ClientName(input.BranchName))
	if err != nil {
		return "", fmt.Errorf("error getting client: %w", err)
	}
	roleID, err := ks.client.CreateClientRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*client.ID,
		keycloakRole,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant role: %w", err)
	}

	return roleID, nil
}

// DeleteRole deletes a role from the specified client
func (ks *KeycloakService) DeleteRole(ctx context.Context, input *command.RoleIDInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeleteClientRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		command.ClientName(input.BranchName),
		input.RoleID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete tenant role: %w", err)
	}

	return nil
}

// GetTenantRole gets a specific role from the tenant
func (ks *KeycloakService) GetRole(ctx context.Context, input *command.RoleIDInput) (*auth.Role, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	client, err := ks.GetClientByClientID(ctx, input.TenantDomain, command.ClientName(input.BranchName))
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	role, err := ks.client.GetClientRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*client.ID,
		input.RoleID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get client role: %w", err)
	}

	return &auth.Role{
		ID:          *role.ID,
		Name:        *role.Name,
		Description: *role.Description,
	}, nil
}

// GetTenantRoles gets all roles from the tenant
func (ks *KeycloakService) GetRoles(ctx context.Context, input *baseCmd.BaseInput) ([]auth.Role, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	roles, err := ks.client.GetClientRoles(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		command.ClientName(input.BranchName),
		gocloak.GetRoleParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant roles: %w", err)
	}

	var authRoles []auth.Role
	for _, role := range roles {
		authRoles = append(authRoles, auth.Role{
			ID:          *role.ID,
			Name:        *role.Name,
			Description: *role.Description,
		})
	}

	return authRoles, nil
}

func (ks *KeycloakService) GetUserRoles(ctx context.Context, input *command.UserRolesInput) (roles []auth.Role, err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	realmRoles, err := ks.client.GetClientRolesByUserID(ctx, token.AccessToken,
		input.TenantDomain, command.ClientName(input.BranchName), input.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting realm roles: %v", err)
	}
	roles = make([]auth.Role, 0)
	for i := range realmRoles {
		roles = append(roles, auth.Role{
			Name:        *realmRoles[i].Name,
			Description: *realmRoles[i].Description,
		})
	}
	return
}
func (ks *KeycloakService) AssignRoles(ctx context.Context, input *command.AssignRolesInput) (err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	client, err := ks.GetClientByClientID(ctx, input.TenantDomain, command.ClientName(input.BranchName))
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	addRoles := make([]gocloak.Role, 0)
	for i := range input.Roles {
		role, err := ks.client.GetClientRole(ctx, token.AccessToken, input.TenantDomain, *client.ID, input.Roles[i])
		if err != nil {
			return fmt.Errorf("error getting realm role: %v", err)
		}
		addRoles = append(addRoles, *role)
	}
	err = ks.client.AddClientRolesToUser(ctx, token.AccessToken,
		input.TenantDomain, *client.ID, input.UserID, addRoles)

	return
}
func (ks *KeycloakService) RemoveRoles(ctx context.Context, input *command.RemoveRolesInput) (err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	removeRoles := make([]gocloak.Role, 0)
	for i := range input.Roles {
		role, err := ks.client.GetRealmRole(ctx, token.AccessToken, input.TenantDomain, input.Roles[i])
		if err != nil {
			return fmt.Errorf("error getting realm role: %v", err)
		}
		removeRoles = append(removeRoles, *role)
	}
	err = ks.client.DeleteClientRolesFromUser(ctx, token.AccessToken,
		input.TenantDomain, command.ClientName(input.BranchName), input.UserID, removeRoles)
	return
}
