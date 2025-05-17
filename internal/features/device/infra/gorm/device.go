// internal/features/device/infra/gorm/device.go
package repository

import (
	"backend/internal/features/device/domain"
	mqttRepo "backend/internal/features/device/infra/gorm/mqtt"
	gormhelper "backend/shared/gorm"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"backend/shared/util"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type deviceRepository struct {
	db       *gorm.DB
	mqttRepo mqttRepo.MqttRepository
}

func NewRepository(db *gorm.DB) domain.DeviceRepository {
	mqtt := mqttRepo.NewRepository(db)
	return &deviceRepository{
		db:       db,
		mqttRepo: mqtt,
	}
}

func (r *deviceRepository) Create(ctx context.Context, input *domain.Device) (password string, err error) {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the connection for the specific branch
	txWithBranch, err := gormhelper.SetBranchDB(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to set branch DB: %w", err)
	}

	// Create device in platform schema
	if err := txWithBranch.WithContext(ctx).Create(deviceMapper(input)).Error; err != nil {
		return "", err
	}

	// Generate a secure random password for the device
	password, err = util.GenerateSecurePassword(16) // 16 bytes = 128 bits
	if err != nil {
		return "", fmt.Errorf("failed to generate device password: %w", err)
	}

	// Create MQTT user
	if err := r.mqttRepo.CreateMqttUser(ctx, input.ID(), password); err != nil {
		return "", err
	}

	return password, tx.Commit().Error
}

func (r *deviceRepository) Update(ctx context.Context, input *domain.Device) error {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return fmt.Errorf("failed to set branch DB: %w", err)
	}
	return db.WithContext(ctx).Updates(deviceMapper(input)).Error
}

func (r *deviceRepository) FindByID(ctx context.Context, deviceID string) (*domain.Device, error) {
	var dbDevice DeviceDB
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to set branch DB: %w", err)
	}
	err = db.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", deviceID).
		First(&dbDevice).Error
	if err != nil {
		return nil, err
	}

	// Extract branch and tenant info from context
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		return nil, err
	}
	return domain.NewFromRepository(
		dbDevice.ID,
		tenant.Domain(),
		branch,
		dbDevice.Name,
		dbDevice.Description,
		dbDevice.Parent,
		dbDevice.UpdatedAt,
		&dbDevice.DeletedAt.Time,
	), nil
}

func (r *deviceRepository) FindAll(ctx context.Context) ([]*domain.Device, error) {
	var dbDevices []DeviceDB

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
	err = db.Model(&DeviceDB{}).Where(gormhelper.DeleteFilter()).Count(&total).Error
	if err != nil {
		return nil, err
	}
	pg.Total = total
	ctx = pagination.WithPagination(ctx, pg)

	err = db.WithContext(ctx).
		Where(gormhelper.DeleteFilter()).
		Offset(pg.Offset).
		Limit(pg.PageSize).
		Find(&dbDevices).Error
	if err != nil {
		return nil, err
	}

	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		return nil, err
	}
	devices := make([]*domain.Device, len(dbDevices))
	for i, dbDevice := range dbDevices {
		devices[i] = domain.NewFromRepository(
			dbDevice.ID,
			tenant.Domain(),
			branch,
			dbDevice.Name,
			dbDevice.Description,
			dbDevice.Parent,
			dbDevice.UpdatedAt,
			&dbDevice.DeletedAt.Time,
		)
	}
	return devices, nil
}

func (r *deviceRepository) Remove(ctx context.Context, deviceID string) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the connection for the specific branch
	txWithBranch, err := gormhelper.SetBranchDB(ctx, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set branch DB: %w", err)
	}

	// Delete device from platform schema
	resp := txWithBranch.WithContext(ctx).Where(
		gormhelper.DeleteFilter()+" AND id = ?", deviceID).Delete(&DeviceDB{})

	if resp.Error != nil {
		tx.Rollback()
		return resp.Error
	}

	if resp.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("no delete")
	}

	// Delete MQTT user
	if err := r.mqttRepo.DeleteMqttUser(ctx, deviceID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *deviceRepository) HasChildren(ctx context.Context, deviceID string) (bool, error) {
	// Get the connection for the specific branch
	db, err := gormhelper.SetBranchDB(ctx, r.db)
	if err != nil {
		return false, fmt.Errorf("failed to set branch DB: %w", err)
	}
	var count int64
	err = db.WithContext(ctx).
		Model(&DeviceDB{}).
		Where(gormhelper.DeleteFilter()+" AND parent = ?", deviceID).
		Count(&count).Error

	return count > 0, err
}
