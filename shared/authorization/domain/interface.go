package domain

import (
	"context"
	"ddd/shared/authorization/domain/command"
)

type PermissionProvider interface {
	HasPermission(ctx context.Context, input *command.PermissionInput) (bool, error)
}
