package repository

import (
	"backend/internal/features/asset/domain"
	"time"

	"gorm.io/gorm"
)

type AssetDependencyDB struct {
	FeatureID string  `gorm:"uniqueIndex:idx_dependency"`
	Feature   string  `gorm:"uniqueIndex:idx_dependency"`
	AssetID   string  `gorm:"type:string;references:ID;constraint:OnDelete:RESTRICT"`
	Asset     AssetDB `gorm:"foreignKey:AssetID;references:ID"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (AssetDependencyDB) TableName() string {
	return "asset_dependencies"
}

// Mapper domain model to db model
func assetDependenciesMapper(data *domain.AssetDependency) (dbAssetDependecies *AssetDependencyDB) {
	dbAssetDependecies = &AssetDependencyDB{
		FeatureID: data.FeatureID,
		Feature:   string(data.Feature),
		AssetID:   data.AssetID,
	}
	return
}
