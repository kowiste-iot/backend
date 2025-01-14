package domain

import (
	"context"
	"ddd/internal/features/user/domain/command"
	baseCmd "ddd/shared/base/command"
)

type UserRepository interface {
	Create(ctx context.Context, input *User) error
	Update(ctx context.Context, input *command.UpdateUserInput) error
	FindByID(ctx context.Context, input *command.UserIDInput) (*User, error)
	FindAll(ctx context.Context, input *baseCmd.BaseInput) ([]*User, error)
	Remove(ctx context.Context, input *command.UserIDInput) error
}
