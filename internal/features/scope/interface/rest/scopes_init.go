package scopehandler

import (
	"backend/internal/features/scope/app"
	"backend/shared/base"

	"github.com/gin-gonic/gin"
)

type ScopeHandler struct {
	base         *base.BaseService
	scopeService app.ScopeService
}

func New(base *base.BaseService, scopeService app.ScopeService) *ScopeHandler {
	return &ScopeHandler{
		base:         base,
		scopeService: scopeService,
	}
}
func (sh *ScopeHandler) Init(rg *gin.RouterGroup) *gin.RouterGroup {

	scopes := rg.Group("scopes")
	{
		scopes.GET("", sh.ListRoles)
	}
	return scopes
}
