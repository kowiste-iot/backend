package actionhandler

import (
	"backend/internal/features/action/domain"
)

type CreateActionRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type UpdateActionRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type ActionResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Parent      string `json:"parent,omitempty"`
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func ToActionResponse(a *domain.Action) ActionResponse {
	return ActionResponse{
		ID:          a.ID(),
		TenantID:    a.TenantID(),
		Name:        a.Name(),
		Parent:      a.Parent(),
		Enabled:     a.Enabled(),
		Description: a.Description(),
		UpdatedAt:   a.UpdatedAt().Unix(),
	}
}

func ToActionResponses(assets []*domain.Action) []ActionResponse {
	responses := make([]ActionResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToActionResponse(a)
	}
	return responses
}
