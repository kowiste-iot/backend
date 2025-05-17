package app

import (
	"backend/internal/features/scope/domain"
	"backend/internal/features/scope/domain/command"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type ScopeService interface {
	CreateScope(ctx context.Context, input *command.CreateScopeInput) (*domain.Scope, error)
	ListScopes(ctx context.Context, input *baseCmd.BaseInput) ([]domain.Scope, error)
}

type ServiceDependencies struct {
	Repo domain.ScopeProvider
}
type scopeService struct {
	scopeProvider domain.ScopeProvider
	*base.BaseService
}

func NewService(base *base.BaseService, deps *ServiceDependencies) ScopeService {
	return &scopeService{
		scopeProvider: deps.Repo,
		BaseService:   base,
	}
}
func (s *scopeService) CreateScope(ctx context.Context, input *command.CreateScopeInput) (scope *domain.Scope, err error) {
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
	sc, err := domain.New(input.Name, input.DisplayName)
	if err != nil {
		return
	}
	scope, err = s.scopeProvider.CreateScope(ctx, &input.BaseInput, *sc)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return
}

func (s *scopeService) ListScopes(ctx context.Context, input *baseCmd.BaseInput) (scopes []domain.Scope, err error) {
	// err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: *input,
	// 	Resource:  domain.ResourceR,
	// 	Scope:     scope.View,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	scopes, err = s.scopeProvider.ListScopes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	return
}
