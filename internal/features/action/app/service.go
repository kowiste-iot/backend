package app

import (
	"backend/internal/features/action/domain"
	"backend/internal/features/action/domain/command"

	appAsset "backend/internal/features/asset/app"
	assetCmd "backend/internal/features/asset/domain/command"

	resourceDomain "backend/internal/features/resource/domain"
	scopeDomain "backend/internal/features/scope/domain"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
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
	repo     domain.ActionRepository
	assetDep appAsset.AssetDependencyService
	*base.BaseService
}

const (
	featureName string = "action"
)

func NewService(base *base.BaseService, repo domain.ActionRepository, assetDep appAsset.AssetDependencyService) *actionService {
	return &actionService{
		repo:        repo,
		assetDep:    assetDep,
		BaseService: base,
	}
}
func (s *actionService) CreateAction(ctx context.Context, input *command.CreateActionInput) (*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Action,
		Scope:     scopeDomain.Create,
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

	//Update asset parent dependecy
	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput:  input.BaseInput,
		Feature:    featureName,
		Action:     assetCmd.DependencyActionCreate,
		FeatureID:  action.ID(),
		NewAssetID: action.Parent(),
	})
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (s *actionService) GetAction(ctx context.Context, input *command.ActionIDInput) (*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Action,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	action, err := s.repo.FindByID(ctx, input.ActionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get action: %w", err)
	}
	return action, nil
}

func (s *actionService) ListActions(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resourceDomain.Action,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	actions, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list actions: %w", err)
	}
	return actions, nil
}

func (s *actionService) UpdateAction(ctx context.Context, input *command.UpdateActionInput) (*domain.Action, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Action,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}
	action, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get action: %w", err)
	}

	oldParent := action.Parent()

	err = action.Update(input.Name, input.Parent, input.Description, input.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to get action: %w", err)
	}
	if err := s.repo.Update(ctx, action); err != nil {
		return nil, fmt.Errorf("failed to update action: %w", err)
	}

	if oldParent != action.Parent() {
		//Update asset parent dependecy
		err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
			BaseInput:       input.BaseInput,
			PreviousAssetID: oldParent,
			Feature:         featureName,
			Action:          assetCmd.DependencyActionUpdate,
			FeatureID:       action.ID(),
			NewAssetID:      action.Parent(),
		})
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (s *actionService) DeleteAction(ctx context.Context, input *command.ActionIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Action,
		Scope:     scopeDomain.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, input.ActionID)
	if err != nil {
		return err
	}

	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput: input.BaseInput,
		Feature:   featureName,
		Action:    assetCmd.DependencyActionDelete,
		FeatureID: input.ActionID,
	})
	if err != nil {
		return err
	}

	return nil
}
