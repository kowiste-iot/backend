package repository

import (
	"context"
	"ddd/internal/features/tenant/domain"
	gormhelper "ddd/shared/gorm"
	"ddd/shared/http/httputil"
	"ddd/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type tenantRepository struct {
	db *gorm.DB
}

type Tenant struct {
	ID          string `gorm:"primaryKey"`
	AuthID      string `gorm:"index"`
	Name        string
	Domain      string `gorm:"index"`
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewTenantRepository(db *gorm.DB) domain.TenantRepository {
	db.AutoMigrate(&Tenant{})
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, asset *domain.Tenant) error {
	dbAsset := Tenant{
		ID:          asset.ID(),
		AuthID:      asset.AuhtID(),
		Name:        asset.Name(),
		Description: asset.Description(),
	}
	return r.db.WithContext(ctx).Save(&dbAsset).Error
}
func (r *tenantRepository) Update(ctx context.Context, asset *domain.Tenant) error {
	dbAsset := Tenant{
		ID:          asset.ID(),
		AuthID:      asset.AuhtID(),
		Name:        asset.Name(),
		Description: asset.Description(),
	}
	return r.db.WithContext(ctx).Save(&dbAsset).Error
}

func (r *tenantRepository) FindByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	var dbTenant Tenant

	err := r.db.WithContext(ctx).Where("auth_id = ?", tenantID).First(&dbTenant).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbTenant.ID,
		dbTenant.AuthID,
		dbTenant.Name,
		dbTenant.Domain,
		dbTenant.Description,
		dbTenant.UpdatedAt,
		&dbTenant.DeletedAt.Time,
	), nil

}

func (r *tenantRepository) FindAll(ctx context.Context) ([]*domain.Tenant, error) {
	var dbAssets []Tenant
	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		return nil, errors.New("not tenant id")
	}
	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	if err := r.db.Model(&Tenant{}).Where(gormhelper.TenantFilter(tenant.ID())).Count(&total).Error; err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err := r.db.WithContext(ctx).
		Where(gormhelper.TenantFilter(tenant.ID())).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbAssets).Error
	if err != nil {
		return nil, err
	}

	assets := make([]*domain.Tenant, len(dbAssets))
	for i, dbAsset := range dbAssets {
		assets[i] = domain.NewFromRepository(
			dbAsset.ID,
			dbAsset.AuthID,
			dbAsset.Name,
			dbAsset.Domain,
			dbAsset.Description,
			dbAsset.UpdatedAt,
			&dbAsset.DeletedAt.Time,
		)
	}
	return assets, nil
}

func (r *tenantRepository) Remove(ctx context.Context, assetID string) error {
	tenantID, ok := httputil.GetTenant(ctx)
	if !ok {
		return errors.New("not tenant id")
	}
	resp := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, assetID).Delete(&Tenant{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *tenantRepository) HasChildren(ctx context.Context, assetID string) (bool, error) {
	var count int64
	tenantID, ok := httputil.GetTenant(ctx)
	if !ok {
		return false, errors.New("not tenant id")
	}

	err := r.db.WithContext(ctx).
		Model(&Tenant{}).
		Where("tenant_id = ? AND parent = ?", tenantID, assetID).
		Count(&count).Error

	return count > 0, err
}
