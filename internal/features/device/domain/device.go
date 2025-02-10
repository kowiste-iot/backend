package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	id          string
	tenantID    string
	branchName  string
	name        string
	parent      string
	description string
	updatedAt   time.Time
	deletedAt   *time.Time
}

func New(tenantID, branchName, name, parent, description string) (device *Device, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	device = &Device{
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

func NewFromRepository(id, tenantID, branchName, name, description string, parent string, updatedAt time.Time, deletedAt *time.Time) *Device {
	return &Device{
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

func (a *Device) Update(name, parent, description string) (err error) {
	a.name = name
	a.description = description
	a.updatedAt = time.Now()
	a.parent = parent

	return
}
func (a *Device) Delete() {
	now := time.Now()
	a.deletedAt = &now
}

func (a *Device) IsDeleted() bool {
	return a.deletedAt != nil
}

// Getters
func (a *Device) ID() string            { return a.id }
func (a *Device) TenantID() string      { return a.tenantID }
func (a *Device) BranchName() string    { return a.branchName }
func (a *Device) Name() string          { return a.name }
func (a *Device) Parent() string        { return a.parent }
func (a *Device) Description() string   { return a.description }
func (a *Device) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Device) DeletedAt() *time.Time { return a.deletedAt }

var (
	ErrDeviceNotFound = errors.New("parent not found")
	ErrInvalidTenantID   = errors.New("invalid tenant id")
	ErrInvalidName       = errors.New("invalid name")
)
