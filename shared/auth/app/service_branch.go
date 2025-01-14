package app

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/command"
	baseCmd "ddd/shared/base/command"
)

func (s *Service) CreateBranch(ctx context.Context, input *command.CreateBranchInput) (id string, err error) {


	// Create branch in auth provider
	id, err = s.tenantProvider.CreateBranch(ctx, input)
	if err != nil {
		return
	}

	return
}

func (s *Service) UpdateBranch(ctx context.Context,input *command.UpdateBranchInput) error {
	return s.tenantProvider.UpdateBranch(ctx, input)
}

func (s *Service) DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error {
	return s.tenantProvider.DeleteBranch(ctx, input)
}

func (s *Service) GetBranch(ctx context.Context,  input *baseCmd.BaseInput) (*auth.Branch, error) {
	return s.tenantProvider.GetBranch(ctx, input)
}

func (s *Service) AssignUserToBranch(ctx context.Context,  input *command.UserToBranch) error {
	return s.tenantProvider.AssignUserToBranch(ctx, input)
}

func (s *Service) RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error {
	return s.tenantProvider.RemoveUserFromBranch(ctx, input)
}

func (s *Service) GetBranchUsers(ctx context.Context, input *baseCmd.BaseInput) ([]auth.User, error) {
	return s.tenantProvider.GetBranchUsers(ctx, input)
}
