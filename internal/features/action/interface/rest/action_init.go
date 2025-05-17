package actionhandler

import (
	"backend/internal/features/action/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type ActionHandler struct {
	base          *base.BaseService
	actionService app.ActionService
}

func New(base *base.BaseService, actionService app.ActionService) *ActionHandler {
	return &ActionHandler{
		base:          base,
		actionService: actionService,
	}
}

func (ah *ActionHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	actions := rg.Group("actions")
	{
		actions.POST("", ah.CreateAction)
		actions.GET("", ah.ListActions)
		actions.GET(":id", ah.GetAction)
		actions.PUT(":id", ah.UpdateAction)
		actions.DELETE(":id", ah.DeleteAction)
	}

	return actions
}
