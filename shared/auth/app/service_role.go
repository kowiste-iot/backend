package app

import (
	"context"
	"ddd/shared/auth/domain/command"
	"ddd/shared/auth/domain/permission"
	"ddd/shared/auth/domain/policy"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/role"
	"ddd/shared/auth/domain/scope"
	baseCmd "ddd/shared/base/command"
	"ddd/shared/util"
	"fmt"
	"slices"
)

// GetTenantRoles retrieves all roles for a tenant
func (s *Service) GetRoles(ctx context.Context, input *baseCmd.BaseInput) ([]role.Role, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resource.Role,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	return s.tenantProvider.GetRoles(ctx, input)
}

// GetTenantRole retrieves a specific role from a tenant
func (s *Service) GetRole(ctx context.Context, input *command.RoleIDInput) (*role.Role, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Role,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	return s.tenantProvider.GetRole(ctx, input)
}

// CreateRole creates a new role for a tenant
func (s *Service) CreateRole(ctx context.Context, input *command.CreateRoleInput) (id string, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Role,
		Scope:     scope.Create,
	})
	if err != nil {
		return
	}
	if s.isDefaultRole(input.Name) {
		return "", fmt.Errorf("cannot create role with reserved name: %s", input.Name)
	}
	id, err = s.tenantProvider.CreateRole(ctx, input)
	if err != nil {
		return
	}
	pol := policy.Policy{
		Name:             fmt.Sprintf("%s-policy", input.Name),
		Description:      fmt.Sprintf("Policy for %s ", util.CapitalizeFirst(input.Name)),
		Type:             policy.TypeRole,
		Roles:            []string{id},
		Logic:            permission.LogicPositive,
		DecisionStrategy: permission.DecisionAffirmative,
	}
	client, err := s.clientProvider.GetClientByClientID(ctx, input.TenantDomain, command.ClientName(input.BranchName))
	if err != nil {
		return "", fmt.Errorf("error getting client: %w", err)
	}
	_, err = s.policyProvider.CreatePolicy(ctx, input.TenantDomain, *client.ID, pol)
	if err != nil {
		return "", fmt.Errorf("failed to create policy for %s: %w", input.Name, err)
	}
	return
}

// DeleteRole deletes a role from a tenant
func (s *Service) DeleteRole(ctx context.Context, input *command.RoleIDInput) (err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Role,
		Scope:     scope.Delete,
	})
	if err != nil {
		return err
	}
	if s.isDefaultRole(input.RoleID) {
		return fmt.Errorf("cannot delete default role: %s", input.RoleID)
	}
	inputPolicy := command.PolicyNameInput{
		BaseInput: input.BaseInput,
		PolicyName:  command.PolicyName(input.RoleID),
	}
	err = s.policyProvider.DeletePolicy(ctx, &inputPolicy)
	if err != nil {
		return fmt.Errorf("failed to delete policy %w", err)
	}
	err = s.tenantProvider.DeleteRole(ctx, input)
	if err != nil {
		return err
	}
	return
}

// AssignRoles assigns roles to a user
func (s *Service) AssignRoles(ctx context.Context, input *command.AssignRolesInput) error {
	// err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  resource.Role,
	// 	Scope:     scope.Update,
	// })
	// if err != nil {
	// 	return err
	// }
	return s.tenantProvider.AssignRoles(ctx, input)
}

// RemoveRoles removes roles from a user
func (s *Service) RemoveRoles(ctx context.Context, input *command.RemoveRolesInput) error {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Role,
		Scope:     scope.Update,
	})
	if err != nil {
		return err
	}
	return s.tenantProvider.RemoveRoles(ctx, input)
}

// GetUserRoles gets all roles assigned to a user
func (s *Service) GetUserRoles(ctx context.Context, input *command.UserRolesInput) ([]role.Role, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Role,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	return s.tenantProvider.GetUserRoles(ctx, input)
}

// isDefaultRole checks if a role name or ID matches any default roles
func (s *Service) isDefaultRole(identifier string) bool {
	return slices.ContainsFunc(role.AllRoles(s.tenantConfig.Authorization.Roles), func(role role.Role) bool {
		return role.Name == identifier
	})
}
