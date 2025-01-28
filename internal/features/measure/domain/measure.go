package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Measure struct {
	id          string
	tenantID    string
	branchName  string
	name        string
	parent      string
	description string
	updatedAt   time.Time
	deletedAt   *time.Time
}

func New(tenantID, branchName, name, parent, description string) (measure *Measure, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	measure = &Measure{
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

func NewFromRepository(id, tenantID, branchName, name, description string, parent string, updatedAt time.Time, deletedAt *time.Time) *Measure {
	return &Measure{
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

func (a *Measure) Update(name, parent, description string) (err error) {
	a.name = name
	a.description = description
	a.updatedAt = time.Now()
	a.parent = parent

	return
}
func (a *Measure) Delete() {
	now := time.Now()
	a.deletedAt = &now
}

func (a *Measure) IsDeleted() bool {
	return a.deletedAt != nil
}

// Getters
func (a *Measure) ID() string            { return a.id }
func (a *Measure) TenantID() string      { return a.tenantID }
func (a *Measure) BranchName() string    { return a.branchName }
func (a *Measure) Name() string          { return a.name }
func (a *Measure) Parent() string        { return a.parent }
func (a *Measure) Description() string   { return a.description }
func (a *Measure) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Measure) DeletedAt() *time.Time { return a.deletedAt }

var (
	ErrMeasureNotFound = errors.New("parent not found")
	ErrInvalidTenantID   = errors.New("invalid tenant id")
	ErrInvalidName       = errors.New("invalid name")
)
