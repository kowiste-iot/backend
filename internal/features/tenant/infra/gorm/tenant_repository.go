package repository

import (
	"backend/internal/features/tenant/domain"
	gormhelper "backend/shared/gorm"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"context"
	"errors"

	"gorm.io/gorm"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) domain.TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, data *domain.Tenant) error {
	return r.db.WithContext(ctx).Save(tenantMapper(data)).Error
}
func (r *tenantRepository) Update(ctx context.Context, data *domain.Tenant) error {

	return r.db.WithContext(ctx).Save(tenantMapper(data)).Error
}

func (r *tenantRepository) FindByID(ctx context.Context, tenantDomain string) (*domain.Tenant, error) {
	var dbTenant TenantDB

	err := r.db.WithContext(ctx).Where("domain = ?", tenantDomain).First(&dbTenant).Error
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
	var dbAssets []TenantDB
	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		return nil, errors.New("not tenant id")
	}
	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	if err := r.db.Model(&TenantDB{}).Where(gormhelper.TenantFilter(tenant.ID())).Count(&total).Error; err != nil {
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
	resp := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, assetID).Delete(&TenantDB{})
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
		Model(&TenantDB{}).
		Where("tenant_id = ? AND parent = ?", tenantID, assetID).
		Count(&count).Error

	return count > 0, err
}
