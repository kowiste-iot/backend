package resourcehandler

import (
	"ddd/shared/auth/domain/resource"
)

type ResourceResponse struct {
	Name        string `json:"name"`
	ReadOnly    bool   `json:"readonly"`
	Description string `json:"description,omitempty"`
}


// Conversion helpers
func ToResourcesResponse(role resource.Resource) ResourceResponse {
	return ResourceResponse{
		Name:        role.Name,
	}
}

func ToResourcesResponses(resources []resource.Resource) []ResourceResponse {
	responses := make([]ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = ToResourcesResponse(resource)
	}
	return responses
}
