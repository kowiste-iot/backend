package resourcehandler

import (
	"ddd/shared/auth/domain/resource"
)

type ResourceResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	DisplayName string              `json:"displayName"`
	Roles       map[string][]string `json:"roles"`
}

// Conversion helpers
func ToResourcesResponse(resource resource.ResourcePermission) ResourceResponse {
	return ResourceResponse{
		ID:          resource.ID,
		Name:        resource.Name,
		DisplayName: resource.DisplayName,
		Roles:       resource.Roles,
	}
}

func ToResourcesResponses(resources []resource.ResourcePermission) []ResourceResponse {
	responses := make([]ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = ToResourcesResponse(resource)
	}
	return responses
}
