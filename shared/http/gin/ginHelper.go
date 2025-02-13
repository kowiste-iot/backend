package ginhelp

import (
	"backend/shared/pagination"
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetPaginationGin(ctx context.Context, c *gin.Context) context.Context {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	sortField := c.DefaultQuery("sort_by", "")
	sortDir := c.DefaultQuery("sort_dir", "")

	return pagination.SetPagination(ctx, page, size, sortField, sortDir)
}

// func GetTenantID(c *gin.Context) (context.Context, bool) {
// 	tenantID, ok := c.Params.Get("tenantid")
// 	if !ok {
// 		return c.Request.Context(), false
// 	}
// 	return httputil.SetTenant(c.Request.Context(), tenantID), true
// }

func GetAuhtHeader(c *gin.Context) (token string, err error) {
	token = c.GetHeader("Authorization")
	if token == "" {
		err = errors.New("no authorization token")
		return
	}
	token = token[7:]
	return
}
