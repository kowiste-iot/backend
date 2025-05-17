package repository

import (
	"backend/internal/features/asset/domain"
	gormhelper "backend/shared/gorm"
	"backend/shared/pagination"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type assetRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) Create(ctx context.Context, input *domain.Asset) (err error) {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Create(assetMapper((input))).Error
}

func (r *assetRepository) Update(ctx context.Context, input *domain.Asset) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Updates(assetMapper(input)).Error
}

func (r *assetRepository) FindByID(ctx context.Context, assetID string) (*domain.Asset, error) {
	var dbAsset AssetDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", assetID).
		First(&dbAsset).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbAsset.ID,
		"dbAsset.TenantID",
		"dbAsset.BranchID",
		dbAsset.Name,
		dbAsset.Description,
		dbAsset.Parent,
		dbAsset.UpdatedAt,
		&dbAsset.DeletedAt.Time,
	), nil

}

func (r *assetRepository) FindAll(ctx context.Context) ([]*domain.Asset, error) {
	var dbAssets []AssetDB

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	var total int64
	err = db.Model(&AssetDB{}).Where(gormhelper.DeleteFilter()).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbAssets).Error
	if err != nil {
		return nil, err
	}

	assets := make([]*domain.Asset, len(dbAssets))
	for i, dbAsset := range dbAssets {
		assets[i] = domain.NewFromRepository(
			dbAsset.ID,
			"dbAsset.TenantID",
			"dbAsset.BranchID",
			dbAsset.Name,
			dbAsset.Description,
			dbAsset.Parent,
			dbAsset.UpdatedAt,
			&dbAsset.DeletedAt.Time,
		)
	}
	return assets, nil
}

func (r *assetRepository) Remove(ctx context.Context, assetID string) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	resp := db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", assetID).Delete(&AssetDB{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *assetRepository) HasChildren(ctx context.Context, assetID string) (bool, error) {
	var count int64
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return false, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).
		Model(&AssetDependencyDB{}).
		Where(gormhelper.DeleteFilter()+" AND parent = ?", assetID).
		Count(&count).Error

	return count > 0, err
}
