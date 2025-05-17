package gormhelper

import (
	"backend/shared/http/httputil"
	"context"
	"fmt"

	"gorm.io/gorm"
)

func TenantFilter(tenantID string) string {
	return fmt.Sprintf("tenant_id = '%s' AND deleted_at IS NULL", tenantID)
}

func DeleteFilter() string {
	return "deleted_at IS NULL"
}

func GetBranchName(tenantID, branchID string) string {
	return fmt.Sprintf("branch_%s_%s", tenantID, branchID)
}
func SetBranchDB(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	t, b, err := httputil.GetBase(ctx)
	if err != nil {
		return nil, err
	}
	err = db.Exec(fmt.Sprintf("SET search_path TO %s", GetBranchName(t.Domain(), b))).Error
	if err != nil {
		return nil, err
	}
	return db, nil
}
