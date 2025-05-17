package app

import (
	"backend/internal/features/permission/domain"
	"backend/internal/features/permission/domain/command"
	appRole "backend/internal/features/user/app"
	roleDomain "backend/internal/features/user/domain"
	appScope "backend/internal/features/scope/app"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type PermissionService interface {
	CreatePermission(ctx context.Context, input *command.CreatePermissionInput) (*domain.Permission, error)
	UpdatePermission(ctx context.Context, input *command.UpdatePermissionInput) (*domain.Permission, error)
	ListPermissions(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Permission, error)
}
type Config struct {
	DefaultRoles []string
}

type ServiceDependencies struct {
	Repo   domain.PermissionProvider
	Role   appRole.RoleService
	Scope  appScope.ScopeService
	Config *Config
}
type permissionService struct {
	permissionProvider domain.PermissionProvider
	scopeProvider      appScope.ScopeService
	roles              appRole.RoleService
	config             *Config
	*base.BaseService
}

func NewService(base *base.BaseService, deps *ServiceDependencies) PermissionService {
	return &permissionService{
		permissionProvider: deps.Repo,
		scopeProvider:      deps.Scope,
		roles:              deps.Role,
		BaseService:        base,
		config:             deps.Config,
	}
}
func (s *permissionService) CreatePermission(ctx context.Context, input *command.CreatePermissionInput) (permission *domain.Permission, err error) {
	// err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  domain.ResourceR,
	// 	Scope:     scope.Create,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	p, err := domain.New(input.Name, input.Description, input.Type, input.DecisionStrategy, input.Resources, input.Scopes, input.Policies)
	if err != nil {
		return
	}
	if input.ResourceType != "" {
		p.SetResourceType(input.ResourceType)
	}
	scopes, err := s.scopeProvider.ListScopes(ctx, &input.BaseInput)
	if err != nil {
		return
	}
	permission, err = s.permissionProvider.CreatePermission(ctx, scopes, &input.BaseInput, p)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return
}

func (s *permissionService) ListPermissions(ctx context.Context, input *baseCmd.BaseInput) (permissions []domain.Permission, err error) {
	// err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: *input,
	// 	Resource:  domain.ResourceR,
	// 	Scope:     scope.View,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	permissions, err = s.permissionProvider.ListPermissions(ctx, input)
	if err != nil {
		return nil, err
	}

	return
}
func (s *permissionService) UpdatePermission(ctx context.Context, input *command.UpdatePermissionInput) (permission *domain.Permission, err error) {
	// err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  domain.ResourceR,
	// 	Scope:     scope.Create,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	//Delete previous permissions
	err = s.permissionProvider.DeletePermission(ctx, &input.BaseInput, input.ResourceID)
	if err != nil {
		return nil, err
	}

	roles, err := s.roles.ListRoles(ctx, &input.BaseInput)
	if err != nil {
		return
	}
	rolesMap := make(map[string]roleDomain.Role)
	for i := range roles {
		rolesMap[roles[i].Name] = roles[i]
	}
	for roleName, scopes := range input.Roles {
		scopesIDs := make([]string, 0)
		for j := range scopes {
			scopesIDs = append(scopesIDs, scopes[j].ID)
		}
		role, found := rolesMap[roleName]
		if !found {
			return nil, fmt.Errorf("role not found")
		}
		_, err = s.CreatePermission(ctx, &command.CreatePermissionInput{
			BaseInput:        input.BaseInput,
			Name:             domain.NameNonAdmin(roleName, input.ResourceName),
			Description:      fmt.Sprintf("Permission for %s resource with %s role", input.ResourceName, roleName),
			Type:             domain.TypeScope,
			Resources:        input.ResourceID,
			Scopes:           scopesIDs,
			Policies:         []string{role.PolicyID},
			DecisionStrategy: domain.DecisionAffirmative,
		})
		if err != nil {
			return nil, err
		}
	}

	return
}
