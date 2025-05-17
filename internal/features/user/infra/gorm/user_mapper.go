package repository

import (
	"backend/internal/features/user/domain"
	"time"

	"gorm.io/gorm"
)

 
type UserDB struct {
	ID         string `gorm:"primaryKey"`
	AuthUserID string `gorm:"uniqueIndex"`
	Email      string
	FirstName  string
	LastName   string
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type RoleDB struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type UserRoleDB struct {
	UserID    string `gorm:"primaryKey"`
	RoleID    string `gorm:"primaryKey"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName sets the table name for UserDB
func (UserDB) TableName() string {
	return "users"
}

// TableName sets the table name for RoleDB
func (RoleDB) TableName() string {
	return "roles"
}

// TableName sets the table name for UserRoleDB
func (UserRoleDB) TableName() string {
	return "user_roles"
}

// Mapper domain model to db model
func userMapper(data *domain.User) (dbUser *UserDB) {
	dbUser = &UserDB{
		ID:         data.ID(),
		AuthUserID: data.AuthID(),
		Email:      data.Email(),
		FirstName:  data.FirstName(),
		LastName:   data.LastName(),
		UpdatedAt:  data.UpdatedAt(),
	}
	return
}
