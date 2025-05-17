package repository

import (
	"backend/internal/features/measure/domain"
	gormhelper "backend/shared/gorm"
	"backend/shared/pagination"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type measureRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.MeasureRepository {
	return &measureRepository{db: db}
}

func (r *measureRepository) Create(ctx context.Context, input *domain.Measure) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Create(measureMapper(input)).Error
}

func (r *measureRepository) Update(ctx context.Context, input *domain.Measure) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Updates(measureMapper(input)).Error
}

func (r *measureRepository) FindByID(ctx context.Context, measureID string) (*domain.Measure, error) {
	var dbMeasure MeasureDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", measureID).
		First(&dbMeasure).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbMeasure.ID,
		"dbMeasure.TenantID",
		"dbMeasure.BranchID",
		dbMeasure.Name,
		dbMeasure.Description,
		dbMeasure.Parent,
		dbMeasure.UpdatedAt,
		&dbMeasure.DeletedAt.Time,
	), nil

}

func (r *measureRepository) FindAll(ctx context.Context) ([]*domain.Measure, error) {
	var dbMeasures []MeasureDB

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
	err = db.Model(&MeasureDB{}).Where(gormhelper.DeleteFilter()).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbMeasures).Error
	if err != nil {
		return nil, err
	}

	measures := make([]*domain.Measure, len(dbMeasures))
	for i, dbMeasure := range dbMeasures {
		measures[i] = domain.NewFromRepository(
			dbMeasure.ID,
			"dbMeasure.TenantID",
			"dbMeasure.BranchID",
			dbMeasure.Name,
			dbMeasure.Description,
			dbMeasure.Parent,
			dbMeasure.UpdatedAt,
			&dbMeasure.DeletedAt.Time,
		)
	}
	return measures, nil
}

func (r *measureRepository) Remove(ctx context.Context, measureID string) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	resp := db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", measureID).Delete(&MeasureDB{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *measureRepository) HasChildren(ctx context.Context, measureID string) (bool, error) {
	var count int64
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return false, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).
		Model(&MeasureDB{}).
		Where(gormhelper.DeleteFilter()+" AND parent = ?", measureID).
		Count(&count).Error

	return count > 0, err
}
