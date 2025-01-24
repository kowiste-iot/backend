package app

import (
	"context"
	"ddd/shared/auth/domain/scope"
	baseCmd "ddd/shared/base/command"
)

// GetTenantRoles retrieves all roles for a tenant
func (s *Service) GetScopes(ctx context.Context, input *baseCmd.BaseInput) ([]scope.Scope, error) {
	// err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: *input,
	// 	Resource:  resource.Role,
	// 	Scope:     scope.View,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	return s.scopeProvider.ListScopes(ctx, input)
}
