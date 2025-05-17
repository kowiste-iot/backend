package dashboardhandler

import (
	"backend/internal/features/dashboard/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type WidgetHandler struct {
	base          *base.BaseService
	widgetService app.WidgetService
}

func NewWidget(base *base.BaseService, widgetService app.WidgetService) *WidgetHandler {
	return &WidgetHandler{
		base:          base,
		widgetService: widgetService,
	}
}

func (wh *WidgetHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	widgets := rg.Group(":id/widgets")
	{
		widgets.POST("", wh.CreateWidget)
		widgets.GET("", wh.ListWidgets)
		widgets.GET(":wid", wh.GetWidget)
		widgets.PUT(":wid", wh.UpdateWidget)
		widgets.PUT(":wid/position", wh.UpdateWidgetPosition)
		widgets.DELETE(":wid", wh.DeleteWidget)
	}
	return widgets
}
