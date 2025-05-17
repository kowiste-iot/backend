package domain

import (
	"backend/internal/features/user/domain/command"
	baseCmd "backend/shared/base/command"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, input *User) error
	Update(ctx context.Context, input *User) error
	FindByID(ctx context.Context, input *command.UserIDInput) (*User, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*User, error)
	Remove(ctx context.Context, input *command.UserIDInput) error
}
type IdentityProvider interface {
	CreateUser(ctx context.Context, input *command.CreateUserInput) (string, error)
	UpdateUser(ctx context.Context, input *command.UpdateUserInput) error
	DeleteUser(ctx context.Context, input *command.UserIDInput) error
	GetUser(ctx context.Context, input *command.UserIDInput) (*User, error)
	GetUserRoles(ctx context.Context, input *command.UserRolesInput) ([]Role, error)
	AssignRolesToUser(ctx context.Context, input *command.AssignRolesInput) error
	RemoveRolesFromUser(ctx context.Context, input *command.RemoveRolesInput) error
}
