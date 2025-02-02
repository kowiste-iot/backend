// app/user_service.go
package app

import (
	"context"

	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
	"backend/shared/auth/domain/resource"
	"backend/shared/auth/domain/scope"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"

	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, input *command.CreateUserInput) (*domain.User, error)
	GetUser(ctx context.Context, input *command.UserIDInput) (*domain.User, error)
	ListUsers(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.User, error)
	UpdateUser(ctx context.Context, input *command.UpdateUserInput) (*domain.User, error)
	DeleteUser(ctx context.Context, input *command.UserIDInput) error
}
type ServiceDependencies struct {
	Repo domain.UserRepository
	Auth domain.IdentityProvider
}
type userService struct {
	repo domain.UserRepository
	auth domain.IdentityProvider
	*base.BaseService
}

func NewService(base *base.BaseService, dep *ServiceDependencies) UserService {
	return &userService{
		repo:        dep.Repo,
		auth:        dep.Auth,
		BaseService: base,
	}
}

func (s *userService) CreateUser(ctx context.Context, input *command.CreateUserInput) (user *domain.User, err error) {
	// err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  resource.User,
	// 	Scope:     scope.Create,
	// })
	// if err != nil {
	// 	return
	// }
	user, err = domain.New(input.TenantDomain, input.BranchName, input.Email, input.FirstName, input.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	//user in authentication
	id, err := s.auth.CreateUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to save user auth: %w", err)
	}
	user.SetAuthID(id)
	//assign roles
	err = s.auth.AssignRolesToUser(ctx, &command.AssignRolesInput{
		BaseInput: input.BaseInput,
		UserID:    id,
		Roles:     input.Roles,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to assign user role auth: %w", err)
	}
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

func (s *userService) UpdateUser(ctx context.Context, input *command.UpdateUserInput) (user *domain.User, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.User,
		Scope:     scope.Update,
	})
	if err != nil {
		return
	}
	user, err = s.repo.FindByID(ctx, &command.UserIDInput{
		BaseInput: baseCmd.BaseInput{
			TenantDomain: input.TenantDomain,
			BranchName:   input.BranchName,
		},
		UserID: input.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	err = user.Update(input.Email, input.FirstName, input.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	err = s.auth.UpdateUser(ctx, &command.UpdateUserInput{
		BaseInput: input.BaseInput,
		ID:        user.AuthID(),
		Email:     user.Email(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Roles:     input.Roles,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update auth user: %w", err)
	}
	if err := s.repo.Update(ctx, input); err != nil {
		return nil, fmt.Errorf("failed to save updated user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, input *command.UserIDInput) (err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resource.User,
		Scope:     scope.Update,
	})
	if err != nil {
		return
	}
	userRepo, err := s.repo.FindByID(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to get auth user: %w", err)
	}
	err = s.auth.DeleteUser(ctx, &command.UserIDInput{
		BaseInput: input.BaseInput,
		UserID:    userRepo.AuthID(),
	})
	if err != nil {
		return fmt.Errorf("failed to delete auth user: %w", err)
	}
	err = s.repo.Remove(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
