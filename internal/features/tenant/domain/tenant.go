package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	id          string
	authID      string
	name        string
	domain      string
	description string
	updatedAt   time.Time
	deletedAt   *time.Time
}


func NewTenant(name, domain, description string) (tenant *Tenant, err error) {
	
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	tenant = &Tenant{
		id:          id.String(),
		name:        name,
		domain:      domain,
		description: description,
		updatedAt:   time.Now(),
	}
	return
}

func NewFromRepository(id, authID, name, domain, description string, updatedAt time.Time, deletedAt *time.Time) *Tenant {
	return &Tenant{
		id:          id,
		authID:      authID,
		name:        name,
		domain:      domain,
		description: description,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}
}

func (a *Tenant) Update(name, domain, description string) (err error) {

	a.name = name
	a.domain = domain
	a.description = description
	a.updatedAt = time.Now()

	return
}
func (a *Tenant) SetAuthID(id string) {
	a.authID = id
}
func (a *Tenant) Delete() {
	now := time.Now()
	a.deletedAt = &now
}

func (a *Tenant) IsDeleted() bool {
	return a.deletedAt != nil
}

// Getters
func (a *Tenant) ID() string            { return a.id }
func (a *Tenant) AuhtID() string        { return a.authID }
func (a *Tenant) Name() string          { return a.name }
func (a *Tenant) Domain() string        { return a.domain }
func (a *Tenant) Description() string   { return a.description }
func (a *Tenant) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Tenant) DeletedAt() *time.Time { return a.deletedAt }

var (
	ErrAssetNotFound   = errors.New("tenant not found")
	ErrInvalidTenantID = errors.New("invalid tenant id")
	ErrInvalidName     = errors.New("invalid name")
)
