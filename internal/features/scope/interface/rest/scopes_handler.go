package scopehandler

import (
	"backend/internal/features/scope/app"
	baseCmd "backend/shared/base/command"

	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/logger"
	"backend/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScopeHandler struct {
	logger       logger.Logger
	scopeService app.ScopeService
}

type Dependencies struct {
	Logger       logger.Logger
	ScopeService app.ScopeService
}

func New(deps Dependencies) *ScopeHandler {
	return &ScopeHandler{
		logger:       deps.Logger,
		scopeService: deps.ScopeService,
	}
}

func (h *ScopeHandler) ListRoles(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	scopes, err := h.scopeService.ListScopes(ctx, &inputBase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToScopeResponses(scopes),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}
