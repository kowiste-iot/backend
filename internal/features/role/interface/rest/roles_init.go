package rolehandler

import (
	"backend/internal/features/user/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	base        *base.BaseService
	roleService app.RoleService
}

func New(base *base.BaseService, roleService app.RoleService) *RoleHandler {
	return &RoleHandler{
		base:        base,
		roleService: roleService,
	}
}

func (rh *RoleHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	roles := rg.Group("roles")
	{
		roles.POST("", rh.CreateRole)
		roles.GET("", rh.ListRoles)
		roles.GET(":name", rh.GetRole)
		roles.DELETE(":name", rh.DeleteRole)
	}
	return roles
}
