package base

import (
	"context"
	auth "ddd/shared/auth/domain"
	authCmd "ddd/shared/auth/domain/command"
	"ddd/shared/base/command"
	"ddd/shared/http/httputil"
	"fmt"

	"ddd/shared/logger"
	"errors"
)

type BaseService struct {
	Logger logger.Logger
	Auth   auth.AuthProvider
}

func (b *BaseService) CheckPermission(ctx context.Context, input *command.CheckPermissionInput) (err error) {
	token, ok := httputil.GetToken(ctx)
	if !ok {
		return fmt.Errorf("not token present")
	}
	hasPermission, err := b.Auth.ValidatePermissionUser(ctx, token,
		authCmd.ClientName(input.BranchName), input.Resource, input.Scope)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("insufficient permissions")
	}
	return nil
}
