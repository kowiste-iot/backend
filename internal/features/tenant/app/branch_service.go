package app

import (
	"context"
	"ddd/internal/features/tenant/domain"
	"ddd/internal/features/tenant/domain/command"
	appAuth "ddd/shared/auth/app"
	authCmd "ddd/shared/auth/domain/command"
	"ddd/shared/base"
	baseCmd "ddd/shared/base/command"
	"ddd/shared/validator"
	"fmt"
)

type BranchService interface {
	CreateBranch(ctx context.Context, input *command.CreateBranchInput) (*domain.Branch, error)
	GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*domain.Branch, error)
	ListBranches(ctx context.Context, tenantID string) ([]*domain.Branch, error)
	UpdateBranch(ctx context.Context, input *command.UpdateBranchInput) (*domain.Branch, error)
	DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error
}

type branchService struct {
	repo domain.BranchRepository
	auth *appAuth.Service
	*base.BaseService
}

func NewBranchService(base *base.BaseService, auth *appAuth.Service, repo domain.BranchRepository) BranchService {
	return &branchService{
		repo:        repo,
		auth:        auth,
		BaseService: base,
	}
}

func (s *branchService) CreateBranch(ctx context.Context, input *command.CreateBranchInput) (*domain.Branch, error) {
	// err := s.CheckPermission(ctx, resource.Tenant, scope.Create)
	// if err != nil {
	// 	return nil, err
	// }
	err := validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	branch, err := domain.NewBranch(input.TenantDomain, input.Name, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	// Create Keycloak group
	branchID, err := s.auth.CreateBranch(ctx, &authCmd.CreateBranchInput{
		TenantID:    input.TenantDomain,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auth branch: %w", err)
	}

	branch.SetAuthBranchID(branchID)

	if err := s.repo.Create(ctx, input.TenantDomain, branch); err != nil {
		// Cleanup Keycloak group if DB save fails
		i := baseCmd.NewInput(input.TenantDomain, branchID)
		if delErr := s.auth.DeleteBranch(ctx, &i); delErr != nil {
			s.Logger.Error(ctx, err, "failed to cleanup auth group after branch creation failure", nil)
		}
		return nil, fmt.Errorf("failed to save branch: %w", err)
	}

	return branch, nil
}

func (s *branchService) GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*domain.Branch, error) {

	branch, err := s.repo.FindByID(ctx, input.TenantDomain, input.BranchName)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	return branch, nil
}

func (s *branchService) ListBranches(ctx context.Context, tenantID string) ([]*domain.Branch, error) {

	branches, err := s.repo.FindAll(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}
	return branches, nil
}

func (s *branchService) UpdateBranch(ctx context.Context, cmd *command.UpdateBranchInput) (*domain.Branch, error) {

	branch, err := s.repo.FindByID(ctx, cmd.TenantDomain, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	err = branch.Update(cmd.Name, cmd.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}

	// Update Keycloak group
	err = s.auth.UpdateBranch(ctx, &authCmd.UpdateBranchInput{
		TenantID:    cmd.TenantDomain,
		ID:          cmd.ID,
		Name:        cmd.Name,
		Description: cmd.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update auth group: %w", err)
	}

	if err := s.repo.Update(ctx, cmd.TenantDomain, branch); err != nil {
		return nil, fmt.Errorf("failed to save branch: %w", err)
	}

	return branch, nil
}

func (s *branchService) DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error {

	_, err := s.repo.FindByID(ctx, input.TenantDomain, input.BranchName)
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}

	err = s.auth.DeleteBranch(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete auth branch: %w", err)
	}

	if err := s.repo.Remove(ctx, input.TenantDomain, input.BranchName); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}
