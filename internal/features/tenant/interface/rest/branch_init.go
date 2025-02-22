package tenanthandler

import (
	"backend/internal/features/tenant/app"
	"backend/internal/interfaces/http/middleware"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	base          *base.BaseService
	branchService app.BranchService
	middleware    *middleware.MiddlewareManager
}

func NewBranch(base *base.BaseService, branchService app.BranchService, middleware *middleware.MiddlewareManager) *BranchHandler {
	return &BranchHandler{
		base:          base,
		branchService: branchService,
		middleware:    middleware,
	}
}

func (bh *BranchHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	branchID := rg.Group(":branchid")
	{
		branchID.Use(bh.middleware.Recovery())
		branchID.Use(bh.middleware.Auth())
	
	}

	return branchID
}
