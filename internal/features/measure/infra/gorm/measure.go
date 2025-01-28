package repository

import (
	"context"
	"ddd/internal/features/measure/domain"
	baseCmd "ddd/shared/base/command"
	gormhelper "ddd/shared/gorm"
	"ddd/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type measureRepository struct {
	db *gorm.DB
}

type Measure struct {
	ID          string `gorm:"primaryKey"`
	TenantID    string `gorm:"index"`
	BranchID    string `gorm:"index"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewRepository(db *gorm.DB) domain.MeasureRepository {
	db.AutoMigrate(&Measure{})
	return &measureRepository{db: db}
}

func (r *measureRepository) Create(ctx context.Context, input *domain.Measure) error {
	dbMeasure := Measure{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Create(&dbMeasure).Error
}

func (r *measureRepository) Update(ctx context.Context, input *domain.Measure) error {
	dbMeasure := Measure{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Updates(&dbMeasure).Error
}

func (r *measureRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, measureID string) (*domain.Measure, error) {
	var dbMeasure Measure

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", measureID).
		First(&dbMeasure).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbMeasure.ID,
		dbMeasure.TenantID,
		dbMeasure.BranchID,
		dbMeasure.Name,
		dbMeasure.Description,
		dbMeasure.Parent,
		dbMeasure.UpdatedAt,
		&dbMeasure.DeletedAt.Time,
	), nil

}

func (r *measureRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Measure, error) {
	var dbMeasures []Measure

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Measure{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
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
			dbMeasure.TenantID,
			dbMeasure.BranchID,
			dbMeasure.Name,
			dbMeasure.Description,
			dbMeasure.Parent,
			dbMeasure.UpdatedAt,
			&dbMeasure.DeletedAt.Time,
		)
	}
	return measures, nil
}

func (r *measureRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, measureID string) error {

	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", measureID).Delete(&Measure{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *measureRepository) HasChildren(ctx context.Context, input *baseCmd.BaseInput, measureID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Measure{}).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND parent = ?", measureID).
		Count(&count).Error

	return count > 0, err
}
