package domain

import (
	"context"
	"ddd/shared/authentication/domain/command"
	"ddd/shared/authorization/domain"
)

type IdentityProvider interface {
	CreateUser(ctx context.Context, input *command.CreateUserInput) (string, error)
	UpdateUser(ctx context.Context, input *command.UpdateUserInput) error
	DeleteUser(ctx context.Context, input *command.UserIDInput) error
	GetUser(ctx context.Context, input *command.UserIDInput) (*User, error)
	GetUserRoles(ctx context.Context, input *command.UserRolesInput) ([]domain.Role, error)
	AssignRolesToUser(ctx context.Context, input *command.AssignRolesInput) error
	RemoveRolesFromUser(ctx context.Context, input *command.RemoveRolesInput) error
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
