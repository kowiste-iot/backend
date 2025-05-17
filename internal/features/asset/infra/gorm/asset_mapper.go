package repository

import (
	"backend/internal/features/asset/domain"
	"time"

	"gorm.io/gorm"
)

type AssetDB struct {
	ID          string  `gorm:"primaryKey"`
	Parent      *string `gorm:"type:string;references:ID"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (AssetDB) TableName() string {
	return "assets"
}

// Mapper domain model to db model
func assetMapper( data *domain.Asset) (dbBranch *AssetDB) {
	dbBranch = &AssetDB{
		ID:          data.ID(),
		Parent:      data.Parent(),
		Name:        data.Name(),
		Description: data.Description(),
	}
	return
}
