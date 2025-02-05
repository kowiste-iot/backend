package resourcehandler

import (
	"backend/internal/features/resource/domain"
	scopeDomain "backend/internal/features/scope/domain"
)

type UpdateResourceRequest struct {
	ID          string                         `json:"id" binding:"required"`
	Name        string                         `json:"name" binding:"required"`
	DisplayName string                         `json:"displayname"`
	Roles       map[string][]scopeDomain.Scope `json:"roles"`
}
type ResourceResponse struct {
	ID          string                         `json:"id"`
	Name        string                         `json:"name"`
	DisplayName string                         `json:"displayName"`
	Roles       map[string][]scopeDomain.Scope `json:"roles"`
}

// Conversion helpers
func ToResourcesResponse(resource domain.ResourcePermission) ResourceResponse {
	return ResourceResponse{
		ID:          resource.ID,
		Name:        resource.Name,
		DisplayName: resource.DisplayName,
		Roles:       resource.Roles,
	}
}

func ToResourcesResponses(resources []domain.ResourcePermission) []ResourceResponse {
	responses := make([]ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = ToResourcesResponse(resource)
	}
	return responses
}
