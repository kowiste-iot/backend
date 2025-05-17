// app/user_service.go
package app

import (
	"context"

	resourceDomain "backend/internal/features/resource/domain"
	scopeDomain "backend/internal/features/scope/domain"
	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
	"backend/internal/features/user/dto"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"

	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, input *command.CreateUserInput) (*dto.UserDTO, error)
	GetUser(ctx context.Context, input *command.UserIDInput) (*dto.UserDTO, error)
	ListUsers(ctx context.Context, input *baseCmd.BaseInput) ([]*dto.UserDTO, error)
	UpdateUser(ctx context.Context, input *command.UpdateUserInput) (*dto.UserDTO, error)
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

func NewUserService(base *base.BaseService, dep *ServiceDependencies) UserService {
	return &userService{
		repo:        dep.Repo,
		auth:        dep.Auth,
		BaseService: base,
	}
}

func (s *userService) CreateUser(ctx context.Context, input *command.CreateUserInput) (user *dto.UserDTO, err error) {
	// err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
	// 	BaseInput: input.BaseInput,
	// 	Resource:  resource.User,
	// 	Scope:     scope.Create,
	// })
	// if err != nil {
	// 	return
	// }
	u, err := domain.NewUser(input.TenantDomain, input.BranchName, input.Email, input.FirstName, input.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	//user in authentication
	id, err := s.auth.CreateUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to save user auth: %w", err)
	}
	u.SetAuthID(id)
	//assign roles
	err = s.auth.AssignRolesToUser(ctx, &command.AssignRolesInput{
		BaseInput: input.BaseInput,
		UserID:    id,
		Roles:     input.Roles,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to assign user role auth: %w", err)
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	return dto.ToDTO(u, input.Roles), nil
}

func (s *userService) GetUser(ctx context.Context, input *command.UserIDInput) (*dto.UserDTO, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.User,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindByID(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	outRoles, err := s.getRolesString(ctx, &command.UserIDInput{
		BaseInput: input.BaseInput,
		UserID:    user.AuthID(),
	})
	if err != nil {
		return nil, err
	}
	return dto.ToDTO(user, outRoles), nil
}

func (s *userService) ListUsers(ctx context.Context, input *baseCmd.BaseInput) (usersDTO []*dto.UserDTO, err error) {
	users, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	for i := range users {
		outRoles, err := s.getRolesString(ctx, &command.UserIDInput{
			BaseInput: *input,
			UserID:    users[i].AuthID(),
		})
		if err != nil {
			return nil, err
		}
		usersDTO = append(usersDTO, dto.ToDTO(users[i], outRoles))
	}
	return
}

func (s *userService) UpdateUser(ctx context.Context, input *command.UpdateUserInput) (user *dto.UserDTO, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.User,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return
	}
	u, err := s.repo.FindByID(ctx, &command.UserIDInput{
		BaseInput: baseCmd.BaseInput{
			TenantDomain: input.TenantDomain,
			BranchName:   input.BranchName,
		},
		UserID: input.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	err = u.Update(input.Email, input.FirstName, input.LastName)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	//TODO: keycloak return 405 method not allowed...
	// err = s.auth.UpdateUser(ctx, &command.UpdateUserInput{
	// 	BaseInput: input.BaseInput,
	// 	ID:        u.AuthID(),
	// 	Email:     u.Email(),
	// 	FirstName: u.FirstName(),
	// 	LastName:  u.LastName(),
	// 	Roles:     input.Roles,
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to update auth user: %w", err)
	// }
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to save updated user: %w", err)
	}

	return dto.ToDTO(u, input.Roles), nil
}

func (s *userService) DeleteUser(ctx context.Context, input *command.UserIDInput) (err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.User,
		Scope:     scopeDomain.Update,
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

func (s userService) getRolesString(ctx context.Context, input *command.UserIDInput) (roles []string, err error) {
	r, err := s.auth.GetUserRoles(ctx, &command.UserRolesInput{
		BaseInput: *&input.BaseInput,
		UserID:    input.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	for j := range roles {
		roles = append(roles, r[j].Name)
	}
	return
}
