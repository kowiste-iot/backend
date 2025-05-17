package ingesthandler

import (
	"backend/internal/features/ingest/domain"
	"backend/shared/http/httputil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Ingest data
// @Description Ingest data into the system
// @Tags ingest
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID" example:"org123"
// @Param data body IngestDataRequest true "Data to ingest"
// @Success 202 {object} IngestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/ingest [post]
func (h *IngestHandler) IngestData(c *gin.Context) {
	var req IngestDataRequest
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

	msg := &domain.Message{
		ID:       req.ID,
		TenantID: tenant.Domain(),
		BranchID: branch,
		Time:     time.Now(),
		Data:     req.Data,
	}

	if err := h.ingestService.Start(); err != nil {
		h.base.Logger.Error(ctx, err, "Failed to start ingest service", map[string]interface{}{
			"error":    err.Error(),
			"tenantID": tenant.Domain(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start ingest service"})
		return
	}

	c.JSON(http.StatusAccepted, ToIngestResponse(msg))
}
