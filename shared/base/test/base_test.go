package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/internal/features/resource/domain"
	auhtMock "backend/shared/authorization/test/mock"
	"backend/shared/base"
	"backend/shared/base/command"
	"backend/shared/http/httputil"
	"backend/shared/logger/test/mock"
)

func TestBaseService_CheckPermission(t *testing.T) {
	ctx := context.Background()
	bInput := command.BaseInput{
		TenantDomain: "madrid",
		BranchName:   "seat",
	}

	baseService := base.New(mock.NewMockLogger(), auhtMock.NewMockPermissionProvider())
	err := baseService.CheckPermission(ctx, &command.CheckPermissionInput{})
	assert.Error(t, err)

	ctx = httputil.SetToken(ctx, "ye4ygbthwhrtrt")
	err = baseService.CheckPermission(ctx, &command.CheckPermissionInput{
		BaseInput: bInput,
		Resource:  domain.Asset,
		Scope:     "view",
	})
	assert.Error(t, err)

}
