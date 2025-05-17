package domain

import (
	userDomain "backend/internal/features/user/domain"
	"backend/shared/authentication/domain/command"
	baseCmd "backend/shared/base/command"
	"context"
)

type IBranch interface {
	CreateBranch(ctx context.Context, input *command.CreateBranchInput) (string, error)
	UpdateBranch(ctx context.Context, input *command.UpdateBranchInput) error
	DeleteBranch(ctx context.Context, input *baseCmd.BaseInput) error
	GetBranch(ctx context.Context, input *baseCmd.BaseInput) (*Branch, error)
	GetBranchUsers(ctx context.Context, input *baseCmd.BaseInput) ([]userDomain.User, error)
	AssignUserToBranch(ctx context.Context, input *command.UserToBranch) error
	RemoveUserFromBranch(ctx context.Context, input *command.UserToBranch) error
}

type Branch struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

const (
	AdminBranch     string = "admin"
	DefaultBranch   string = "default"
	UndefinedBranch string = "undefined"
)

func ForbiddenBranch() []string {
	return []string{
		AdminBranch,
		DefaultBranch,
		UndefinedBranch,
	}

}
