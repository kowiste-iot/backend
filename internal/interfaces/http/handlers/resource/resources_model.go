package resourcehandler

import (
	"ddd/shared/auth/domain/resource"
)

type ResourceResponse struct {
	Name        string `json:"name"`
	DisplayName string `json:"displaName"`
}

// Conversion helpers
func ToResourcesResponse(resource resource.Resource) ResourceResponse {
	return ResourceResponse{
		Name:        resource.Name,
		DisplayName: resource.DisplayName,
	}
}

func ToResourcesResponses(resources []resource.Resource) []ResourceResponse {
	responses := make([]ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = ToResourcesResponse(resource)
	}
	return responses
}
