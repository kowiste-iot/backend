package app

import (
	"backend/internal/features/permission/domain"
	"backend/internal/features/permission/domain/command"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type PermissionService interface {
	CreatePermission(ctx context.Context, input *command.CreatePermissionInput) (*domain.Permission, error)
	ListPermissions(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Permission, error)
}
type Config struct {
	DefaultRoles []string
}

type ServiceDependencies struct {
	Repo   domain.PermissionProvider
	Config *Config
}
type permissionService struct {
	permissionProvider domain.PermissionProvider
	config             *Config
	*base.BaseService
}

func NewService(base *base.BaseService, deps *ServiceDependencies) PermissionService {
	return &permissionService{
		permissionProvider: deps.Repo,
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
	permission, err = s.permissionProvider.CreatePermission(ctx, &input.BaseInput, *p)
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
