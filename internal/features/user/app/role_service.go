package app

import (
	resourceDomain "backend/internal/features/resource/domain"
	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
	scopeDomain "backend/internal/features/scope/domain"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
	"slices"
)

type RoleService interface {
	CreateRole(ctx context.Context, input *command.CreateRoleInput) (*domain.Role, error)
	CreateDefaultRoles(ctx context.Context, input *command.CreateRoleInput) (*domain.Role, error)
	GetRole(ctx context.Context, input *command.RoleIDInput) (*domain.Role, error)
	ListRoles(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Role, error)
	DeleteRole(ctx context.Context, input *command.RoleIDInput) error
}
type Config struct {
	DefaultRoles []string
}
type roleService struct {
	roleProvider domain.RoleProvider
	config       *Config
	*base.BaseService
}

func NewRoleService(base *base.BaseService, repo domain.RoleProvider, config Config) RoleService {
	return &roleService{
		roleProvider: repo,
		BaseService:  base,
		config:       &config,
	}
}
func (s *roleService) CreateRole(ctx context.Context, input *command.CreateRoleInput) (*domain.Role, error) {
	// err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  resource.ResourceAsset,
	// 	Scope:     scope.Create,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	err := validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	isDefault := s.isDefaultRole(input.Name)
	if isDefault {
		return nil, fmt.Errorf("default role")
	}
	role, err := s.roleProvider.CreateRole(ctx, &command.CreateRoleInput{
		BaseInput:   input.BaseInput,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	return domain.NewRole(role.ID, input.Name), nil
}

// CreateDefaultRoles use for bypass validation
func (s *roleService) CreateDefaultRoles(ctx context.Context, input *command.CreateRoleInput) (role *domain.Role, err error) {
	// err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  resource.ResourceAsset,
	// 	Scope:     scope.Create,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	role, err = s.roleProvider.CreateRole(ctx, &command.CreateRoleInput{
		BaseInput:   input.BaseInput,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create initial role %s: %w", input.Name, err)
	}

	return role, nil
}

func (s *roleService) GetRole(ctx context.Context, input *command.RoleIDInput) (*domain.Role, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Role,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	asset, err := s.roleProvider.GetRole(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	return asset, nil
}

func (s *roleService) ListRoles(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Role, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resourceDomain.Role,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	roles, err := s.roleProvider.GetRoles(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return roles, nil
}

func (s *roleService) DeleteRole(ctx context.Context, input *command.RoleIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Role,
		Scope:     scopeDomain.Delete,
	})
	if err != nil {
		return err
	}
	role, err := s.roleProvider.GetRole(ctx, input)
	if err != nil {
		return err
	}
	isDefault := s.isDefaultRole(role.Name)
	if isDefault {
		return fmt.Errorf("default role")
	}
	err = s.roleProvider.DeleteRole(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

// isDefaultRole checks if a role name or ID matches any default roles
func (s *roleService) isDefaultRole(identifier string) bool {
	return slices.ContainsFunc(domain.AllRoles(s.config.DefaultRoles), func(role domain.Role) bool {
		return role.Name == identifier
	})
}
