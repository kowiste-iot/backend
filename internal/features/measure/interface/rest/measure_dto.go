package measurehandler

import (
	"backend/internal/features/measure/domain"
)

type CreateMeasureRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type UpdateMeasureRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type MeasureResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Parent      string `json:"parent,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func ToMeasureResponse(a *domain.Measure) MeasureResponse {
	return MeasureResponse{
		ID:          a.ID(),
		TenantID:    a.TenantID(),
		Name:        a.Name(),
		Parent:      a.Parent(),
		Description: a.Description(),
		UpdatedAt:   a.UpdatedAt().Unix(),
	}
}

func ToMeasureResponses(assets []*domain.Measure) []MeasureResponse {
	responses := make([]MeasureResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToMeasureResponse(a)
	}
	return responses
}
