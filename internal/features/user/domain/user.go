package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id        string
	tenantID  string
	branch    string
	authID    string
	email     string
	firstName string
	lastName  string
	roles     []string
	updatedAt time.Time
	deletedAt *time.Time
}

func NewUser(tenantID, branch, email, firstName, lastName string) (user *User, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	user = &User{
		id:        id.String(),
		tenantID:  tenantID,
		branch:    branch,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		roles:     []string{},
		updatedAt: time.Now(),
	}
	return
}

func NewWithRoles(tenantID, branch, email, firstName, lastName string, roles []string) (user *User, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	user = &User{
		id:        id.String(),
		tenantID:  tenantID,
		branch:    branch,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		roles:     roles,
		updatedAt: time.Now(),
	}
	return
}

func NewFromRepository(id, tenantID, authID, email, firstName, lastName string, updatedAt time.Time, deletedAt *time.Time) *User {
	return &User{
		id:        id,
		tenantID:  tenantID,
		authID:    authID,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		roles:     []string{},
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func NewFromRepositoryWithRoles(id, tenantID, authID, email, firstName, lastName string, updatedAt time.Time, deletedAt *time.Time, roles []string) *User {
	return &User{
		id:        id,
		tenantID:  tenantID,
		authID:    authID,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		roles:     roles,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

func (u *User) Update(email, firstName, lastName string) error {
	u.email = email
	u.firstName = firstName
	u.lastName = lastName
	u.updatedAt = time.Now()
	return nil
}

func (u *User) SetAuthID(id string) {
	u.authID = id
}

func (u *User) Delete() {
	now := time.Now()
	u.deletedAt = &now
}

func (u *User) IsDeleted() bool {
	return u.deletedAt != nil
}

// Role Management
func (u *User) SetRoles(roles []string) {
	u.roles = roles
	u.updatedAt = time.Now()
}

func (u *User) AddRole(role string) {
	// Check if role already exists
	for _, r := range u.roles {
		if r == role {
			return
		}
	}
	u.roles = append(u.roles, role)
	u.updatedAt = time.Now()
}

func (u *User) RemoveRole(role string) {
	var newRoles []string
	for _, r := range u.roles {
		if r != role {
			newRoles = append(newRoles, r)
		}
	}
	u.roles = newRoles
	u.updatedAt = time.Now()
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

// Getters
func (u *User) ID() string            { return u.id }
func (u *User) TenantID() string      { return u.tenantID }
func (u *User) Branch() string        { return u.branch }
func (u *User) AuthID() string        { return u.authID }
func (u *User) Email() string         { return u.email }
func (u *User) FirstName() string     { return u.firstName }
func (u *User) LastName() string      { return u.lastName }
func (u *User) Roles() []string       { return u.roles }
func (u *User) UpdatedAt() time.Time  { return u.updatedAt }
func (u *User) DeletedAt() *time.Time { return u.deletedAt }

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidTenantID = errors.New("invalid tenant id")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidRole     = errors.New("invalid role")
)
