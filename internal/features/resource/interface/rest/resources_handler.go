package resourcehandler

import (
	"backend/internal/features/resource/domain/command"
	baseCmd "backend/shared/base/command"
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	_, err = h.resource.UpdateResource(ctx, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}

	c.Status(http.StatusOK)
}
