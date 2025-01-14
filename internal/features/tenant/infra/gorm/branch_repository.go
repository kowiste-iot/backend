package repository

import (
	"context"
	"ddd/internal/features/tenant/domain"
	"ddd/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type branchRepository struct {
	db *gorm.DB
}

type Branch struct {
	ID           string `gorm:"primaryKey"`
	TenantID     string `gorm:"index"`
	AuthBranchID string `gorm:"index"`
	Name         string
	Description  string
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func NewBranchRepository(db *gorm.DB) domain.BranchRepository {
	db.AutoMigrate(&Branch{})
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(ctx context.Context, tenantID string, branch *domain.Branch) error {
	dbBranch := Branch{
		ID:           branch.ID(),
		TenantID:     tenantID,
		AuthBranchID: branch.AuthBranchID(),
		Name:         branch.Name(),
		Description:  branch.Description(),
		UpdatedAt:    branch.UpdatedAt(),
	}
	return r.db.WithContext(ctx).Create(&dbBranch).Error
}
func (r *branchRepository) Update(ctx context.Context, tenantID string, branch *domain.Branch) error {
	dbBranch := Branch{
		ID:           branch.ID(),
		TenantID:     tenantID,
		AuthBranchID: branch.AuthBranchID(),
		Name:         branch.Name(),
		Description:  branch.Description(),
		UpdatedAt:    branch.UpdatedAt(),
	}
	return r.db.WithContext(ctx).Updates(&dbBranch).Error
}

func (r *branchRepository) FindByID(ctx context.Context, tenantID, branchID string) (*domain.Branch, error) {
	var dbBranch Branch

	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", branchID, tenantID).
		First(&dbBranch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBranchNotFound
		}
		return nil, err
	}

	return domain.NewBranchFromRepository(
		dbBranch.ID,
		dbBranch.TenantID,
		dbBranch.AuthBranchID,
		dbBranch.Name,
		dbBranch.Description,
		dbBranch.UpdatedAt,
		&dbBranch.DeletedAt.Time,
	), nil
}

func (r *branchRepository) FindAll(ctx context.Context, tenantID string) ([]*domain.Branch, error) {
	var dbBranches []Branch

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}

	// Count total for pagination
	var total int64
	if err := r.db.Model(&Branch{}).
		Where("tenant_id = ?", tenantID).
		Count(&total).Error; err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	// Get paginated results
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbBranches).Error
	if err != nil {
		return nil, err
	}

	branches := make([]*domain.Branch, len(dbBranches))
	for i, dbBranch := range dbBranches {
		branches[i] = domain.NewBranchFromRepository(
			dbBranch.ID,
			dbBranch.TenantID,
			dbBranch.AuthBranchID,
			dbBranch.Name,
			dbBranch.Description,
			dbBranch.UpdatedAt,
			&dbBranch.DeletedAt.Time,
		)
	}
	return branches, nil
}

func (r *branchRepository) Remove(ctx context.Context, tenantID, branchID string) error {

	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", branchID, tenantID).
		Delete(&Branch{})

	if result.RowsAffected == 0 {
		return domain.ErrBranchNotFound
	}
	return result.Error
}
