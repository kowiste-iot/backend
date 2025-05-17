package repository

import (
	"backend/internal/features/alert/domain"
	"time"

	"gorm.io/gorm"
)

type AlertDB struct {
	ID          string `gorm:"primaryKey"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Enabled     bool
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (AlertDB) TableName() string {
	return "alerts"
}

// Mapper domain model to db model
func alertMapper(data *domain.Alert) (dbDashboard *AlertDB) {
	dbDashboard = &AlertDB{
		ID:          data.ID(),
		Parent:      data.Parent(),
		Name:        data.Name(),
		Enabled:     data.Enabled(),
		Description: data.Description(),
	}
	return
}
