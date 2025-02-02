package keycloak

import (
	roleDomain "backend/internal/features/role/domain"
	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
	"context"
	"fmt"

	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

type Useukeycloak struct {
	*keycloak.Keycloak
}

func New(core *keycloak.Keycloak) *Useukeycloak {
	return &Useukeycloak{
		Keycloak: core,
	}
}

func (uk Useukeycloak) CreateUser(ctx context.Context, input *command.CreateUserInput) (string, error) {

	token, err := uk.GetValidToken(ctx)
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

	userID, err := uk.Client.CreateUser(ctx, token.AccessToken, input.TenantDomain, kcUser)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}
func (uk Useukeycloak) UpdateUser(ctx context.Context, input *command.UpdateUserInput) (err error) {
	return
}
func (uk Useukeycloak) DeleteUser(ctx context.Context, input *command.UserIDInput) (err error) {
	return

}
func (uk Useukeycloak) GetUser(ctx context.Context, input *command.UserIDInput) (user *domain.User, err error) {
	return

}
func (uk Useukeycloak) GetUserRoles(ctx context.Context, input *command.UserRolesInput) (roles []roleDomain.Role, err error) {
	return

}
func (uk Useukeycloak) AssignRolesToUser(ctx context.Context, input *command.AssignRolesInput) (err error) {
	return

}
func (uk Useukeycloak) RemoveRolesFromUser(ctx context.Context, input *command.RemoveRolesInput) (err error) {
	return

}
