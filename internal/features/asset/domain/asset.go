package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	id          string
	tenantID    string
	branchName  string
	name        string
	parent      *string
	description string
	updatedAt   time.Time
	deletedAt   *time.Time
}

func New(tenantID, branchName, name, description string) (asset *Asset, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	asset = &Asset{
		id:          id.String(),
		branchName:  branchName,
		tenantID:    tenantID,
		name:        name,
		description: description,
		updatedAt:   time.Now(),
	}

	return
}

func NewFromRepository(id, tenantID, branchName, name, description string, parent *string, updatedAt time.Time, deletedAt *time.Time) *Asset {
	return &Asset{
		id:          id,
		tenantID:    tenantID,
		branchName:  branchName,
		name:        name,
		parent:      parent,
		description: description,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}
}

func (a *Asset) Update(name, parent, description string) (err error) {
	a.name = name
	a.description = description
	a.updatedAt = time.Now()
	if parent != "" {
		a.parent = &parent
	}

	return
}
func (a *Asset) Delete() {
	now := time.Now()
	a.deletedAt = &now
}
func (a *Asset) WithParent(parent string) {
	a.parent = &parent
}
func (a *Asset) IsDeleted() bool {
	return a.deletedAt != nil
}

// Getters
func (a *Asset) ID() string            { return a.id }
func (a *Asset) TenantID() string      { return a.tenantID }
func (a *Asset) BranchName() string    { return a.branchName }
func (a *Asset) Name() string          { return a.name }
func (a *Asset) Parent() *string       { return a.parent }
func (a *Asset) Description() string   { return a.description }
func (a *Asset) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Asset) DeletedAt() *time.Time { return a.deletedAt }

var (
	ErrAssetNotFound   = errors.New("asset not found")
	ErrInvalidTenantID = errors.New("invalid tenant id")
	ErrInvalidName     = errors.New("invalid name")
)
