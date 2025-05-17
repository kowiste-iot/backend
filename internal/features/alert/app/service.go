package app

import (
	"backend/internal/features/alert/domain"
	"backend/internal/features/alert/domain/command"

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

type AlertService interface {
	CreateAlert(ctx context.Context, input *command.CreateAlertInput) (*domain.Alert, error)
	GetAlert(ctx context.Context, input *command.AlertIDInput) (*domain.Alert, error)
	ListAlerts(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Alert, error)
	UpdateAlert(ctx context.Context, input *command.UpdateAlertInput) (*domain.Alert, error)
	DeleteAlert(ctx context.Context, input *command.AlertIDInput) error
}
type alertService struct {
	repo     domain.AlertRepository
	assetDep appAsset.AssetDependencyService
	*base.BaseService
}

const (
	featureName string = "alert"
)

func NewService(base *base.BaseService, repo domain.AlertRepository, assetDep appAsset.AssetDependencyService) *alertService {
	return &alertService{
		repo:        repo,
		assetDep:    assetDep,
		BaseService: base,
	}
}

func (s *alertService) CreateAlert(ctx context.Context, input *command.CreateAlertInput) (*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Alert,
		Scope:     scopeDomain.Create,
	})
	if err != nil {
		return nil, err
	}
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	alert, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Parent, input.Description, input.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	err = s.repo.Create(ctx, alert)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	//Update asset parent dependecy
	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput:  input.BaseInput,
		Feature:    featureName,
		Action:     assetCmd.DependencyActionCreate,
		FeatureID:  alert.ID(),
		NewAssetID: alert.Parent(),
	})
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func (s *alertService) GetAlert(ctx context.Context, input *command.AlertIDInput) (*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Alert,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	alert, err := s.repo.FindByID(ctx, input.AlertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	return alert, nil
}

func (s *alertService) ListAlerts(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resourceDomain.Alert,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	alerts, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}
	return alerts, nil
}

func (s *alertService) UpdateAlert(ctx context.Context, input *command.UpdateAlertInput) (*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Alert,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}
	alert, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	oldParent := alert.Parent()

	err = alert.Update(input.Name, input.Parent, input.Description, input.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	if err := s.repo.Update(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to update alert: %w", err)
	}

	if oldParent != alert.Parent() {
		//Update asset parent dependecy
		err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
			BaseInput:       input.BaseInput,
			PreviousAssetID: oldParent,
			Feature:         featureName,
			Action:          assetCmd.DependencyActionUpdate,
			FeatureID:       alert.ID(),
			NewAssetID:      alert.Parent(),
		})
		if err != nil {
			return nil, err
		}
	}

	return alert, nil
}
func (s *alertService) DeleteAlert(ctx context.Context, input *command.AlertIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Alert,
		Scope:     scopeDomain.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, input.AlertID)
	if err != nil {
		return err
	}

	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput: input.BaseInput,
		Feature:   featureName,
		Action:    assetCmd.DependencyActionDelete,
		FeatureID: input.AlertID,
	})
	if err != nil {
		return err
	}

	return nil
}
