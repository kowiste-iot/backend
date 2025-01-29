package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Alert struct {
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

func New(tenantID, branchName, name, parent, description string, enable bool) (alert *Alert, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	alert = &Alert{
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

func NewFromRepository(id, tenantID, branchName, name, description string, parent string, enabled bool, updatedAt time.Time, deletedAt *time.Time) *Alert {
	return &Alert{
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

func (a *Alert) Update(name, parent, description string, enabled bool) (err error) {
	a.name = name
	a.description = description
	a.updatedAt = time.Now()
	a.parent = parent
	a.enabled = enabled

	return
}
func (a *Alert) Delete() {
	now := time.Now()
	a.deletedAt = &now
}

func (a *Alert) IsDeleted() bool {
	return a.deletedAt != nil
}

// Getters
func (a *Alert) ID() string            { return a.id }
func (a *Alert) TenantID() string      { return a.tenantID }
func (a *Alert) BranchName() string    { return a.branchName }
func (a *Alert) Name() string          { return a.name }
func (a *Alert) Parent() string        { return a.parent }
func (a *Alert) Enabled() bool         { return a.enabled }
func (a *Alert) Description() string   { return a.description }
func (a *Alert) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Alert) DeletedAt() *time.Time { return a.deletedAt }

var (
	ErrAlertNotFound  = errors.New("alert not found")
	ErrInvalidTenantID = errors.New("invalid tenant id")
	ErrInvalidName     = errors.New("invalid name")
)
