package resourcehandler

import (
	"backend/internal/features/resource/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	base     *base.BaseService
	resource app.ResourceService
}

func New(base *base.BaseService, resourceService app.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		base:     base,
		resource: resourceService,
	}
}

func (rh *ResourceHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	resource := rg.Group("resources")
	{
		resource.GET("", rh.ListResources)
		resource.PUT(":id", rh.UpdateResource)
	}
	return resource
}
