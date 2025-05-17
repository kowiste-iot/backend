package devicehandler

import (
	"net/http"

	"backend/internal/features/device/domain"
	"backend/internal/features/device/domain/command"
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
// @Router /api/v1/{tenantid}/devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req CreateDeviceRequest
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
	input := command.CreateDeviceInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		Name:        req.Name,
		Description: req.Description,
		Parent:      req.Parent,
	}

	result, password, err := h.deviceService.CreateDevice(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to create measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create measure"})
		return
	}

	c.JSON(http.StatusCreated,gin.H{"data":ToDevicePasswordResponse(result, password) } )
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
// @Router /api/v1/{tenantid}/devices/{id} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	deviceID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.DeviceIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		DeviceID:  deviceID,
	}
	result, err := h.deviceService.GetDevice(ctx, &input)
	if err != nil {
		if err == domain.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measure not found"})
			return
		}

		h.base.Logger.Error(c.Request.Context(), err, "Failed to get measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenant.Domain(),
			"deviceID": deviceID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get measure"})
		return
	}

	c.JSON(http.StatusOK, ToDeviceResponse(result))
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
// @Router /api/v1/{tenantid}/devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := baseCmd.NewInput(tenant.Domain(), branch)
	results, err := h.deviceService.ListDevices(ctx, &input)
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
		Data:       ToDeviceResponses(results),
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
// @Router /api/v1/{tenantid}/devices/{id} [put]

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	deviceID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateDeviceInput{
		BaseInput:   inputBase,
		ID:          deviceID,
		Name:        req.Name,
		Parent:      req.Parent,
		Description: req.Description,
	}
	result, err := h.deviceService.UpdateDevice(ctx, &input)
	if err != nil {
		if err == domain.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measure not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to update measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"deviceID": deviceID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update measure"})
		return
	}

	c.JSON(http.StatusOK, ToDeviceResponse(result))
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
// @Router /api/v1/{tenantid}/devices/{id} [delete]

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	deviceID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.DeviceIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		DeviceID:  deviceID,
	}
	err = h.deviceService.DeleteDevice(ctx, &input)
	if err != nil {
		if err == domain.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "measure not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to delete measure", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"deviceID": deviceID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
