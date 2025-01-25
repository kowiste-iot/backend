package rolehandler

import (
	"ddd/shared/auth/domain/role"
)

// Requests
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description"`
}

type AssignRoleRequest struct {
	UserID string   `json:"userId" binding:"required,uuid"`
	Roles  []string `json:"roles" binding:"required"`
}

// Responses
type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ReadOnly    bool   `json:"readonly"`
	Description string `json:"description,omitempty"`
}

type RoleAssignmentResponse struct {
	UserID string   `json:"userId"`
	Roles  []string `json:"roles"`
}

// Conversion helpers
func ToRoleResponse(role role.Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		ReadOnly:    role.IsDefaultRole(),
		Description: role.Description,
	}
}

func ToRoleResponses(roles []role.Role) []RoleResponse {
	responses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToRoleResponse(role)
	}
	return responses
}
