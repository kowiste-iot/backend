package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Action struct {
	id          string
	tenantID    string
	branchName  string
	name        string
	parent      string
	enabled     bool
	description string
	updatedAt   time.Time
	deletedAt   *time.Time
}

func New(tenantID, branchName, name, parent, description string, enable bool) (action *Action, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	action = &Action{
		id:          id.String(),
		branchName:  branchName,
		tenantID:    tenantID,
		parent:      parent,
		enabled:     enable,
		name:        name,
		description: description,
		updatedAt:   time.Now(),
	}

	return
}

func NewFromRepository(id, tenantID, branchName, name, description string, parent string, enabled bool, updatedAt time.Time, deletedAt *time.Time) *Action {
	return &Action{
		id:          id,
		tenantID:    tenantID,
		branchName:  branchName,
		name:        name,
		parent:      parent,
		enabled:     enabled,
		description: description,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}
}

func (a *Action) Update(name, parent, description string, enabled bool) (err error) {
	a.name = name
	a.description = description
	a.updatedAt = time.Now()
	a.parent = parent
	a.enabled = enabled

	return
}
func (a *Action) Delete() {
	now := time.Now()
	a.deletedAt = &now
}

func (a *Action) IsDeleted() bool {
	return a.deletedAt != nil
}

// Getters
func (a *Action) ID() string            { return a.id }
func (a *Action) TenantID() string      { return a.tenantID }
func (a *Action) BranchName() string    { return a.branchName }
func (a *Action) Name() string          { return a.name }
func (a *Action) Parent() string        { return a.parent }
func (a *Action) Enabled() bool         { return a.enabled }
func (a *Action) Description() string   { return a.description }
func (a *Action) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Action) DeletedAt() *time.Time { return a.deletedAt }

var (
	ErrActionNotFound  = errors.New("action not found")
	ErrInvalidTenantID = errors.New("invalid tenant id")
	ErrInvalidName     = errors.New("invalid name")
)
