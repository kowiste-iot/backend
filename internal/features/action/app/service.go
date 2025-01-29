package app

import (
	"context"
	"ddd/internal/features/action/domain"
	"ddd/internal/features/action/domain/command"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
	"ddd/shared/base"
	baseCmd "ddd/shared/base/command"
	"ddd/shared/validator"
	"fmt"
)

type ActionService interface {
	CreateAction(ctx context.Context, input *command.CreateActionInput) (*domain.Action, error)
	GetAction(ctx context.Context, input *command.ActionIDInput) (*domain.Action, error)
	ListActions(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Action, error)
	UpdateAction(ctx context.Context, input *command.UpdateActionInput) (*domain.Action, error)
	DeleteAction(ctx context.Context, input *command.ActionIDInput) error
}
type actionService struct {
	repo domain.ActionRepository
	*base.BaseService
}

func NewService(base *base.BaseService, repo domain.ActionRepository) *actionService {
	return &actionService{
		repo:        repo,
		BaseService: base,
	}
}
func (s *actionService) CreateAction(ctx context.Context, input *command.CreateActionInput) (*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Action,
		Scope:     scope.Create,
	})
	if err != nil {
		return nil, err
	}
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	action, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Parent, input.Description, input.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	err = s.repo.Create(ctx, action)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return action, nil
}

func (s *actionService) GetAction(ctx context.Context, input *command.ActionIDInput) (*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Action,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	action, err := s.repo.FindByID(ctx, &input.BaseInput, input.ActionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get action: %w", err)
	}
	return action, nil
}

func (s *actionService) ListActions(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resource.Action,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	actions, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list actions: %w", err)
	}
	return actions, nil
}

func (s *actionService) UpdateAction(ctx context.Context, input *command.UpdateActionInput) (*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Action,
		Scope:     scope.Update,
	})
	if err != nil {
		return nil, err
	}
	action, err := s.repo.FindByID(ctx, &input.BaseInput, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get action: %w", err)
	}
	err = action.Update(input.Name, input.Parent, input.Description, input.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to get action: %w", err)
	}
	if err := s.repo.Update(ctx, action); err != nil {
		return nil, fmt.Errorf("failed to update action: %w", err)
	}

	return action, nil
}
func (s *actionService) DeleteAction(ctx context.Context, input *command.ActionIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Action,
		Scope:     scope.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, &input.BaseInput, input.ActionID)
	if err != nil {
		return err
	}
	return nil
}
