package repository

import (
	"context"
	"ddd/internal/features/device/domain"
	baseCmd "ddd/shared/base/command"
	gormhelper "ddd/shared/gorm"
	"ddd/shared/pagination"
	"errors"
	"time"

	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

type Device struct {
	ID          string `gorm:"primaryKey"`
	TenantID    string `gorm:"index"`
	BranchID    string `gorm:"index"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewRepository(db *gorm.DB) domain.DeviceRepository {
	db.AutoMigrate(&Device{})
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, input *domain.Device) error {
	dbDevice := Device{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Create(&dbDevice).Error
}

func (r *deviceRepository) Update(ctx context.Context, input *domain.Device) error {
	dbDevice := Device{
		ID:          input.ID(),
		TenantID:    input.TenantID(),
		BranchID:    input.BranchName(),
		Parent:      input.Parent(),
		Name:        input.Name(),
		Description: input.Description(),
	}
	return r.db.WithContext(ctx).Updates(&dbDevice).Error
}

func (r *deviceRepository) FindByID(ctx context.Context, input *baseCmd.BaseInput, deviceID string) (*domain.Device, error) {
	var dbDevice Device

	err := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", deviceID).
		First(&dbDevice).Error
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbDevice.ID,
		dbDevice.TenantID,
		dbDevice.BranchID,
		dbDevice.Name,
		dbDevice.Description,
		dbDevice.Parent,
		dbDevice.UpdatedAt,
		&dbDevice.DeletedAt.Time,
	), nil

}

func (r *deviceRepository) FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Device, error) {
	var dbDevices []Device

	pg, ok := pagination.GetPagination(ctx)
	if !ok {
		return nil, errors.New("pagination not found in context")
	}
	var total int64
	err := r.db.Model(&Device{}).Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = r.db.WithContext(ctx).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbDevices).Error
	if err != nil {
		return nil, err
	}

	devices := make([]*domain.Device, len(dbDevices))
	for i, dbDevice := range dbDevices {
		devices[i] = domain.NewFromRepository(
			dbDevice.ID,
			dbDevice.TenantID,
			dbDevice.BranchID,
			dbDevice.Name,
			dbDevice.Description,
			dbDevice.Parent,
			dbDevice.UpdatedAt,
			&dbDevice.DeletedAt.Time,
		)
	}
	return devices, nil
}

func (r *deviceRepository) Remove(ctx context.Context, input *baseCmd.BaseInput, deviceID string) error {

	resp := r.db.WithContext(ctx).Where(
		gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND id = ?", deviceID).Delete(&Device{})
	if resp.RowsAffected == 0 {
		return errors.New("no delete")
	}
	return resp.Error
}

func (r *deviceRepository) HasChildren(ctx context.Context, input *baseCmd.BaseInput, deviceID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Device{}).
		Where(gormhelper.TenantBranchFilter(input.TenantDomain, input.BranchName)+" AND parent = ?", deviceID).
		Count(&count).Error

	return count > 0, err
}
