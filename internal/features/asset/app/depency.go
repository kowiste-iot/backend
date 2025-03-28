package app

import (
	"backend/internal/features/asset/domain"
	"backend/internal/features/asset/domain/command"
	"backend/shared/base"
	"backend/shared/validator"
	"context"
	"fmt"
)

type AssetDependencyService interface {
	UpdateDependency(ctx context.Context, input *command.DependencyChangeInput) error
}

type assetDependencyService struct {
	repo domain.AssetDependencyRepository
	*base.BaseService
}

func NewAssetDependencyService(base *base.BaseService, repo domain.AssetDependencyRepository) AssetDependencyService {
	return &assetDependencyService{
		repo:        repo,
		BaseService: base,
	}
}

func (s *assetDependencyService) UpdateDependency(ctx context.Context, input *command.DependencyChangeInput) error {
	err := validator.Validate(input)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	switch input.Action {
	case command.DependencyActionCreate:
		dependency := domain.NewAssetDependency(
			input.TenantDomain,
			input.BranchName,
			input.FeatureID,
			input.Feature,
			input.NewAssetID,
		)
		return s.repo.Create(ctx, dependency)

	case command.DependencyActionUpdate:
		dependency, err := s.repo.FindByFeatureID(ctx, input.TenantDomain, input.BranchName, input.FeatureID)
		if err != nil {
			return fmt.Errorf("failed to find dependency: %w", err)
		}
		dependency.UpdateAsset(input.NewAssetID)
		return s.repo.Update(ctx, dependency)

	case command.DependencyActionDelete:
		return s.repo.Remove(ctx, input.TenantDomain, input.BranchName, input.FeatureID)

	default:
		return fmt.Errorf("invalid action: %s", input.Action)
	}
}
