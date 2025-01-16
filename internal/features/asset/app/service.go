package app

import (
	"context"
	"ddd/internal/features/asset/domain"
	"ddd/internal/features/asset/domain/command"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
	"ddd/shared/base"
	baseCmd "ddd/shared/base/command"
	"errors"
	"fmt"
)

type AssetService interface {
	CreateAsset(ctx context.Context, input *command.CreateAssetInput) (*domain.Asset, error)
	GetAsset(ctx context.Context, input *command.AssetIDInput) (*domain.Asset, error)
	ListAssets(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Asset, error)
	UpdateAsset(ctx context.Context, input *command.UpdateAssetInput) (*domain.Asset, error)
	DeleteAsset(ctx context.Context, input *command.AssetIDInput) error
}
type assetService struct {
	repo domain.AssetRepository
	*base.BaseService
}

func NewService(base *base.BaseService, repo domain.AssetRepository) AssetService {
	return &assetService{
		repo:        repo,
		BaseService: base,
	}
}
func (s *assetService) CreateAsset(ctx context.Context, input *command.CreateAssetInput) (*domain.Asset, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Asset,
		Scope:     scope.Create,
	})
	if err != nil {
		return nil, err
	}

	asset, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}
	if input.Parent != "" {
		asset.WithParent(input.Parent)
	}
	err = s.repo.Create(ctx, &command.CreateAssetInput{
		BaseInput: baseCmd.BaseInput{
			TenantDomain: input.TenantDomain,
			BranchName:   input.BranchName,
		},
		ID:          asset.ID(),
		Name:        asset.Name(),
		Parent:      *asset.Parent(),
		Description: asset.Description(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	return asset, nil
}

func (s *assetService) GetAsset(ctx context.Context, input *command.AssetIDInput) (*domain.Asset, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Asset,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	asset, err := s.repo.FindByID(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	return asset, nil
}

func (s *assetService) ListAssets(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Asset, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resource.Asset,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	assets, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return assets, nil
}

func (s *assetService) UpdateAsset(ctx context.Context, input *command.UpdateAssetInput) (*domain.Asset, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Asset,
		Scope:     scope.Update,
	})
	if err != nil {
		return nil, err
	}
	asset, err := s.repo.FindByID(ctx, &command.AssetIDInput{
		BaseInput: input.BaseInput,
		AssetID:   input.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	err = asset.Update(input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	if err := s.repo.Update(ctx, &command.UpdateAssetInput{
		BaseInput: baseCmd.BaseInput{
			TenantDomain: input.TenantDomain,
			BranchName:   input.BranchName,
		},
		ID:          asset.ID(),
		Name:        asset.Name(),
		Parent:      *asset.Parent(),
		Description: asset.Description(),
	}); err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	return asset, nil
}
func (s *assetService) DeleteAsset(ctx context.Context, input *command.AssetIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.Asset,
		Scope:     scope.Delete,
	})
	if err != nil {
		return err
	}
	has, err := s.repo.HasChildren(ctx, input)
	if err != nil {
		return err
	}
	if has {
		return errors.New("asset has children and cannot be deleted")
	}
	err = s.repo.Remove(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
