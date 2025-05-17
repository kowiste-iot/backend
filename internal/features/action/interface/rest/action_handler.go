package actionhandler

import (
	"net/http"

	"backend/internal/features/action/domain"
	"backend/internal/features/action/domain/command"
	baseCmd "backend/shared/base/command"
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/pagination"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new measure
// @Description Create a new measure for the tenant
// @Tags measures
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param measure body CreateMeasureRequest true "Measure creation request"
// @Success 201 {object} domain.Measure
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/actions [post]
func (h *ActionHandler) CreateAction(c *gin.Context) {
	var req CreateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.CreateActionInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		Name:        req.Name,
		Description: req.Description,
		Parent:      req.Parent,
	}

	result, err := h.actionService.CreateAction(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to create measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create measure"})
		return
	}

	c.JSON(http.StatusCreated, ToActionResponse(result))
}

// @Summary Get an measure by ID
// @Description Get an measure by its ID for the tenant
// @Tags measures
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Measure ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} domain.Measure
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/actions/{id} [get]
func (h *ActionHandler) GetAction(c *gin.Context) {
	actionID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.ActionIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		ActionID:  actionID,
	}
	result, err := h.actionService.GetAction(ctx, &input)
	if err != nil {
		if err == domain.ErrActionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measure not found"})
			return
		}

		h.base.Logger.Error(c.Request.Context(), err, "Failed to get measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenant.Domain(),
			"actionID": actionID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get measure"})
		return
	}

	c.JSON(http.StatusOK, ToActionResponse(result))
}

// @Summary List all measures
// @Description List all measures for the tenant
// @Tags measures
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Success 200 {object} struct{measures []domain.Measure} "Array of measures wrapped in measures field"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/actions [get]
func (h *ActionHandler) ListActions(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := baseCmd.NewInput(tenant.Domain(), branch)
	results, err := h.actionService.ListActions(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to get measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list measures"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToActionResponses(results),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update an measure
// @Description Update an existing measure by ID
// @Tags measures
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Measure ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param measure body UpdateMeasureRequest true "Measure update request"
// @Success 200 {object} domain.Measure
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/actions/{id} [put]

func (h *ActionHandler) UpdateAction(c *gin.Context) {
	actionID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	var req UpdateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateActionInput{
		BaseInput:   inputBase,
		ID:          actionID,
		Name:        req.Name,
		Parent:      req.Parent,
		Description: req.Description,
	}
	result, err := h.actionService.UpdateAction(ctx, &input)
	if err != nil {
		if err == domain.ErrActionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measure not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to update measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"actionID": actionID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update measure"})
		return
	}

	c.JSON(http.StatusOK, ToActionResponse(result))
}

// @Summary Delete an measure
// @Description Delete an measure by ID
// @Tags measures
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Measure ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/actions/{id} [delete]

func (h *ActionHandler) DeleteAction(c *gin.Context) {
	actionID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.ActionIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		ActionID:  actionID,
	}
	err = h.actionService.DeleteAction(ctx, &input)
	if err != nil {
		if err == domain.ErrActionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "measure not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to delete measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"actionID": actionID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
