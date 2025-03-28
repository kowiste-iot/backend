package services

import (
	appAsset "backend/internal/features/asset/app"
	repoAsset "backend/internal/features/asset/infra/gorm"
)

func (c *Container) initializeAssetServices(s *Services) error {
	assetDepRepo := repoAsset.NewDependencyRepository(c.base.DB)
	s.AssetDepService = appAsset.NewAssetDependencyService(c.base, assetDepRepo)

	assetRepo := repoAsset.NewRepository(c.base.DB)
	s.AssetService = appAsset.NewService(c.base, assetRepo)

	return nil
}
