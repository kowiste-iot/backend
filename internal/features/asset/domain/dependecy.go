package domain

import "context"

type AssetDependency struct {
	TenantID  string
	BranchID  string
	FeatureID string
	Feature   string
	AssetID   string
}

func NewAssetDependency(tenantID, branchID, featureID, feature, assetID string) *AssetDependency {
	return &AssetDependency{
		TenantID:  tenantID,
		BranchID:  branchID,
		FeatureID: featureID,
		Feature:   feature,
		AssetID:   assetID,
	}
}

func (d *AssetDependency) UpdateAsset(newAssetID string) {
	d.AssetID = newAssetID
}

type AssetDependencyRepository interface {
	Create(ctx context.Context, dependency *AssetDependency) error
	Update(ctx context.Context, dependency *AssetDependency) error
	Remove(ctx context.Context, tenantID, branchID, featureID string) error
	FindByFeatureID(ctx context.Context, tenantID, branchID, featureID string) (*AssetDependency, error)
	FindByAssetID(ctx context.Context, tenantID, branchID, assetID string) ([]*AssetDependency, error)
}
