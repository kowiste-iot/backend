package measurehandler

import (
	"backend/internal/features/measure/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type MeasureHandler struct {
	base           *base.BaseService
	measureService app.MeasureService
}

func New(base *base.BaseService, measureService app.MeasureService) *MeasureHandler {
	return &MeasureHandler{
		base:           base,
		measureService: measureService,
	}
}
func (mh *MeasureHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	measures := rg.Group("measures")
	{
		measures.POST("", mh.CreateMeasure)
		measures.GET("", mh.ListMeasures)
		measures.GET(":id", mh.GetMeasure)
		measures.PUT(":id", mh.UpdateMeasure)
		measures.DELETE(":id", mh.DeleteMeasure)
	}
	return measures
}
