package scopehandler

import (
	"ddd/shared/auth/domain/scope"
)

// Responses
type ScopeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// Conversion helpers
func ToScopeResponse(sc scope.Scope) ScopeResponse {
	return ScopeResponse{
		ID:          sc.ID,
		Name:        sc.Name,
		DisplayName: sc.DisplayName,
	}
}

func ToScopeResponses(roles []scope.Scope) []ScopeResponse {
	responses := make([]ScopeResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToScopeResponse(role)
	}
	return responses
}
