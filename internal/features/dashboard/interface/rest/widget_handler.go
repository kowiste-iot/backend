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

// @Summary Create a new widget
// @Description Create a new widget for the tenant
// @Tags widgets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param widget body CreateWidgetRequest true "Widget creation request"
// @Success 201 {object} domain.Widget
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{did}/widgets [post]
func (h *WidgetHandler) CreateWidget(c *gin.Context) {
	dashboardID := c.Param("id")

	var req CreateWidgetRequest
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
	ctx = ginhelp.SetPaginationGin(ctx, c)
	// Convert from request link data to command link data
	// Assuming your command package has a WidgetLink type
	commandLinks := make([]command.WidgetLink, len(req.Data.Link))
	for i, link := range req.Data.Link {
		commandLinks[i] = command.WidgetLink{
			Measure: link.Measure,
			Tag:     link.Tag,
			Legend:  link.Legend,
		}
	}

	// Create input for service
	input := command.CreateWidgetInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		DashboardID: dashboardID, TypeWidget: req.Type,
		Label:       req.Data.Label,
		ShowLabel:   req.Data.ShowLabel,
		ShowEmotion: req.Data.ShowEmotion,
		TrueEmotion: req.Data.TrueEmotion,
		Link:        commandLinks, // Using the command link type here
		Options:     req.Data.Options,
	}

	result, err := h.widgetService.CreateWidget(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to create widget", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenantID,
			"dashboardID": dashboardID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create widget"})
		return
	}

	c.JSON(http.StatusCreated, ToWidgetResponse(result, 0))
}

// @Summary Get a widget by ID
// @Description Get a widget by its ID for the tenant
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param did path string true "Dashboard ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param id path string true "Widget ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} domain.Widget
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{did}/widgets/{id} [get]
func (h *WidgetHandler) GetWidget(c *gin.Context) {
	dashboardID := c.Param("id")
	widgetID := c.Param("wid")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.WidgetIDInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		DashboardID: dashboardID,
		WidgetID:    widgetID,
	}
	result, err := h.widgetService.GetWidget(ctx, &input)
	if err != nil {
		if err == domain.ErrDashboardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
			return
		}

		h.base.Logger.Error(c.Request.Context(), err, "Failed to get widget", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenant.Domain(),
			"dashboardID": dashboardID,
			"widgetID":    widgetID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get widget"})
		return
	}

	c.JSON(http.StatusOK, ToWidgetResponse(result, 0))
}

// @Summary List all widgets in dashboard
// @Description List all widgets for the dashboard
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Success 200 {object} struct{dashboards []domain.Widget} "Array of widget wrapped in dashboards field"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{id}/widgets [get]
func (h *WidgetHandler) ListWidgets(c *gin.Context) {
	dashboardID := c.Param("id")
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.DashboardIDInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		DashboardID: dashboardID,
	}
	results, err := h.widgetService.ListWidgets(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to get widget", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenantID,
			"dashboardID": dashboardID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list widgets"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToWidgetResponses(results),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update an widget
// @Description Update an existing widget by ID
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Widget ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param did path string true "Dashboard ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param dashboard body UpdateDashboardRequest true "Widget update request"
// @Success 200 {object} domain.Widget
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{did}/widgets/{id} [put]

func (h *WidgetHandler) UpdateWidget(c *gin.Context) {
	dashboardID := c.Param("id")
	widgetID := c.Param("wid")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}

	var req UpdateWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert from request link data to command link data
	commandLinks := make([]command.WidgetLink, len(req.Data.Link))
	for i, link := range req.Data.Link {
		commandLinks[i] = command.WidgetLink{
			Measure: link.Measure,
			Tag:     link.Tag,
			Legend:  link.Legend,
		}
	}

	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateWidgetInput{
		BaseInput:   inputBase,
		ID:          widgetID,
		DashboardID: dashboardID,
		TypeWidget:  req.Type,
		X:           req.X, Y: req.Y,
		W: req.W, H: req.H,
		Label:       req.Data.Label,
		ShowLabel:   req.Data.ShowLabel,
		ShowEmotion: req.Data.ShowEmotion,
		TrueEmotion: req.Data.TrueEmotion,
		Link:        commandLinks,
		Options:     req.Data.Options,
	}

	result, err := h.widgetService.UpdateWidget(ctx, &input)
	if err != nil {
		if err == domain.ErrDashboardNotFound || err == domain.ErrWidgetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to update widget", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenantID,
			"dashboardID": dashboardID,
			"widgetID":    widgetID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update widget"})
		return
	}

	c.JSON(http.StatusOK, ToWidgetResponse(result, 0))
}
// @Summary Update a widget position
// @Description Update an existing widget position by ID
// @Tags dashboards
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Widget ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param did path string true "Dashboard ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param dashboard body UpdateDashboardRequest true "Widget update request"
// @Success 200 {object} domain.Widget
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{did}/widgets/{id} [put]

func (h *WidgetHandler) UpdateWidgetPosition(c *gin.Context) {
	dashboardID := c.Param("id")
	widgetID := c.Param("wid")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}

	var req UpdateWidgetPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}



	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateWidgetPositionInput{
		BaseInput:   inputBase,
		ID:          widgetID,
		DashboardID: dashboardID,
		X:           req.X, Y: req.Y,
		W: req.W, H: req.H,
	}

	result, err := h.widgetService.UpdateWidgetPosition(ctx, &input)
	if err != nil {
		if err == domain.ErrDashboardNotFound || err == domain.ErrWidgetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Widget not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to update widget position", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenantID,
			"dashboardID": dashboardID,
			"widgetID":    widgetID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update widget position"})
		return
	}

	c.JSON(http.StatusOK, ToWidgetResponse(result, 0))
}

// @Summary Delete an widget
// @Description Delete an widget by ID
// @Tags widgets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param did path string true "Dashboard ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param id path string true "Widget ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/dashboards/{did}/widgets/{id} [delete]

func (h *WidgetHandler) DeleteWidget(c *gin.Context) {
	dashboardID := c.Param("id")
	widgetID := c.Param("wid")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.WidgetIDInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		DashboardID: dashboardID,
		WidgetID:    widgetID,
	}
	err = h.widgetService.DeleteWidget(ctx, &input)
	if err != nil {
		if err == domain.ErrWidgetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "widget not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to delete widget", map[string]interface{}{
			"error":       err.Error(),
			"tenantID":    tenantID,
			"dashboardID": dashboardID,
			"widgetID":    widgetID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
