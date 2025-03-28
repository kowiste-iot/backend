package repository

import (
	"backend/internal/features/asset/domain"
	gormhelper "backend/shared/gorm"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AssetDependency struct {
    TenantID  string `gorm:"index;uniqueIndex:idx_dependency"`
    BranchID  string `gorm:"index;uniqueIndex:idx_dependency"`
    FeatureID string `gorm:"uniqueIndex:idx_dependency"`
    Feature   string `gorm:"uniqueIndex:idx_dependency"`
    AssetID   string `gorm:"type:string;references:ID;constraint:OnDelete:RESTRICT"`
    Asset     Asset  `gorm:"foreignKey:AssetID;references:ID"` 
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type assetDependencyRepository struct {
	db *gorm.DB
}

func NewDependencyRepository(db *gorm.DB) domain.AssetDependencyRepository {
	db.AutoMigrate(&AssetDependency{})
	return &assetDependencyRepository{db: db}
}

func (r *assetDependencyRepository) Create(ctx context.Context, dependency *domain.AssetDependency) error {
    // Verify asset exists
    var count int64
    err := r.db.WithContext(ctx).Model(&Asset{}).
        Where(gormhelper.TenantBranchFilter(dependency.TenantID, dependency.BranchID)+" AND id = ?", dependency.AssetID).
        Count(&count).Error
    if err != nil {
        return err
    }
    if count == 0 {
        return errors.New("referenced asset does not exist")
    }

    dbDependency := AssetDependency{
        TenantID:  dependency.TenantID,
        BranchID:  dependency.BranchID,
        FeatureID: dependency.FeatureID,
        Feature:   string(dependency.Feature),
        AssetID:   dependency.AssetID,
    }

    return r.db.WithContext(ctx).Create(&dbDependency).Error
}

func (r *assetDependencyRepository) Update(ctx context.Context, dependency *domain.AssetDependency) error {
    // Verify new asset exists
    var count int64
    err := r.db.WithContext(ctx).Model(&Asset{}).
        Where(gormhelper.TenantBranchFilter(dependency.TenantID, dependency.BranchID)+" AND id = ?", dependency.AssetID).
        Count(&count).Error
    if err != nil {
        return err
    }
    if count == 0 {
        return errors.New("referenced asset does not exist")
    }

    result := r.db.WithContext(ctx).
        Where(gormhelper.TenantBranchFilter(dependency.TenantID, dependency.BranchID)+
            " AND feature_id = ? AND feature = ?", 
            dependency.FeatureID, 
            dependency.Feature).
        Updates(&AssetDependency{
            AssetID: dependency.AssetID,
        })

    if result.RowsAffected == 0 {
        return errors.New("dependency not found")
    }
    return result.Error
}

func (r *assetDependencyRepository) Remove(ctx context.Context, tenantID, branchID, featureID string) error {
	//Unscope delete matched records permanently
	result := r.db.WithContext(ctx).Unscoped().
		Where(gormhelper.TenantBranchFilter(tenantID, branchID)+" AND feature_id = ?", featureID).
		Delete(&AssetDependency{})

	if result.RowsAffected == 0 {
		return errors.New("dependency not found")
	}
	return result.Error
}

func (r *assetDependencyRepository) FindByFeatureID(ctx context.Context, tenantID, branchID, featureID string) (*domain.AssetDependency, error) {
	var dbDependency AssetDependency
	err := r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(tenantID, branchID)+" AND feature_id = ?", featureID).
		First(&dbDependency).Error
	if err != nil {
		return nil, err
	}

	return domain.NewAssetDependency(
		dbDependency.TenantID,
		dbDependency.BranchID,
		dbDependency.FeatureID,
		dbDependency.Feature,
		dbDependency.AssetID,
	), nil
}

func (r *assetDependencyRepository) FindByAssetID(ctx context.Context, tenantID, branchID, assetID string) ([]*domain.AssetDependency, error) {
	var dbDependencies []AssetDependency
	err := r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(tenantID, branchID)+" AND asset_id = ?", assetID).
		Find(&dbDependencies).Error
	if err != nil {
		return nil, err
	}

	dependencies := make([]*domain.AssetDependency, len(dbDependencies))
	for i, dep := range dbDependencies {
		dependencies[i] = domain.NewAssetDependency(
			dep.TenantID,
			dep.BranchID,
			dep.FeatureID,
			dep.Feature,
			dep.AssetID,
		)
	}
	return dependencies, nil
}
