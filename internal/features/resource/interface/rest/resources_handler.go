package resourcehandler

import (
	"backend/internal/features/resource/app"

	"backend/internal/features/resource/domain/command"
	baseCmd "backend/shared/base/command"
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/logger"
	"backend/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	logger   logger.Logger
	resource app.ResourceService
}

type Dependencies struct {
	Logger   logger.Logger
	Resource app.ResourceService
}

func New(deps Dependencies) *ResourceHandler {
	return &ResourceHandler{
		logger:   deps.Logger,
		resource: deps.Resource,
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
	resources, err := h.resource.ListResources(ctx, &inputBase)
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
func (h *ResourceHandler) UpdateResource(c *gin.Context) {
	resourceID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	var req UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateResourceInput{
		BaseInput:   inputBase,
		ID:          resourceID,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Roles:       req.Roles,
	}
	result, err := h.resource.UpdateResource(ctx, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, ToResourcesResponse(*result))
}
