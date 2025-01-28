package app

import (
	"context"
	"ddd/internal/features/measure/domain"
	"ddd/internal/features/measure/domain/command"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
	"ddd/shared/base"
	baseCmd "ddd/shared/base/command"
	"ddd/shared/validator"
	"fmt"
)

type MeasureService interface {
	CreateMeasure(ctx context.Context, input *command.CreateMeasureInput) (*domain.Measure, error)
	GetMeasure(ctx context.Context, input *command.MeasureIDInput) (*domain.Measure, error)
	ListMeasures(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Measure, error)
	UpdateMeasure(ctx context.Context, input *command.UpdateMeasureInput) (*domain.Measure, error)
	DeleteMeasure(ctx context.Context, input *command.MeasureIDInput) error
}
type measureService struct {
	repo domain.MeasureRepository
	*base.BaseService
}

func NewService(base *base.BaseService, repo domain.MeasureRepository) MeasureService {
	return &measureService{
		repo:        repo,
		BaseService: base,
	}
}
func (s *measureService) CreateMeasure(ctx context.Context, input *command.CreateMeasureInput) (*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Measure,
		Scope:     scope.Create,
	})
	if err != nil {
		return nil, err
	}
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	measure, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create measure: %w", err)
	}

	err = s.repo.Create(ctx, measure)
	if err != nil {
		return nil, fmt.Errorf("failed to create measure: %w", err)
	}

	return measure, nil
}

func (s *measureService) GetMeasure(ctx context.Context, input *command.MeasureIDInput) (*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Measure,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	measure, err := s.repo.FindByID(ctx, &input.BaseInput, input.MeasureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get measure: %w", err)
	}
	return measure, nil
}

func (s *measureService) ListMeasures(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resource.Measure,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	measures, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list measures: %w", err)
	}
	return measures, nil
}

func (s *measureService) UpdateMeasure(ctx context.Context, input *command.UpdateMeasureInput) (*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Measure,
		Scope:     scope.Update,
	})
	if err != nil {
		return nil, err
	}
	measure, err := s.repo.FindByID(ctx, &input.BaseInput, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get measure: %w", err)
	}
	err = measure.Update(input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get measure: %w", err)
	}
	if err := s.repo.Update(ctx, measure); err != nil {
		return nil, fmt.Errorf("failed to update measure: %w", err)
	}

	return measure, nil
}
func (s *measureService) DeleteMeasure(ctx context.Context, input *command.MeasureIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Measure,
		Scope:     scope.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, &input.BaseInput, input.MeasureID)
	if err != nil {
		return err
	}
	return nil
}
