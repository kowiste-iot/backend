package command

import "fmt"

type CreateTenantInput struct {
	Name        string `validate:"required,min=3,max=255"`
	Domain      string `validate:"required,min=3,max=255,alphanum"`
	Description string `validate:"omitempty,min=3,max=512"`
	AdminEmail  string `validate:"required,email"`
	Branch      string `validate:"required,min=3,max=255,alphanum"`
}

type UpdateTenantInput struct {
	ID          string `validate:"required,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Domain      string `validate:"required,min=3,max=255,alphanum"`
	Description string `validate:"omitempty,min=3,max=512"`
}

type CreateBranchInput struct {
	TenantDomain string `validate:"required,min=3,max=255,alphanum"`
	Name         string `validate:"required,min=3,max=255"`
	Description  string `validate:"omitempty,min=3,max=512"`
	Default      bool   
}

type UpdateBranchInput struct {
	ID           string `validate:"required,uuidv7"`
	TenantDomain string `validate:"required,min=3,max=255,alphanum"`
	Name         string `validate:"required,min=3,max=255"`
	Description  string `validate:"omitempty,min=3,max=512"`
}
type UserToBranch struct {
	TenantDomain string
	UserID       string
	Branchs      []string
}

func ClientName(branchName string) string {
	return fmt.Sprintf("%s-service", branchName)
}
