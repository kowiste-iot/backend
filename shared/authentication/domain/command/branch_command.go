package command

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
	UserID       string
	Branchs      []string
}
