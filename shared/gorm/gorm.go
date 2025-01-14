package gormhelper

import "fmt"

func TenantFilter(tenantID string) string {
	return fmt.Sprintf("tenant_id = '%s' AND deleted_at IS NULL", tenantID)
}
func TenantBranchFilter(tenantID, branchID string) string {
	return fmt.Sprintf("tenant_id = '%s' AND branch_id = '%s' AND deleted_at IS NULL", tenantID, branchID)
}
