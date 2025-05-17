package app

import (
	"backend/internal/features/asset/domain"
	"backend/internal/features/asset/domain/command"
	resourceDomain "backend/internal/features/resource/domain"
	scopeDomain "backend/internal/features/scope/domain"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
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
		Resource:  resourceDomain.Asset,
		Scope:     scopeDomain.Create,
	})
	if err != nil {
		return nil, err
	}
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	asset, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}
	if input.Parent != "" {
		asset.WithParent(input.Parent)
	}

	err = s.repo.Create(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	return asset, nil
}

func (s *assetService) GetAsset(ctx context.Context, input *command.AssetIDInput) (*domain.Asset, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Asset,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	asset, err := s.repo.FindByID(ctx, input.AssetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	return asset, nil
}

func (s *assetService) ListAssets(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Asset, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resourceDomain.Asset,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	assets, err := s.repo.FindAll(ctx)
 	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return assets, nil
}

func (s *assetService) UpdateAsset(ctx context.Context, input *command.UpdateAssetInput) (*domain.Asset, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Asset,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}
	asset, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	err = asset.Update(input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	if err := s.repo.Update(ctx, asset); err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	return asset, nil
}
func (s *assetService) DeleteAsset(ctx context.Context, input *command.AssetIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Asset,
		Scope:     scopeDomain.Delete,
	})
	if err != nil {
		return err
	}
	has, err := s.repo.HasChildren(ctx, input.AssetID)
	if err != nil {
		return err
	}
	if has {
		return errors.New("asset has children and cannot be deleted")
	}
	err = s.repo.Remove(ctx, input.AssetID)
	if err != nil {
		return err
	}
	return nil
}
