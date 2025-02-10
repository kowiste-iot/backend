package assethandler

import (
	"backend/internal/features/asset/domain"
)

type CreateAssetRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type UpdateAssetRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type AssetResponse struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenantId"`
	Parent      *string `json:"parent,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UpdatedAt   int64   `json:"updatedAt"`
}

func ToAssetResponse(a *domain.Asset) AssetResponse {
	return AssetResponse{
		ID:          a.ID(),
		TenantID:    a.TenantID(),
		Name:        a.Name(),
		Parent:      a.Parent(),
		Description: a.Description(),
		UpdatedAt:   a.UpdatedAt().Unix(),
	}
}

func ToAssetResponses(assets []*domain.Asset) []AssetResponse {
	responses := make([]AssetResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToAssetResponse(a)
	}
	return responses
}
