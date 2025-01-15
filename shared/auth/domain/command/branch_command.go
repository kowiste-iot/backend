package command

import (
	"ddd/shared/base/command"
	"fmt"
)

type CreateBranchInput struct {
	TenantID    string
	Name        string
	Description string
}

type UpdateBranchInput struct {
	TenantID    string
	ID          string
	Name        string
	Description string
}

type UserToBranch struct {
	command.BaseInput
	UserID string
}

func ClientName(branchName string) string {
	return fmt.Sprintf("%s-service", branchName)
}
