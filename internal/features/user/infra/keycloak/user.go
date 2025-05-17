package keycloak

import (
	tenantCmd "backend/internal/features/tenant/domain/command"
	"backend/internal/features/user/domain"
	roleDomain "backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"

	"context"
	"fmt"

	"backend/shared/keycloak"

	"github.com/Nerzal/gocloak/v13"
)

type Userkeycloak struct {
	*keycloak.Keycloak
}

func New(core *keycloak.Keycloak) *Userkeycloak {
	return &Userkeycloak{
		Keycloak: core,
	}
}

func (uk Userkeycloak) CreateUser(ctx context.Context, input *command.CreateUserInput) (string, error) {

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
func (uk Userkeycloak) UpdateUser(ctx context.Context, input *command.UpdateUserInput) (err error) {
	token, err := uk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	kcUser := gocloak.User{
		Username:  &input.Email,
		Email:     &input.Email,
		FirstName: &input.FirstName,
		LastName:  &input.LastName,
	}

	err = uk.Client.UpdateUser(ctx, token.AccessToken, input.TenantDomain, kcUser)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return
}
func (uk Userkeycloak) DeleteUser(ctx context.Context, input *command.UserIDInput) (err error) {
	token, err := uk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = uk.Client.DeleteUser(ctx, token.AccessToken, input.TenantDomain, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return
}
func (uk Userkeycloak) GetUser(ctx context.Context, input *command.UserIDInput) (user *domain.User, err error) {
	token, err := uk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	userKC, err := uk.Client.GetUserByID(ctx, token.AccessToken, input.TenantDomain, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	user, err = domain.NewUser(input.TenantDomain, input.BranchName, *userKC.Email, *userKC.FirstName, *userKC.LastName)
	user.SetAuthID(*userKC.ID)
	return

}
func (uk Userkeycloak) GetUserRoles(ctx context.Context, input *command.UserRolesInput) (roles []roleDomain.Role, err error) {
	token, err := uk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	err = uk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	clientRoles, err := uk.Client.GetClientRolesByUserID(ctx, token.AccessToken,
		input.TenantDomain, *input.ClientID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting roles: %v", err)
	}
	roles = make([]roleDomain.Role, 0)
	for i := range clientRoles {
		roles = append(roles, roleDomain.Role{
			Name:        *clientRoles[i].Name,
			Description: *clientRoles[i].Description,
		})
	}
	return
}
func (uk Userkeycloak) AssignRolesToUser(ctx context.Context, input *command.AssignRolesInput) (err error) {
	token, err := uk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = uk.FetchClient(ctx, &input.BaseInput)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	addRoles := make([]gocloak.Role, 0)
	for i := range input.Roles {
		role, err := uk.Client.GetClientRole(ctx, token.AccessToken, input.TenantDomain, *input.ClientID, input.Roles[i])
		if err != nil {
			return fmt.Errorf("error getting realm role: %v", err)
		}
		addRoles = append(addRoles, *role)
	}
	err = uk.Client.AddClientRolesToUser(ctx, token.AccessToken,
		input.TenantDomain, *input.ClientID, input.UserID, addRoles)

	return
}
func (uk Userkeycloak) RemoveRolesFromUser(ctx context.Context, input *command.RemoveRolesInput) (err error) {
	token, err := uk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	removeRoles := make([]gocloak.Role, 0)
	for i := range input.Roles {
		role, err := uk.Client.GetRealmRole(ctx, token.AccessToken, input.TenantDomain, input.Roles[i])
		if err != nil {
			return fmt.Errorf("error getting realm role: %v", err)
		}
		removeRoles = append(removeRoles, *role)
	}
	err = uk.Client.DeleteClientRolesFromUser(ctx, token.AccessToken,
		input.TenantDomain, tenantCmd.ClientName(input.BranchName), input.UserID, removeRoles)
	return
}
