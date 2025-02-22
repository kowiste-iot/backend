package dashboardhandler

import (
	"backend/internal/features/dashboard/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	base             *base.BaseService
	dashboardService app.DashboardService
}

func New(base *base.BaseService, dashboardService app.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		base:             base,
		dashboardService: dashboardService,
	}

}

func (dh *DashboardHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	dashboards := rg.Group("dashboards")
	{
		dashboards.POST("", dh.CreateDashboard)
		dashboards.GET("", dh.ListDashboards)
		dashboards.GET(":id", dh.GetDashboard)
		dashboards.PUT(":id", dh.UpdateDashboard)
		dashboards.DELETE(":id", dh.DeleteDashboard)

	}
	return dashboards
}
