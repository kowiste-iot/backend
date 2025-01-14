package command

type CreateTenantInput struct {
	Domain      string
	Name        string
	Description string
	AdminEmail string
	DefaultBranch string
}
