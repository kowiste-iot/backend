// interface/rest/ingesthandler/handler.go
package ingesthandler

import (
	"backend/internal/features/ingest/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type IngestHandler struct {
	base          *base.BaseService
	ingestService app.IngestService
}

func New(base *base.BaseService, ingestService app.IngestService) *IngestHandler {
	return &IngestHandler{
		base:          base,
		ingestService: ingestService,
	}
}

func (h *IngestHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	ingest := rg.Group("ingest")
	{
		ingest.POST("", h.IngestData)
	}
	return ingest
}
