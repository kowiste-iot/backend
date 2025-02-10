package app

import (
	"backend/internal/features/dashboard/domain"
	"backend/internal/features/dashboard/domain/command"
	"backend/shared/auth/domain/resource"
	"backend/shared/auth/domain/scope"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type DashboardService interface {
	CreateDashboard(ctx context.Context, input *command.CreateDashboardInput) (*domain.Dashboard, error)
	GetDashboard(ctx context.Context, input *command.DashboardIDInput) (*domain.Dashboard, error)
	ListDashboards(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Dashboard, error)
	UpdateDashboard(ctx context.Context, input *command.UpdateDashboardInput) (*domain.Dashboard, error)
	DeleteDashboard(ctx context.Context, input *command.DashboardIDInput) error
}
type dashboardService struct {
	repo domain.DashboardRepository
	*base.BaseService
}

func NewService(base *base.BaseService, repo domain.DashboardRepository) DashboardService {
	return &dashboardService{
		repo:        repo,
		BaseService: base,
	}
}
func (s *dashboardService) CreateDashboard(ctx context.Context, input *command.CreateDashboardInput) (*domain.Dashboard, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Dashboard,
		Scope:     scope.Create,
	})
	if err != nil {
		return nil, err
	}
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	dashboard, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard: %w", err)
	}

	err = s.repo.Create(ctx, dashboard)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard: %w", err)
	}

	return dashboard, nil
}

func (s *dashboardService) GetDashboard(ctx context.Context, input *command.DashboardIDInput) (*domain.Dashboard, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Dashboard,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	dashboard, err := s.repo.FindByID(ctx, &input.BaseInput, input.DashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	return dashboard, nil
}

func (s *dashboardService) ListDashboards(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Dashboard, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resource.Dashboard,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	dashboards, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards: %w", err)
	}
	return dashboards, nil
}

func (s *dashboardService) UpdateDashboard(ctx context.Context, input *command.UpdateDashboardInput) (*domain.Dashboard, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Dashboard,
		Scope:     scope.Update,
	})
	if err != nil {
		return nil, err
	}
	dashboard, err := s.repo.FindByID(ctx, &input.BaseInput, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	err = dashboard.Update(input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}
	if err := s.repo.Update(ctx, dashboard); err != nil {
		return nil, fmt.Errorf("failed to update dashboard: %w", err)
	}

	return dashboard, nil
}
func (s *dashboardService) DeleteDashboard(ctx context.Context, input *command.DashboardIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Dashboard,
		Scope:     scope.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, &input.BaseInput, input.DashboardID)
	if err != nil {
		return err
	}
	return nil
}
