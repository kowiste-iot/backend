package app

import (
	"backend/internal/features/measure/domain"
	"backend/internal/features/measure/domain/command"

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

type MeasureService interface {
	CreateMeasure(ctx context.Context, input *command.CreateMeasureInput) (*domain.Measure, error)
	GetMeasure(ctx context.Context, input *command.MeasureIDInput) (*domain.Measure, error)
	ListMeasures(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Measure, error)
	UpdateMeasure(ctx context.Context, input *command.UpdateMeasureInput) (*domain.Measure, error)
	DeleteMeasure(ctx context.Context, input *command.MeasureIDInput) error
}
type measureService struct {
	repo     domain.MeasureRepository
	assetDep appAsset.AssetDependencyService
	*base.BaseService
}

const (
	featureName string = "measure"
)

func NewService(base *base.BaseService, repo domain.MeasureRepository, assetDep appAsset.AssetDependencyService) MeasureService {
	return &measureService{
		repo:        repo,
		assetDep:    assetDep,
		BaseService: base,
	}
}

func (s *measureService) CreateMeasure(ctx context.Context, input *command.CreateMeasureInput) (*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Measure,
		Scope:     scopeDomain.Create,
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

	//Update asset parent dependecy
	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput:  input.BaseInput,
		Feature:    featureName,
		Action:     assetCmd.DependencyActionCreate,
		FeatureID:  measure.ID(),
		NewAssetID: measure.Parent(),
	})
	if err != nil {
		return nil, err
	}

	return measure, nil
}

func (s *measureService) GetMeasure(ctx context.Context, input *command.MeasureIDInput) (*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Measure,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	measure, err := s.repo.FindByID(ctx, input.MeasureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get measure: %w", err)
	}
	return measure, nil
}

func (s *measureService) ListMeasures(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resourceDomain.Measure,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	measures, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list measures: %w", err)
	}
	return measures, nil
}

func (s *measureService) UpdateMeasure(ctx context.Context, input *command.UpdateMeasureInput) (*domain.Measure, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Measure,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}

	measure, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get measure: %w", err)
	}

	oldParent := measure.Parent()

	err = measure.Update(input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get measure: %w", err)
	}
	if err := s.repo.Update(ctx, measure); err != nil {
		return nil, fmt.Errorf("failed to update measure: %w", err)
	}

	if oldParent != measure.Parent() {
		//Update asset parent dependecy
		err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
			BaseInput:       input.BaseInput,
			PreviousAssetID: oldParent,
			Feature:         featureName,
			Action:          assetCmd.DependencyActionUpdate,
			FeatureID:       measure.ID(),
			NewAssetID:      measure.Parent(),
		})
		if err != nil {
			return nil, err
		}
	}

	return measure, nil
}

func (s *measureService) DeleteMeasure(ctx context.Context, input *command.MeasureIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Measure,
		Scope:     scopeDomain.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, input.MeasureID)
	if err != nil {
		return err
	}

	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput: input.BaseInput,
		Feature:   featureName,
		Action:    assetCmd.DependencyActionDelete,
		FeatureID: input.MeasureID,
	})
	if err != nil {
		return err
	}

	return nil
}
