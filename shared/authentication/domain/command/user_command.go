package command

import "ddd/shared/base/command"

type CreateUserInput struct {
	command.BaseInput
	ID        string
	Email     string
	FirstName string
	LastName  string
	Roles     []string
}

type UpdateUserInput struct {
	command.BaseInput
	ID        string
	Email     string
	FirstName string
	LastName  string
	Roles     []string
}
type UserIDInput struct {
	command.BaseInput
	UserID string
}
type AssignRolesInput struct {
	command.BaseInput
	UserID string    
	Roles  []string  
}

type RemoveRolesInput struct {
	command.BaseInput
	UserID string    
	Roles  []string 
}

type UserRolesInput struct {
	command.BaseInput
	UserID string 
}