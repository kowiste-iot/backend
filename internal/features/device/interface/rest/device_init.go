package devicehandler

import (
	"backend/internal/features/device/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	base          *base.BaseService
	deviceService app.DeviceService
}

func New(base *base.BaseService, deviceService app.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		base:          base,
		deviceService: deviceService,
	}
}

func (dh *DeviceHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	devices := rg.Group("devices")
	{
		devices.POST("", dh.CreateDevice)
		devices.GET("", dh.ListDevices)
		devices.GET(":id", dh.GetDevice)
		devices.PUT(":id", dh.UpdateDevice)
		devices.DELETE(":id", dh.DeleteDevice)
	}
	return devices
}
