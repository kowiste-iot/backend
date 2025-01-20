package command

type BaseInput struct {
	TenantDomain string `validate:"required"`
	BranchName   string `validate:"required"`
	ClientID     *string
}

func NewInput(tenant, branch string) BaseInput {
	return BaseInput{
		TenantDomain: tenant,
		BranchName:   branch,
	}
}

func (b *BaseInput) WithClientID(id string) BaseInput {
	b.ClientID = &id
	return *b
}
