// app/user_service.go
package app

import (
	"context"
	"ddd/internal/features/user/domain"
	"ddd/internal/features/user/domain/command"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/resource"
	"ddd/shared/auth/domain/scope"
	"ddd/shared/base"
	baseCmd "ddd/shared/base/command"

	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, input *command.CreateUserInput) (*domain.User, error)
	GetUser(ctx context.Context, input *command.UserIDInput) (*domain.User, error)
	ListUsers(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.User, error)
	UpdateUser(ctx context.Context, input *command.UpdateUserInput) (*domain.User, error)
	DeleteUser(ctx context.Context, input *command.UserIDInput) error
}

type userService struct {
	repo domain.UserRepository
	auth auth.IdentityProvider
	*base.BaseService
}

func NewService(base *base.BaseService, auth auth.IdentityProvider, repo domain.UserRepository) UserService {
	return &userService{
		repo:        repo,
		auth:        auth,
		BaseService: base,
	}
}

func (s *userService) CreateUser(ctx context.Context, input *command.CreateUserInput) (user *domain.User, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.User,
		Scope:     scope.Create,
	})
	if err != nil {
		return
	}
	user, err = domain.New(input.TenantDomain, input.BranchName, input.Email, input.FirstName, input.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	//user keycloak
	id, err := s.auth.CreateUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to save user auth: %w", err)
	}
	user.SetAuthID(id)
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, input *command.UserIDInput) (*domain.User, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.User,
		Scope:     scope.View,
	})
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindByID(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.User, error) {
	users, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, input *command.UpdateUserInput) (*domain.User, error) {
	// err := s.CheckPermission(ctx, resource.Asset, scope.Update)
	// if err != nil {
	// 	return nil, err
	// }
	user, err := s.repo.FindByID(ctx, &command.UserIDInput{
		BaseInput: baseCmd.BaseInput{
			TenantDomain: input.TenantDomain,
			BranchName:   input.BranchName,
		},
		UserID: input.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = user.Update(input.Email, input.FirstName, input.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if err := s.repo.Update(ctx, input); err != nil {
		return nil, fmt.Errorf("failed to save updated user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, input *command.UserIDInput) error {
	err := s.repo.Remove(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
