package dashboardhandler

import (
	"net/http"

	"backend/internal/features/dashboard/domain"
	"backend/internal/features/dashboard/domain/command"
	baseCmd "backend/shared/base/command"
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/pagination"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new dashboard
// @Description Create a new dashboard for the tenant
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param dashboard body CreateDashboardRequest true "Dashboard creation request"
// @Success 201 {object} domain.Dashboard
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards [post]
func (h *DashboardHandler) CreateDashboard(c *gin.Context) {
	var req CreateDashboardRequest
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
	input := command.CreateDashboardInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		Name:        req.Name,
		Description: req.Description,
		Parent:      req.Parent,
	}

	result, err := h.dashboardService.CreateDashboard(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to create dashboard", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dashboard"})
		return
	}

	c.JSON(http.StatusCreated, ToDashboardResponse(result))
}

// @Summary Get an dashboard by ID
// @Description Get an dashboard by its ID for the tenant
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Dashboard ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} domain.Dashboard
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{id} [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	dashboardID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.DashboardIDInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		DashboardID: dashboardID,
	}
	result, err := h.dashboardService.GetDashboard(ctx, &input)
	if err != nil {
		if err == domain.ErrDashboardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
			return
		}

		h.base.Logger.Error(c.Request.Context(), err, "Failed to get dashboard", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenant.Domain(),
			"dashboardID": dashboardID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard"})
		return
	}

	c.JSON(http.StatusOK, ToDashboardResponse(result))
}

// @Summary List all dashboards
// @Description List all dashboards for the tenant
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Success 200 {object} struct{dashboards []domain.Dashboard} "Array of dashboards wrapped in dashboards field"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards [get]
func (h *DashboardHandler) ListDashboards(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := baseCmd.NewInput(tenant.Domain(), branch)
	results, err := h.dashboardService.ListDashboards(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to get dashboard", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list dashboards"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToDashboardResponses(results),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update an dashboard
// @Description Update an existing dashboard by ID
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Dashboard ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param dashboard body UpdateDashboardRequest true "Dashboard update request"
// @Success 200 {object} domain.Dashboard
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{id} [put]

func (h *DashboardHandler) UpdateDashboard(c *gin.Context) {
	dashboardID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	var req UpdateDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateDashboardInput{
		BaseInput:   inputBase,
		ID:          dashboardID,
		Name:        req.Name,
		Parent:      req.Parent,
		Description: req.Description,
	}
	result, err := h.dashboardService.UpdateDashboard(ctx, &input)
	if err != nil {
		if err == domain.ErrDashboardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Dashboard not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to update dashboard", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenantID,
			"dashboardID": dashboardID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dashboard"})
		return
	}

	c.JSON(http.StatusOK, ToDashboardResponse(result))
}

// @Summary Delete an dashboard
// @Description Delete an dashboard by ID
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Dashboard ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{id} [delete]

func (h *DashboardHandler) DeleteDashboard(c *gin.Context) {
	dashboardID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.DashboardIDInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		DashboardID: dashboardID,
	}
	err = h.dashboardService.DeleteDashboard(ctx, &input)
	if err != nil {
		if err == domain.ErrDashboardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "dashboard not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to delete dashboard", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenantID,
			"dashboardID": dashboardID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
