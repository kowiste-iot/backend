package domain

import (
	"ddd/shared/validator"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	id          string
	tenantID    string
	branchID    string
	name        string
	parent      *string
	description string
	updatedAt   time.Time
	deletedAt   *time.Time
}
type assetValidation struct {
	ID          string     `validate:"required,uuidv7"`
	TenantID    string     `validate:"required,uuidv7"`
	BranchID    string     `validate:"required,uuidv7"`
	Name        string     `validate:"required,min=3,max=255"`
	Parent      *string    `validate:"omitempty,uuidv7"`
	Description string     `validate:"omitempty,min=3,max=512"`
	UpdatedAt   time.Time  `validate:"required"`
	DeletedAt   *time.Time `validate:"omitempty"`
}

func New(tenantID, branchID, name, description string) (asset *Asset, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	asset = &Asset{
		id:          id.String(),
		branchID:    branchID,
		tenantID:    tenantID,
		name:        name,
		description: description,
		updatedAt:   time.Now(),
	}

	err = validator.Validate(assetValidation{
		ID:          asset.id,
		TenantID:    asset.tenantID,
		BranchID:    asset.branchID,
		Name:        asset.name,
		Parent:      asset.parent,
		Description: asset.description,
		UpdatedAt:   asset.updatedAt,
		DeletedAt:   asset.deletedAt,
	})

	return
}

func NewFromRepository(id, tenantID, branchID, name, description string, parent *string, updatedAt time.Time, deletedAt *time.Time) *Asset {
	return &Asset{
		id:          id,
		tenantID:    tenantID,
		branchID:    branchID,
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

	return validator.Validate(assetValidation{
		ID:          a.id,
		TenantID:    a.tenantID,
		BranchID:    a.branchID,
		Name:        a.name,
		Parent:      a.parent,
		Description: a.description,
		UpdatedAt:   a.updatedAt,
		DeletedAt:   a.deletedAt,
	})
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
func (a *Asset) BranchID() string      { return a.branchID }
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
