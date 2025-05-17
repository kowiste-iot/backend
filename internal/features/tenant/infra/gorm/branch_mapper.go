package repository

import (
	"backend/internal/features/tenant/domain"
	"time"
)

type BranchDB struct {
	ID           string     `gorm:"primaryKey;column:id"`
	TenantID     string     `gorm:"column:tenant_id;not null"`
	AuthBranchID string     `gorm:"column:auth_branch_id"`
	Name         string     `gorm:"column:name;not null"`
	Description  string     `gorm:"column:description"`
	Timezone     string     `gorm:"column:timezone;default:UTC"`
	SchemaName   string     `gorm:"column:schema_name;not null"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (BranchDB) TableName() string {
	return "platform.branches"
}

// Mapper domain model to db model
func branchMapper(schema string, tenant *domain.Tenant, data *domain.Branch) (dbBranch *BranchDB) {
	dbBranch = &BranchDB{
		ID:           data.ID(),
		TenantID:     tenant.ID(),
		AuthBranchID: data.AuthBranchID(),
		Name:         data.Name(),
		SchemaName:   schema,
		Description:  data.Description(),
		UpdatedAt:    data.UpdatedAt(),
	}
	return
}

type BranchInfo struct {
	TenantID    string    `gorm:"primaryKey;column:tenant_id"`
	BranchID    string    `gorm:"primaryKey;column:branch_id"`
	TenantName  string    `gorm:"column:tenant_name"`
	BranchName  string    `gorm:"column:branch_name"`
	Description string    `gorm:"column:description"`
	Timezone    string    `gorm:"column:timezone;default:'UTC'"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (b *BranchInfo) TableName() string {
	return "branch_info"
}

// Mapper domain model to db model
func branchInfoMapper(tenant *domain.Tenant, data *domain.Branch) (dbBranch *BranchInfo) {
	dbBranch = &BranchInfo{
		BranchID:    data.ID(),
		TenantID:    tenant.ID(),
		TenantName:  tenant.Domain(),
		BranchName:  data.Name(),
		Description: data.Description(),
		UpdatedAt:   data.UpdatedAt(),
	}
	return
}
