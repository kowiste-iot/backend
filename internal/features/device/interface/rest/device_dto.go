package devicehandler

import (
	"backend/internal/features/device/domain"
)

type CreateDeviceRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type UpdateDeviceRequest struct {
	Name        string `json:"name" binding:"required"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

type DeviceResponse struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Parent      string `json:"parent,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Password    string `json:"password,omitempty"`
	UpdatedAt   int64  `json:"updatedAt"`
}
type DevicePasswordResponse struct {
	ID          string `json:"id"`
	Password    string `json:"password"`
}

func ToDeviceResponse(a *domain.Device) DeviceResponse {
	return DeviceResponse{
		ID:          a.ID(),
		TenantID:    a.TenantID(),
		Name:        a.Name(),
		Parent:      a.Parent(),
		Description: a.Description(),
		UpdatedAt: a.UpdatedAt().Unix(),
	}
}
func ToDevicePasswordResponse(a *domain.Device, password string) DevicePasswordResponse {
	return DevicePasswordResponse{
		ID:          a.ID(),
		Password:    password,
	}
}

func ToDeviceResponses(assets []*domain.Device) []DeviceResponse {
	responses := make([]DeviceResponse, len(assets))
	for i, a := range assets {
		responses[i] = ToDeviceResponse(a)
	}
	return responses
}
