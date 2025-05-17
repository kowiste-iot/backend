package userhandler

import (
	"backend/internal/features/user/app"

	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	base        *base.BaseService
	userService app.UserService
}

func New(base *base.BaseService, userService app.UserService) *UserHandler {
	return &UserHandler{
		base:        base,
		userService: userService,
	}
}
func (uh *UserHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	users := rg.Group("users")
	{
		users.POST("", uh.CreateUser)
		users.GET("", uh.ListUsers)
		users.GET(":id", uh.GetUser)
		users.PUT(":id", uh.UpdateUser)
		users.DELETE(":id", uh.DeleteUser)
	}

	return users
}
