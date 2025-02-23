
package dashboardhandler

import (
	"backend/internal/features/dashboard/domain"
)

type CreateWidgetRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type UpdateWidgetRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type WidgetResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Parent      string `json:"parent,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func ToWidgetResponse(a *domain.Widget) WidgetResponse {
	return WidgetResponse{
		ID:          a.ID(),
		TenantID:    a.TenantID(),
		Name:        a.Name(),
		UpdatedAt:   a.UpdatedAt().Unix(),
	}
}

func ToWidgetResponses(assets []*domain.Widget) []WidgetResponse {
	responses := make([]WidgetResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToWidgetResponse(a)
	}
	return responses
}
