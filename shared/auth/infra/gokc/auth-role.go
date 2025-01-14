package keycloak

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/command"
	baseCmd "ddd/shared/base/command"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

// CreateTenantRole creates a new role in the specified tenant
func (ks *KeycloakService) CreateRole(ctx context.Context, input *command.CreateRoleInput) (string, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	keycloakRole := gocloak.Role{
		Name:        &input.Name,
		Description: &input.Description,
	}

	// if input.Attributes != nil {
	// 	keycloakRole.Attributes = &role.Attributes
	// }

	roleID, err := ks.client.CreateRealmRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		keycloakRole,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant role: %w", err)
	}

	return roleID, nil
}

// UpdateTenantRole updates an existing role in the specified tenant
func (ks *KeycloakService) UpdateRole(ctx context.Context, input *command.UpdateRoleInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	keycloakRole := gocloak.Role{
		Name:        &input.Name,
		Description: &input.Description,
	}

	// if role.Attributes != nil {
	// 	keycloakRole.Attributes = &role.Attributes
	// }

	err = ks.client.UpdateRealmRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.Name,
		keycloakRole,
	)
	if err != nil {
		return fmt.Errorf("failed to update tenant role: %w", err)
	}

	return nil
}

// DeleteTenantRole deletes a role from the specified tenant
func (ks *KeycloakService) DeleteRole(ctx context.Context, input *command.RoleIDInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeleteRealmRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.BranchName,
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
	role, err := ks.client.GetRealmRole(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.RoleID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant role: %w", err)
	}

	return &auth.Role{
		Name:        *role.Name,
		Description: *role.Description,
		// Attributes:  *role.Attributes,
	}, nil
}

// GetTenantRoles gets all roles from the tenant
func (ks *KeycloakService) GetRoles(ctx context.Context, input *baseCmd.BaseInput) ([]auth.Role, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	roles, err := ks.client.GetRealmRoles(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		gocloak.GetRoleParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant roles: %w", err)
	}

	var authRoles []auth.Role
	for _, role := range roles {
		authRoles = append(authRoles, auth.Role{
			Name:        *role.Name,
			Description: *role.Description,
			// Attributes:  *role.Attributes,
		})
	}

	return authRoles, nil
}

func (ks *KeycloakService) GetUserRoles(ctx context.Context, input *command.UserRolesInput) (roles []auth.Role, err error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	realmRoles, err := ks.client.GetRealmRolesByUserID(ctx, token.AccessToken, input.TenantDomain, input.UserID)
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
	addRoles := make([]gocloak.Role, 0)
	for i := range input.Roles {
		role, err := ks.client.GetRealmRole(ctx, token.AccessToken, input.TenantDomain, input.Roles[i])
		if err != nil {
			return fmt.Errorf("error getting realm role: %v", err)
		}
		addRoles = append(addRoles, *role)
	}
	err = ks.client.AddRealmRoleToUser(ctx, token.AccessToken, input.TenantDomain, input.UserID, addRoles)

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
	err = ks.client.DeleteRealmRoleFromUser(ctx, token.AccessToken, input.TenantDomain, input.UserID, removeRoles)
	return
}
