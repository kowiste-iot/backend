package dto

import "time"

type UserDTO struct {
	ID        string    `json:"id"`
	AuthID    string    `json:"authId"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName" binding:"required"`
	LastName  string    `json:"lastName" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt"`
	Roles     []string  `json:"roles"`
}
