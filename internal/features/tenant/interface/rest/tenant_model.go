package tenanthandler

import (
	"backend/internal/features/tenant/domain"
)

type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Branch      string `json:"branch"`
}

type UpdateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type TenantResponse struct {
	ID          string `json:"id"`
	AuthID      string `json:"authId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func ToTenantResponse(a *domain.Tenant) TenantResponse {
	return TenantResponse{
		ID:          a.ID(),
		AuthID:      a.AuhtID(),
		Name:        a.Name(),
		Description: a.Description(),
		UpdatedAt:   a.UpdatedAt().Unix(),
	}
}

func ToTenantResponses(tenants []*domain.Tenant) []TenantResponse {
	responses := make([]TenantResponse, len(tenants))
	for i, a := range tenants {
		responses[i] = ToTenantResponse(a)
	}
	return responses
}
