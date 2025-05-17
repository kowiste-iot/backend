package assethandler

import (
	"net/http"

	"backend/internal/features/asset/domain"
	"backend/internal/features/asset/domain/command"
	baseCmd "backend/shared/base/command"
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/pagination"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new asset
// @Description Create a new asset for the tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param asset body CreateAssetRequest true "Asset creation request"
// @Success 201 {object} domain.Asset
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/assets [post]
func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var req CreateAssetRequest
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
	input := command.CreateAssetInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		Name:        req.Name,
		Description: req.Description,
		Parent:      req.Parent,
	}

	result, err := h.assetService.CreateAsset(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to create asset", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset"})
		return
	}

	c.JSON(http.StatusCreated, ToAssetResponse(result))
}

// @Summary Get an asset by ID
// @Description Get an asset by its ID for the tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Asset ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} domain.Asset
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/assets/{id} [get]
func (h *AssetHandler) GetAsset(c *gin.Context) {
	assetID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.AssetIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		AssetID:   assetID,
	}
	result, err := h.assetService.GetAsset(ctx, &input)
	if err != nil {
		if err == domain.ErrAssetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
			return
		}

		h.base.Logger.Error(c.Request.Context(), err, "Failed to get asset", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenant.Domain(),
			"assetID":  assetID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get asset"})
		return
	}

	c.JSON(http.StatusOK, ToAssetResponse(result))
}

// @Summary List all assets
// @Description List all assets for the tenant
// @Tags assets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Success 200 {object} struct{assets []domain.Asset} "Array of assets wrapped in assets field"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/assets [get]
func (h *AssetHandler) ListAssets(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := baseCmd.NewInput(tenant.Domain(), branch)
	results, err := h.assetService.ListAssets(ctx, &input)
	if err != nil {
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to get asset", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list assets"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToAssetResponses(results),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update an asset
// @Description Update an existing asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Asset ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Param asset body UpdateAssetRequest true "Asset update request"
// @Success 200 {object} domain.Asset
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/assets/{id} [put]

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	assetID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	var req UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	input := command.UpdateAssetInput{
		BaseInput:   inputBase,
		ID:          assetID,
		Name:        req.Name,
		Parent:      req.Parent,
		Description: req.Description,
	}
	result, err := h.assetService.UpdateAsset(ctx, &input)
	if err != nil {
		if err == domain.ErrAssetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to update asset", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"assetID":  assetID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset"})
		return
	}

	c.JSON(http.StatusOK, ToAssetResponse(result))
}

// @Summary Delete an asset
// @Description Delete an asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param id path string true "Asset ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/assets/{id} [delete]

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	assetID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.AssetIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		AssetID:   assetID,
	}
	err = h.assetService.DeleteAsset(ctx, &input)
	if err != nil {
		if err == domain.ErrAssetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
			return
		}
		tenantID, _ := httputil.GetTenant(ctx)
		h.base.Logger.Error(c.Request.Context(), err, "Failed to delete asset", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenantID,
			"assetID":  assetID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
