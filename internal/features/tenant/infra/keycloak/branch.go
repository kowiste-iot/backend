package keycloak

import (
	"backend/internal/features/tenant/domain"
	"backend/internal/features/tenant/domain/command"
	userDomain "backend/internal/features/user/domain"
	baseCmd "backend/shared/base/command"
	"backend/shared/keycloak"
	"context"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v13"
)

type BranchKeycloak struct {
	*keycloak.Keycloak
}

func NewBranch(core *keycloak.Keycloak) *BranchKeycloak {
	return &BranchKeycloak{
		Keycloak: core,
	}
}

func (rk BranchKeycloak) CreateBranch(ctx context.Context, input *command.CreateBranchInput) (string, error) {

	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	group := mapBranchToGroup(input)
	groupID, err := rk.Client.CreateGroup(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*group,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create branch: %w", err)
	}

	return groupID, nil
}
func (rk BranchKeycloak) DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	err = rk.Client.DeleteGroup(
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
func (rk BranchKeycloak) GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*domain.Branch, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	group, err := rk.Client.GetGroup(
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
func (rk BranchKeycloak) UpdateBranch(ctx context.Context, input *command.UpdateBranchInput) (err error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	group := mapBranchToGroup(&command.CreateBranchInput{
		TenantDomain: input.TenantDomain,
		Name:         input.Name,
		Description:  input.Description,
	})
	group.ID = &input.ID

	err = rk.Client.UpdateGroup(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		*group,
	)
	if err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}
	return 
}
func (rk BranchKeycloak) GetBranchUsers(ctx context.Context, input *baseCmd.BaseInput) ([]userDomain.User, error) {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	users, err := rk.Client.GetGroupMembers(
		ctx,
		token.AccessToken,
		input.TenantDomain,
		input.BranchName,
		gocloak.GetGroupsParams{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch users: %w", err)
	}

	var authUsers []userDomain.User
	for _, user := range users {
		authUsers = append(authUsers, *userDomain.NewFromRepository(*user.ID, input.TenantDomain, *user.ID,
			*user.Email, *user.FirstName, *user.LastName, time.Now(), nil))
	}

	return authUsers, nil
}
func (rk BranchKeycloak) AssignUserToBranch(ctx context.Context, input *command.UserToBranch) error {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	for i := range input.Branchs {
		err = rk.Client.AddUserToGroup(
			ctx,
			token.AccessToken,
			input.TenantDomain,
			input.UserID,
			input.Branchs[i],
		)
		if err != nil {
			return fmt.Errorf("failed to assign user to branch: %w", err)
		}
	}

	return nil
}
func (rk BranchKeycloak) RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error {
	token, err := rk.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	for i := range input.Branchs {
		err = rk.Client.DeleteUserFromGroup(
			ctx,
			token.AccessToken,
			input.TenantDomain,
			input.UserID,
			input.Branchs[i],
		)
		if err != nil {
			return fmt.Errorf("failed to remove user from branch: %w", err)
		}
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

func mapGroupToBranch(group *gocloak.Group) *domain.Branch {
	description := ""
	if group.Attributes != nil {
		if desc, ok := (*group.Attributes)["description"]; ok && len(desc) > 0 {
			description = desc[0]
		}
	}
	return domain.NewBranchFromRepository(*group.ID, "", *group.ID, *group.Name, description, time.Now(), nil)

}
