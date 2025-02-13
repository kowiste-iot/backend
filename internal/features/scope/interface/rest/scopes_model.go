package scopehandler

import "backend/internal/features/scope/domain"

// Responses
type ScopeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// Conversion helpers
func ToScopeResponse(sc domain.Scope) ScopeResponse {
	return ScopeResponse{
		ID:          sc.ID,
		Name:        sc.Name,
		DisplayName: sc.DisplayName,
	}
}

func ToScopeResponses(roles []domain.Scope) []ScopeResponse {
	responses := make([]ScopeResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToScopeResponse(role)
	}
	return responses
}
