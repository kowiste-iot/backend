package auth

import (
	"context"
	userCmd "ddd/internal/features/user/domain/command"

	"github.com/golang-jwt/jwt/v5"
)

type AuthProvider interface {
	ValidateToken(ctx context.Context, token string) (*jwt.Token, error)
	ValidatePermissionService(ctx context.Context, token, clientID, resource, scope string) (bool, error)
	ValidatePermissionUser(ctx context.Context, token, clientID, resource, scope string) (bool, error)
}

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
