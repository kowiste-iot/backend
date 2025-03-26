package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Dashboard struct {
	id          string
	tenantID    string
	branchName  string
	name        string
	parent      string
	description string
	updatedAt   time.Time
	deletedAt   *time.Time
}

func New(tenantID, branchName, name, parent, description string) (dashboard *Dashboard, err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	dashboard = &Dashboard{
		id:          id.String(),
		branchName:  branchName,
		tenantID:    tenantID,
		parent:      parent,
		name:        name,
		description: description,
		updatedAt:   time.Now(),
	}

	return
}

func NewFromRepository(id, tenantID, branchName, name, description string, parent string, updatedAt time.Time, deletedAt *time.Time) *Dashboard {
	return &Dashboard{
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

func (a *Dashboard) Update(name, parent, description string) (err error) {
	a.name = name
	a.description = description
	a.updatedAt = time.Now()
	a.parent = parent

	return
}
func (a *Dashboard) Delete() {
	now := time.Now()
	a.deletedAt = &now
}

func (a *Dashboard) IsDeleted() bool {
	return a.deletedAt != nil
}

// Getters
func (a *Dashboard) ID() string            { return a.id }
func (a *Dashboard) TenantID() string      { return a.tenantID }
func (a *Dashboard) BranchName() string    { return a.branchName }
func (a *Dashboard) Name() string          { return a.name }
func (a *Dashboard) Parent() string        { return a.parent }
func (a *Dashboard) Description() string   { return a.description }
func (a *Dashboard) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Dashboard) DeletedAt() *time.Time { return a.deletedAt }

var (
	ErrDashboardNotFound = errors.New("parent not found")
	ErrInvalidTenantID   = errors.New("invalid tenant id")
	ErrInvalidName       = errors.New("invalid name")
)
