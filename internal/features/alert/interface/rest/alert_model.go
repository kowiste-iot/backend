package alerthandler

import (
	"backend/internal/features/alert/domain"
)

type CreateAlertRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type UpdateAlertRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type AlertResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Parent      string `json:"parent,omitempty"`
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func ToAlertResponse(a *domain.Alert) AlertResponse {
	return AlertResponse{
		ID:          a.ID(),
		TenantID:    a.TenantID(),
		Name:        a.Name(),
		Parent:      a.Parent(),
		Enabled:     a.Enabled(),
		Description: a.Description(),
		UpdatedAt:   a.UpdatedAt().Unix(),
	}
}

func ToAlertResponses(assets []*domain.Alert) []AlertResponse {
	responses := make([]AlertResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToAlertResponse(a)
	}
	return responses
}
