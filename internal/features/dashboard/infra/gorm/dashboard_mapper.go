package repository

import (
	"backend/internal/features/dashboard/domain"
	"time"

	"gorm.io/gorm"
)

type DashboardDB struct {
	ID          string `gorm:"primaryKey"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (DashboardDB) TableName() string {
	return "dashboards"
}

// Mapper domain model to db model
func dashboardMapper(data *domain.Dashboard) (dbDashboard *DashboardDB) {
	dbDashboard = &DashboardDB{
		ID:          data.ID(),

		Parent:      data.Parent(),
		Name:        data.Name(),
		Description: data.Description(),
	}
	return
}
