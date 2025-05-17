package repository

import (
	"backend/internal/features/asset/domain"
	gormhelper "backend/shared/gorm"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type assetDependencyRepository struct {
	db *gorm.DB
}

func NewDependencyRepository(db *gorm.DB) domain.AssetDependencyRepository {
	return &assetDependencyRepository{db: db}
}

func (r *assetDependencyRepository) Create(ctx context.Context, dependency *domain.AssetDependency) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	// Verify asset exists
	var count int64
	err = db.WithContext(ctx).Model(&AssetDB{}).
		Where(gormhelper.DeleteFilter()+" AND id = ?", dependency.AssetID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("referenced asset does not exist")
	}

	return r.db.WithContext(ctx).Create(assetDependenciesMapper(dependency)).Error
}

func (r *assetDependencyRepository) Update(ctx context.Context, dependency *domain.AssetDependency) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	// Verify new asset exists
	var count int64
	err = db.WithContext(ctx).Model(&AssetDB{}).
		Where(gormhelper.DeleteFilter()+" AND id = ?", dependency.AssetID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("referenced asset does not exist")
	}

	result := db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()+
			" AND feature_id = ? AND feature = ?",
			dependency.FeatureID,
			dependency.Feature).
		Updates(&AssetDependencyDB{
			AssetID: dependency.AssetID,
		})

	if result.RowsAffected == 0 {
		return errors.New("dependency not found")
	}
	return result.Error
}

func (r *assetDependencyRepository) Remove(ctx context.Context, tenantID, branchID, featureID string) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	//Unscope delete matched records permanently
	result := db.WithContext(ctx).Unscoped().
		Where(gormhelper.DeleteFilter()+" AND feature_id = ?", featureID).
		Delete(&AssetDependencyDB{})

	if result.RowsAffected == 0 {
		return errors.New("dependency not found")
	}
	return result.Error
}

func (r *assetDependencyRepository) FindByFeatureID(ctx context.Context, tenantID, branchID, featureID string) (*domain.AssetDependency, error) {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	var dbDependency AssetDependencyDB
	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()+" AND feature_id = ?", featureID).
		First(&dbDependency).Error
	if err != nil {
		return nil, err
	}

	return domain.NewAssetDependency(
		"dbDependency.TenantID",
		"dbDependency.BranchID",
		dbDependency.FeatureID,
		dbDependency.Feature,
		dbDependency.AssetID,
	), nil
}

func (r *assetDependencyRepository) FindByAssetID(ctx context.Context, tenantID, branchID, assetID string) ([]*domain.AssetDependency, error) {
		// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil,fmt.Errorf("failed to set branch DB: %w", err)
	}
	var dbDependencies []AssetDependencyDB
	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()+" AND asset_id = ?", assetID).
		Find(&dbDependencies).Error
	if err != nil {
		return nil, err
	}

	dependencies := make([]*domain.AssetDependency, len(dbDependencies))
	for i, dep := range dbDependencies {
		dependencies[i] = domain.NewAssetDependency(
			"dep.TenantID",
			"dep.BranchID",
			dep.FeatureID,
			dep.Feature,
			dep.AssetID,
		)
	}
	return dependencies, nil
}
