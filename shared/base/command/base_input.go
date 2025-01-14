package command

type BaseInput struct {
	TenantDomain string `validate:"required,uuidv7"`
	BranchName   string `validate:"required,uuidv7"`
}

func NewInput(tenant, branch string) BaseInput {
	return BaseInput{
		TenantDomain: tenant,
		BranchName:   branch,
	}
}
