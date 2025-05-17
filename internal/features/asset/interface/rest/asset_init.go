package assethandler

import (
	"backend/internal/features/asset/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	base         *base.BaseService
	assetService app.AssetService
}

func New(base *base.BaseService, assetService app.AssetService) *AssetHandler {
	return &AssetHandler{
		base:         base,
		assetService: assetService,
	}
}
func (ah *AssetHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	assets := rg.Group("assets")
	{
		assets.POST("", ah.CreateAsset)
		assets.GET("", ah.ListAssets)
		assets.GET(":id", ah.GetAsset)
		assets.PUT(":id", ah.UpdateAsset)
		assets.DELETE(":id", ah.DeleteAsset)
	}
	return assets
}
