package alerthandler

import (
	"net/http"

	"backend/internal/features/alert/domain"
	"backend/internal/features/alert/domain/command"
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
// @Router /api/v1/{tenantid}/alerts [post]
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	var req CreateAlertRequest
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
	input := command.CreateAlertInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		Name:        req.Name,
		Description: req.Description,
		Parent:      req.Parent,
	}

	result, err := h.alertService.CreateAlert(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to create measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create measure"})
		return
	}

	c.JSON(http.StatusCreated, ToAlertResponse(result))
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
// @Router /api/v1/{tenantid}/alerts/{id} [get]
func (h *AlertHandler) GetAlert(c *gin.Context) {
	alertID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.AlertIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		AlertID:   alertID,
	}
	result, err := h.alertService.GetAlert(ctx, &input)
	if err != nil {
		if err == domain.ErrAlertNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measure not found"})
			return
		}

		h.base.Logger.Error(c.Request.Context(), err, "Failed to get measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenant.Domain(),
			"alertID":  alertID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get measure"})
		return
	}

	c.JSON(http.StatusOK, ToAlertResponse(result))
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
// @Router /api/v1/{tenantid}/alerts [get]
func (h *AlertHandler) ListAlerts(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := baseCmd.NewInput(tenant.Domain(), branch)
	results, err := h.alertService.ListAlerts(ctx, &input)
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
		Data:       ToAlertResponses(results),
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
// @Router /api/v1/{tenantid}/alerts/{id} [put]

func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	alertID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	var req UpdateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateAlertInput{
		BaseInput:   inputBase,
		ID:          alertID,
		Name:        req.Name,
		Parent:      req.Parent,
		Description: req.Description,
	}
	result, err := h.alertService.UpdateAlert(ctx, &input)
	if err != nil {
		if err == domain.ErrAlertNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measure not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to update measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"alertID":  alertID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update measure"})
		return
	}

	c.JSON(http.StatusOK, ToAlertResponse(result))
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
// @Router /api/v1/{tenantid}/alerts/{id} [delete]

func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	alertID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.AlertIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		AlertID:   alertID,
	}
	err = h.alertService.DeleteAlert(ctx, &input)
	if err != nil {
		if err == domain.ErrAlertNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "measure not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to delete measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"alertID":  alertID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
