package command

type CheckPermissionInput struct {
	BaseInput
	Resource string
	Scope    string
}
