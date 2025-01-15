package keycloak

import (
	"context"
	userCmd "ddd/internal/features/user/domain/command"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/command"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

func (ks *KeycloakService) CreateUser(ctx context.Context, input *userCmd.CreateUserInput) (string, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	enabled := true
	kcUser := gocloak.User{
		Username:      &input.Email,
		Email:         &input.Email,
		Enabled:       &enabled,
		FirstName:     &input.FirstName,
		LastName:      &input.LastName,
		EmailVerified: &enabled,
	}

	userID, err := ks.client.CreateUser(ctx, token.AccessToken, input.TenantDomain, kcUser)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	err=ks.AssignRoles(ctx, &command.AssignRolesInput{
		BaseInput: input.BaseInput,
		UserID:    userID,
		Roles:     input.Roles,
	})
		if err != nil {
		return "", fmt.Errorf("failed to set user role: %w", err)
	}
	return userID, nil
}

func (ks *KeycloakService) UpdateUser(ctx context.Context, input *userCmd.UpdateUserInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	kcUser := gocloak.User{
		Username:  &input.Email,
		Email:     &input.Email,
		FirstName: &input.FirstName,
		LastName:  &input.LastName,
	}

	err = ks.client.UpdateUser(ctx, token.AccessToken, input.TenantDomain, kcUser)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (ks *KeycloakService) DeleteUser(ctx context.Context, input *userCmd.UserIDInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeleteUser(ctx, token.AccessToken, input.TenantDomain, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (ks *KeycloakService) GetUser(ctx context.Context, input *userCmd.UserIDInput) (*auth.User, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	user, err := ks.client.GetUserByID(ctx, token.AccessToken, input.TenantDomain, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &auth.User{
		ID:        *user.ID,
		TenantID:  input.TenantDomain,
		Email:     *user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
	}, nil
}
