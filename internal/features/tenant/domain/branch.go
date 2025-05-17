package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	AdminBranch     string = "admin"
	DefaultBranch   string = "default"
	UndefinedBranch string = "undefined"
)

func ForbiddenBranch() []string {
	return []string{
		AdminBranch,
		DefaultBranch,
		UndefinedBranch,
	}

}

type Branch struct {
	id            string
	tenantID      string
	authBranchID  string
	adminBranchID string
	name          string
	description   string
	updatedAt     time.Time
	deletedAt     *time.Time
}


func NewBranch(tenantID, name, description string) (branch *Branch, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	branch = &Branch{
		id:          id.String(),
		tenantID:    tenantID,
		name:        name,
		description: description,
		updatedAt:   time.Now(),
	}
	return
}

func NewBranchFromRepository(id, tenantID, authBranchID, adminBranchID, name, description string, updatedAt time.Time, deletedAt *time.Time) *Branch {
	return &Branch{
		id:            id,
		tenantID:      tenantID,
		authBranchID:  authBranchID,
		adminBranchID: adminBranchID,
		name:          name,
		description:   description,
		updatedAt:     updatedAt,
		deletedAt:     deletedAt,
	}
}

func (b *Branch) Update(name, description string) error {
	if name == "" {
		return ErrInvalidName
	}

	b.name = name
	b.description = description
	b.updatedAt = time.Now()
	return nil
}

func (b *Branch) SetAuthBranchID(id string) {
	b.authBranchID = id
}
func (b *Branch) SetAdminBranchID(id string) {
	b.adminBranchID = id
}

func (b *Branch) Delete() {
	now := time.Now()
	b.deletedAt = &now
}

func (b *Branch) IsDeleted() bool {
	return b.deletedAt != nil
}

// Getters
func (b *Branch) ID() string            { return b.id }
func (b *Branch) TenantID() string      { return b.tenantID }
func (b *Branch) AuthBranchID() string  { return b.authBranchID }
func (b *Branch) AdminBranchID() string  { return b.adminBranchID }
func (b *Branch) Name() string          { return b.name }
func (b *Branch) Description() string   { return b.description }
func (b *Branch) UpdatedAt() time.Time  { return b.updatedAt }
func (b *Branch) DeletedAt() *time.Time { return b.deletedAt }

// Domain errors (can be added to existing errors.go)
var (
	ErrBranchNotFound  = errors.New("branch not found")
	ErrInvalidBranchID = errors.New("invalid branch id")
)
