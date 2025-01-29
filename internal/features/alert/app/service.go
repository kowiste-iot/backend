package app

import (
	"context"
	"ddd/internal/features/alert/domain"
	"ddd/internal/features/alert/domain/command"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
	"ddd/shared/base"
	baseCmd "ddd/shared/base/command"
	"ddd/shared/validator"
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
	repo domain.AlertRepository
	*base.BaseService
}

func NewService(base *base.BaseService, repo domain.AlertRepository) *alertService {
	return &alertService{
		repo:        repo,
		BaseService: base,
	}
}
func (s *alertService) CreateAlert(ctx context.Context, input *command.CreateAlertInput) (*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Alert,
		Scope:     scope.Create,
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

	return alert, nil
}

func (s *alertService) GetAlert(ctx context.Context, input *command.AlertIDInput) (*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Alert,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	alert, err := s.repo.FindByID(ctx, &input.BaseInput, input.AlertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	return alert, nil
}

func (s *alertService) ListAlerts(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resource.Alert,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	alerts, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}
	return alerts, nil
}

func (s *alertService) UpdateAlert(ctx context.Context, input *command.UpdateAlertInput) (*domain.Alert, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Alert,
		Scope:     scope.Update,
	})
	if err != nil {
		return nil, err
	}
	alert, err := s.repo.FindByID(ctx, &input.BaseInput, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	err = alert.Update(input.Name, input.Parent, input.Description, input.Enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	if err := s.repo.Update(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to update alert: %w", err)
	}

	return alert, nil
}
func (s *alertService) DeleteAlert(ctx context.Context, input *command.AlertIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Alert,
		Scope:     scope.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, &input.BaseInput, input.AlertID)
	if err != nil {
		return err
	}
	return nil
}
