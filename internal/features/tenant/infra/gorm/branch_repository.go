package repository

import (
	"backend/internal/features/tenant/domain"
	gormhelper "backend/shared/gorm"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) domain.BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(ctx context.Context, tenantID string, branch *domain.Branch) (err error) {
	//TODO: should be good remove tenant from parameter and pass it in context?

	dbName := gormhelper.GetBranchName(tenantID, branch.Name())
	tenant, exist := httputil.GetTenant(ctx)
	if !exist {
		return errors.New("tenant not found")
	}
	//setting ctx values to be use
	ctx = httputil.SetBranch(ctx, branch.Name())
	ctx = httputil.SetTenant(ctx, tenant)
	err = r.createBranchSchema(ctx, dbName, branch)
	if err != nil {
		return
	}

	return r.db.WithContext(ctx).Create(branchMapper(dbName, tenant, branch)).Error

}

func (r *branchRepository) createBranchSchema(ctx context.Context, dbName string, branch *domain.Branch) error {
	// Create branch schema name

	// Read the branch.sql file and pass the schema name for templating
	sqlContent, err := domain.GetBranchSchemaSQL(dbName, 0)
	if err != nil {
		return fmt.Errorf("failed to get branch schema SQL: %w", err)
	}
	// Execute the SQL to create the schema and tables
	err = r.db.Exec(string(sqlContent)).Error
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	// Get a connection to the newly created schema
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	tenant, _ := httputil.GetTenant(ctx)

	// Save branch info to the branch_info table in the new schema
	return db.WithContext(ctx).Create(branchInfoMapper(tenant, branch)).Error
}
func (r *branchRepository) Update(ctx context.Context, tenantID string, branch *domain.Branch) error {
	tenant, _ := httputil.GetTenant(ctx)
	return r.db.WithContext(ctx).Updates(branchInfoMapper(tenant, branch)).Error
}

func (r *branchRepository) FindByID(ctx context.Context, tenantID, branchID string) (*domain.Branch, error) {
	var dbBranch BranchInfo

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
		dbBranch.BranchID,
		dbBranch.TenantID,
		"",
		"dbBranch.AdminBranchID",
		dbBranch.BranchName,
		dbBranch.Description,
		dbBranch.UpdatedAt,
		nil,
	), nil
}

func (r *branchRepository) FindAll(ctx context.Context, tenantID string) ([]*domain.Branch, error) {
	var dbBranches []BranchInfo

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}

	// Count total for pagination
	var total int64
	if err := r.db.Model(&BranchInfo{}).
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
			dbBranch.BranchID,
			dbBranch.TenantID,
			"dbBranch.AuthBranchID",
			"dbBranch.AdminBranchID",
			dbBranch.BranchName,
			dbBranch.Description,
			dbBranch.UpdatedAt,
			nil,
		)
	}
	return branches, nil
}

func (r *branchRepository) Remove(ctx context.Context, tenantID, branchID string) error {

	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", branchID, tenantID).
		Delete(&BranchInfo{})

	if result.RowsAffected == 0 {
		return domain.ErrBranchNotFound
	}
	return result.Error
}
