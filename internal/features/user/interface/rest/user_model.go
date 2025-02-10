package userhandler

type CreateUserRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	FirstName string   `json:"firstName" binding:"required"`
	LastName  string   `json:"lastName" binding:"required"`
	Roles     []string `json:"roles" binding:"required"`
}

type UpdateUserRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	FirstName string   `json:"firstName" binding:"required"`
	LastName  string   `json:"lastName" binding:"required"`
	Roles     []string `json:"roles" binding:"required"`
}
