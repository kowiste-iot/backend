package base

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/http/httputil"

	"ddd/shared/logger"
	"errors"
	"fmt"
)

type BaseService struct {
	Logger logger.Logger
	Auth   auth.AuthProvider
}

func (b *BaseService) CheckPermission(c context.Context, resource string, scope string) (err error) {
	token, ok := httputil.GetToken(c)
	if !ok {
		return errors.New("token not found")
	}
	branch, ok := httputil.GetBranch(c)
	if !ok {
		return errors.New("branch not found")
	}
	branchClient := fmt.Sprintf("%s-service", branch)
	hasPermission, err := b.Auth.ValidatePermissionUser(c, token, branchClient, resource, scope)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("insufficient permissions")
	}
	return nil
}
