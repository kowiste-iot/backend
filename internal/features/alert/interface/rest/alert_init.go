package alerthandler

import (
	"backend/internal/features/alert/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	base         *base.BaseService
	alertService app.AlertService
}

func New(base *base.BaseService, alerService app.AlertService) *AlertHandler {
	return &AlertHandler{
		base:         base,
		alertService: alerService,
	}
}
func (ah *AlertHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	alerts := rg.Group("alerts")
	{
		alerts.POST("", ah.CreateAlert)
		alerts.GET("", ah.ListAlerts)
		alerts.GET(":id", ah.GetAlert)
		alerts.PUT(":id", ah.UpdateAlert)
		alerts.DELETE(":id", ah.DeleteAlert)
	}

	return alerts
}
