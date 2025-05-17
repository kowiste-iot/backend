package repository

import (
	"backend/internal/features/measure/domain"
	"time"

	"gorm.io/gorm"
)

type MeasureDB struct {
	ID          string `gorm:"primaryKey"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (MeasureDB) TableName() string {
	return "measures"
}

// Mapper domain model to db model
func measureMapper(data *domain.Measure) (dbMeasure *MeasureDB) {
	dbMeasure = &MeasureDB{
		ID:          data.ID(),
		Parent:      data.Parent(),
		Name:        data.Name(),
		Description: data.Description(),
	}
	return
}
