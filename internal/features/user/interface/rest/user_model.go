package userhandler

import "backend/internal/features/user/domain"

type CreateUserRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	FirstName string   `json:"firstName" binding:"required"`
	LastName  string   `json:"lastName" binding:"required"`
	Roles     []string `json:"roles" binding:"required"`
}

type UpdateUserRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	FirstName string   `json:"firstName" binding:"required"`
	LastName  string   `json:"lastName" binding:"required"`
	Roles     []string `json:"roles" binding:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	AuthID    string `json:"authId"`
	Email     string `json:"email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	UpdatedAt int64  `json:"updatedAt"`
}

func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID(),
		AuthID:    u.AuthID(),
		Email:     u.Email(),
		FirstName: u.FirstName(),
		LastName:  u.LastName(),
		UpdatedAt: u.UpdatedAt().Unix(),
	}
}

func ToUserResponses(users []*domain.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponse(u)
	}
	return responses
}
