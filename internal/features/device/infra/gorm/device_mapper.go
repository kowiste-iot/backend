package repository

import (
	"backend/internal/features/device/domain"
	"time"

	"gorm.io/gorm"
)

type DeviceDB struct {
	ID          string `gorm:"primaryKey"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (DeviceDB) TableName() string {
	return "devices"
}

// Mapper domain model to db model
func deviceMapper(data *domain.Device) (dbDevice *DeviceDB) {
	dbDevice = &DeviceDB{
		ID:          data.ID(),
		Parent:      data.Parent(),
		Name:        data.Name(),
		Description: data.Description(),
		UpdatedAt:   data.UpdatedAt(),
	}
	return
}
