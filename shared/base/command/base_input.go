package command

//TODO: add client.ID to this base is neccesary in a lot of places
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
