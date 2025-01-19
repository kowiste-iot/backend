package repository

import (
	"context"
	"ddd/internal/features/asset/domain"
	baseCmd "ddd/shared/base/command"
	gormhelper "ddd/shared/gorm"
	"ddd/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type assetRepository struct {
	db *gorm.DB
}

type Asset struct {
	ID          string  `gorm:"primaryKey"`
	TenantID    string  `gorm:"index"`
	BranchID    string  `gorm:"index"`
	Parent      *string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewRepository(db *gorm.DB) domain.AssetRepository {
	db.AutoMigrate(&Asset{})
	return &assetRepository{db: db}
}

func (r *assetRepository) Create(ctx context.Context, input *domain.Asset) error {
	dbAsset := Asset{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Create(&dbAsset).Error
}

func (r *assetRepository) Update(ctx context.Context, input *domain.Asset) error {
	dbAsset := Asset{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Updates(&dbAsset).Error
}

func (r *assetRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, assetID string) (*domain.Asset, error) {
	var dbAsset Asset

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", assetID).
		First(&dbAsset).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbAsset.ID,
		dbAsset.TenantID,
		dbAsset.BranchID,
		dbAsset.Name,
		dbAsset.Description,
		dbAsset.Parent,
		dbAsset.UpdatedAt,
		&dbAsset.DeletedAt.Time,
	), nil

}

func (r *assetRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Asset, error) {
	var dbAssets []Asset

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Asset{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
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
			dbAsset.TenantID,
			dbAsset.BranchID,
			dbAsset.Name,
			dbAsset.Description,
			dbAsset.Parent,
			dbAsset.UpdatedAt,
			&dbAsset.DeletedAt.Time,
		)
	}
	return assets, nil
}

func (r *assetRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, assetID string) error {

	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", assetID).Delete(&Asset{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *assetRepository) HasChildren(ctx context.Context, input *baseCmd.BaseInput, assetID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Asset{}).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND parent = ?", assetID).
		Count(&count).Error

	return count > 0, err
}
