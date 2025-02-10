package tenanthandler

import (
	"backend/internal/features/tenant/domain"
)

type CreateBranchRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateBranchRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type BranchResponse struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenantId"`
	AuthBranchID string `json:"authBranchId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	UpdatedAt    int64  `json:"updatedAt"`
}

func ToBranchResponse(b *domain.Branch) BranchResponse {
	return BranchResponse{
		ID:           b.ID(),
		TenantID:     b.TenantID(),
		AuthBranchID: b.AuthBranchID(),
		Name:         b.Name(),
		Description:  b.Description(),
		UpdatedAt:    b.UpdatedAt().Unix(),
	}
}

func ToBranchResponses(branches []*domain.Branch) []BranchResponse {
	responses := make([]BranchResponse, len(branches))
	for i, b := range branches {
		responses[i] = ToBranchResponse(b)
	}
	return responses
}
