package auth

import (
	userCmd "backend/internal/features/user/domain/command"
	"context"
)

type IdentityProvider interface {
	CreateUser(ctx context.Context, input *userCmd.CreateUserInput) (string, error)
	UpdateUser(ctx context.Context, input *userCmd.UpdateUserInput) error
	DeleteUser(ctx context.Context, input *userCmd.UserIDInput) error
	GetUser(ctx context.Context, input *userCmd.UserIDInput) (*User, error)
}

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	TenantID  string   `json:"tenantId"`
	Roles     []string `json:"roles"`
	AuthID    string   `json:"authId,omitempty"`
}
