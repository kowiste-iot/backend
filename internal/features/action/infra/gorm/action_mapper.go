package repository

import (
	"backend/internal/features/action/domain"
	"time"

	"gorm.io/gorm"
)

type ActionDB struct {
	ID          string `gorm:"primaryKey"`
	Parent      string `gorm:"type:string;references:ID"`
	Name        string
	Enabled     bool
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ActionDB) TableName() string {
	return "actions"
}

// Mapper domain model to db model
func actiondMapper(data *domain.Action) (dbBranch *ActionDB) {
	dbBranch = &ActionDB{
		ID:          data.ID(),
		Parent:      data.Parent(),
		Name:        data.Name(),
		Description: data.Description(),
	}
	return
}
