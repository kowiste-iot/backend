package keycloak

import (
	"context"
	auth "ddd/shared/auth/domain"
	"ddd/shared/auth/domain/command"
	baseCmd "ddd/shared/base/command"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

func (ks *KeycloakService) CreateBranch(ctx context.Context, input *command.CreateBranchInput) (string, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	group := mapBranchToGroup(input)
	groupID, err := ks.client.CreateGroup(
		ctx,
		token.AccessToken,
		input.TenantID,
		*group,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create branch: %w", err)
	}

	return groupID, nil
}

func (ks *KeycloakService) UpdateBranch(ctx context.Context, input *command.UpdateBranchInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	group := mapBranchToGroup(&command.CreateBranchInput{
		TenantID:    input.TenantID,
		Name:        input.Name,
		Description: input.Description,
	})
	group.ID = &input.ID

	err = ks.client.UpdateGroup(
		ctx,
		token.AccessToken,
		input.TenantID,
		*group,
	)
	if err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}

	return nil
}

func (ks *KeycloakService) DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = ks.client.DeleteGroup(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.BranchName,
	)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

func (ks *KeycloakService) GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*auth.Branch, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	group, err := ks.client.GetGroup(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.BranchName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	return mapGroupToBranch(group), nil
}

func (ks *KeycloakService) GetBranchUsers(ctx context.Context, input *baseCmd.BaseInput) ([]auth.User, error) {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	users, err := ks.client.GetGroupMembers(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.BranchName,
		gocloak.GetGroupsParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch users: %w", err)
	}

	var authUsers []auth.User
	for _, user := range users {
		authUsers = append(authUsers, auth.User{
			ID:        *user.ID,
			FirstName: *user.FirstName,
			LastName:  *user.LastName,
			Email:     *user.Email,
		})
	}

	return authUsers, nil
}

func (ks *KeycloakService) AssignUserToBranch(ctx context.Context, input *command.UserToBranch) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.AddUserToGroup(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.UserID,
		input.BranchName,
	)
	if err != nil {
		return fmt.Errorf("failed to assign user to branch: %w", err)
	}

	return nil
}

func (ks *KeycloakService) RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error {
	token, err := ks.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = ks.client.DeleteUserFromGroup(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.UserID,
		input.BranchName,
	)
	if err != nil {
		return fmt.Errorf("failed to remove user from branch: %w", err)
	}

	return nil
}

// Helper functions for mapping between domain and Keycloak types
func mapBranchToGroup(branch *command.CreateBranchInput) *gocloak.Group {
	return &gocloak.Group{
		Name: &branch.Name,
		Attributes: &map[string][]string{
			"description": {branch.Description},
		},
	}
}

func mapGroupToBranch(group *gocloak.Group) *auth.Branch {
	description := ""
	if group.Attributes != nil {
		if desc, ok := (*group.Attributes)["description"]; ok && len(desc) > 0 {
			description = desc[0]
		}
	}

	return &auth.Branch{
		ID:          *group.ID,
		Name:        *group.Name,
		Description: description,
		Attributes:  make(map[string]string),
	}
}
