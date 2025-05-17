package domain

import (
	"backend/shared/authorization/domain/command"
	"context"
)

type PermissionProvider interface {
	HasPermission(ctx context.Context, input *command.PermissionInput) (bool, error)
}
