package base

import (
	authzDomain "backend/shared/authorization/domain"
	authzCmd "backend/shared/authorization/domain/command"
	"backend/shared/base/command"
	"backend/shared/http/httputil"
	"context"
	"fmt"

	"backend/shared/logger"
	"errors"

	"gorm.io/gorm"
)

type BaseService struct {
	Logger logger.Logger
	DB     *gorm.DB
	// Auth   validation.AuthProvider
	Perm authzDomain.PermissionProvider
}

func New(log logger.Logger, permission authzDomain.PermissionProvider, db *gorm.DB) *BaseService {
	return &BaseService{
		Logger: log,
		DB:     db,
		Perm:   permission,
	}
}

func (b *BaseService) CheckPermission(ctx context.Context, input *command.CheckPermissionInput) (err error) {
	token, ok := httputil.GetToken(ctx)
	if !ok {
		return fmt.Errorf("not token present")
	}
	hasPermission, err := b.Perm.HasPermission(ctx, &authzCmd.PermissionInput{
		Token:    token,
		Resource: input.Resource,
		Action:   input.Scope,
		TenantID: input.TenantDomain,
		BranchID: input.BranchName,
	})
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("insufficient permissions")
	}
	return nil
}
