package dashboardhandler

import (
	"backend/internal/features/dashboard/domain"
)

type CreateDashboardRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type UpdateDashboardRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type DashboardResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Parent      string `json:"parent,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func ToDashboardResponse(a *domain.Dashboard) DashboardResponse {
	return DashboardResponse{
		ID:          a.ID(),
		TenantID:    a.TenantID(),
		Name:        a.Name(),
		Parent:      a.Parent(),
		Description: a.Description(),
		UpdatedAt:   a.UpdatedAt().Unix(),
	}
}

func ToDashboardResponses(assets []*domain.Dashboard) []DashboardResponse {
	responses := make([]DashboardResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToDashboardResponse(a)
	}
	return responses
}
