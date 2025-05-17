package repository

import (
	"backend/internal/features/tenant/domain"
	"time"

	"gorm.io/gorm"
)

type TenantDB struct {
	ID          string `gorm:"primaryKey"`
	AuthID      string `gorm:"index"`
	Name        string
	Domain      string `gorm:"index"`
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (b *TenantDB) TableName() string {
	return "platform.tenants"
}

// Mapper domain model to db model
func tenantMapper(data *domain.Tenant) (dbTenant *TenantDB) {
	dbTenant = &TenantDB{
		ID:          data.ID(),
		AuthID:      data.AuhtID(),
		Domain:      data.Domain(),
		Name:        data.Name(),
		Description: data.Description(),
	}
	return
}
