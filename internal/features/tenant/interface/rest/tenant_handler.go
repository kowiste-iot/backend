package tenanthandler

import (
	"net/http"

	"backend/internal/features/tenant/domain"
	"backend/internal/features/tenant/domain/command"
	ginhelp "backend/shared/http/gin"
	"backend/shared/pagination"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new tenant
// @Description Create a new tenant for the tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param tenant body CreateTenantRequest true "Tenant creation request"
// @Success 201 {object} domain.Tenant
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	ctx := c.Request.Context()
	h.base.Logger.Debug(ctx, "Starting tenant creation request", nil)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.base.Logger.Error(ctx, err, "Failed to bind JSON request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := command.CreateTenantInput{
		Name:        req.Name,
		Domain:      req.Domain,
		Description: req.Description,
		AdminEmail:  req.Email,
		Branch:      req.Branch,
	}

	result, err := h.tenantService.CreateTenant(ctx, &input)
	if err != nil {
		h.base.Logger.Error(ctx, err, "Failed to create tenant", map[string]interface{}{
			"domain": input.Domain,
			"name":   input.Name,
			"err":    err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create tenant",
			"data":  err.Error(),
		})
		return
	}
	h.base.Logger.Info(ctx, "Tenant created successfully", map[string]interface{}{
		"tenantID": result.ID(), // Assuming you have an ID getter
		"domain":   result.Domain(),
	})
	c.JSON(http.StatusCreated, ToTenantResponse(result))
}

// @Summary Get a tenant by ID
// @Description Get a tenant by its ID for the tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Tenant ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} domain.Tenant
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/tenants/{id} [get]
func (h *TenantHandler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")
	ctx := c.Request.Context()

	result, err := h.tenantService.GetTenant(ctx, tenantID)
	if err != nil {
		if err == domain.ErrAssetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
			return
		}
		h.base.Logger.Error(c.Request.Context(), err, "Failed to get tenant", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant"})
		return
	}

	c.JSON(http.StatusOK, ToTenantResponse(result))
}

// @Summary List all tenants
// @Description List all tenants for the tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Success 200 {object} struct{tenants []domain.Tenant} "Array of tenants wrapped in tenants field"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/tenants [get]
func (h *TenantHandler) ListTenants(c *gin.Context) {
	ctx := c.Request.Context()

	ctx = ginhelp.SetPaginationGin(ctx, c)

	results, err := h.tenantService.ListTenants(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tenants"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToTenantResponses(results),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update a tenant
// @Description Update an existing tenant by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Tenant ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param tenant body UpdateTenantRequest true "Tenant update request"
// @Success 200 {object} domain.Tenant
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	tenantID := c.Param("id")
	ctx := c.Request.Context()

	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := command.UpdateTenantInput{
		ID:          tenantID,
		Name:        req.Name,
		Description: req.Description,
	}
	result, err := h.tenantService.UpdateTenant(ctx, &input)
	if err != nil {
		if err == domain.ErrAssetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
			return
		}

		h.base.Logger.Error(c.Request.Context(), err, "Failed to update tenant", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tenant"})
		return
	}

	c.JSON(http.StatusOK, ToTenantResponse(result))
}

// @Summary Delete a tenant
// @Description Delete a tenant by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Tenant ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	tenantID := c.Param("id")
	ctx := c.Request.Context()

	err := h.tenantService.DeleteTenant(ctx, tenantID)
	if err != nil {
		if err == domain.ErrAssetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
			return
		}
		h.base.Logger.Error(c.Request.Context(), err, "Failed to delete tenant", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
