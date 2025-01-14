package command

import "ddd/shared/base/command"

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

type UserToBranch struct{
	command.BaseInput
	UserID string
}