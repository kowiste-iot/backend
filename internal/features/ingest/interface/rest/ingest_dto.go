package ingesthandler

import (
    "backend/internal/features/ingest/domain"
)

type IngestDataRequest struct {
    ID   string                 `json:"id" binding:"required"`
    Data map[string]interface{} `json:"data" binding:"required"`
}

type IngestResponse struct {
    ID        string                 `json:"id"`
    TenantID  string                 `json:"tenantId"`
    BranchID  string                 `json:"branchId"`
    Timestamp int64                  `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
}

func ToIngestResponse(m *domain.Message) IngestResponse {
    return IngestResponse{
        ID:        m.ID,
        TenantID:  m.TenantID,
        BranchID:  m.BranchID,
        Timestamp: m.Time.Unix(),
        Data:      m.Data,
    }
}