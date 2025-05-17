package tenanthandler

import (
	"backend/internal/features/tenant/app"
	"backend/internal/interfaces/http/middleware"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	base          *base.BaseService
	tenantService app.TenantService
	middleware    *middleware.MiddlewareManager
}

func New(base *base.BaseService, tenantService app.TenantService, middleware *middleware.MiddlewareManager) *TenantHandler {
	return &TenantHandler{
		base:          base,
		tenantService: tenantService,
		middleware:    middleware,
	}
}

func (th *TenantHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	tenant := rg.Group("tenant")
	{
		tenant.Use(th.middleware.Auth())
		tenant.POST("", th.CreateTenant)
		tenant.GET("", th.ListTenants)
		tenant.GET(":id", th.GetTenant)
		tenant.PUT(":id", th.UpdateTenant)
		tenant.DELETE(":id", th.DeleteTenant)
	}

	return tenant
}
