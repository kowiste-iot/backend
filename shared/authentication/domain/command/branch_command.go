package command

import (
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
	TenantDomain string
	UserID   string
	Branchs   []string
}

func ClientName(branchName string) string {
	return fmt.Sprintf("%s-service", branchName)
}
