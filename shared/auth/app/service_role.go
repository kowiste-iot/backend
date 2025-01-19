package app

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/command"
	baseCmd "ddd/shared/base/command"
	"fmt"
	"slices"
)

// GetTenantRoles retrieves all roles for a tenant
func (s *Service) GetRoles(ctx context.Context, input *baseCmd.BaseInput) ([]auth.Role, error) {
	return s.tenantProvider.GetRoles(ctx, input)
}

// GetTenantRole retrieves a specific role from a tenant
func (s *Service) GetRole(ctx context.Context, input *command.RoleIDInput) (*auth.Role, error) {
	return s.tenantProvider.GetRole(ctx, input)
}

// CreateRole creates a new role for a tenant
func (s *Service) CreateRole(ctx context.Context, input *command.CreateRoleInput) (string, error) {
	if s.isDefaultRole(input.Name) {
		return "", fmt.Errorf("cannot create role with reserved name: %s", input.Name)
	}
	return s.tenantProvider.CreateRole(ctx, input)
}

// DeleteRole deletes a role from a tenant
func (s *Service) DeleteRole(ctx context.Context, input *command.RoleIDInput) error {
	if s.isDefaultRole(input.RoleID) {
		return fmt.Errorf("cannot delete default role: %s", input.RoleID)
	}
	return s.tenantProvider.DeleteRole(ctx, input)
}

// AssignRoles assigns roles to a user
func (s *Service) AssignRoles(ctx context.Context, input *command.AssignRolesInput) error {
	return s.tenantProvider.AssignRoles(ctx, input)
}

// RemoveRoles removes roles from a user
func (s *Service) RemoveRoles(ctx context.Context, input *command.RemoveRolesInput) error {
	return s.tenantProvider.RemoveRoles(ctx, input)
}

// GetUserRoles gets all roles assigned to a user
func (s *Service) GetUserRoles(ctx context.Context, input *command.UserRolesInput) ([]auth.Role, error) {
	return s.tenantProvider.GetUserRoles(ctx, input)
}

// isDefaultRole checks if a role name or ID matches any default roles
func (s *Service) isDefaultRole(identifier string) bool {
	return slices.ContainsFunc(auth.AllRoles(s.tenantConfig.Authorization.Roles), func(role auth.Role) bool {
		return role.Name == identifier
	})
}
