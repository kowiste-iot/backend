package rolehandler

import (
	auth "ddd/shared/auth/domain"
)

// Requests
type CreateRoleRequest struct {
	Role        string `json:"role" binding:"required,min=3,max=255"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Role        string `json:"role" binding:"required,min=3,max=255"`
	Description string `json:"description"`
}

type AssignRoleRequest struct {
	UserID string   `json:"userId" binding:"required,uuid"`
	Roles  []string `json:"roles" binding:"required"`
}

// Responses
type RoleResponse struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type RoleAssignmentResponse struct {
	UserID string   `json:"userId"`
	Roles  []string `json:"roles"`
}

// Conversion helpers
func ToRoleResponse(role *auth.Role) RoleResponse {
	return RoleResponse{
		Name:        role.Name,
		Description: role.Description,
	}
}

func ToRoleResponses(roles []*auth.Role) []RoleResponse {
	responses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToRoleResponse(role)
	}
	return responses
}
