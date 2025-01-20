package resourcehandler

import (
	authApp "ddd/shared/auth/app"
	baseCmd "ddd/shared/base/command"
	ginhelp "ddd/shared/http/gin"
	"ddd/shared/http/httputil"
	"ddd/shared/logger"
	"ddd/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	logger      logger.Logger
	authService *authApp.Service
}

type Dependencies struct {
	Logger      logger.Logger
	AuthService *authApp.Service
}

func New(deps Dependencies) *ResourceHandler {
	return &ResourceHandler{
		logger:      deps.Logger,
		authService: deps.AuthService,
	}
}

func (h *ResourceHandler) ListResources(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	resources, err := h.authService.GetResources(ctx, &inputBase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToResourcesResponses(resources),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}
