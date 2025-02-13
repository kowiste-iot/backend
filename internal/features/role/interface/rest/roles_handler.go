package rolehandler

import (
	"backend/internal/features/role/app"
	"backend/internal/features/role/domain/command"
	baseCmd "backend/shared/base/command"
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/logger"
	"backend/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	logger      logger.Logger
	roleService app.RoleService
}

type Dependencies struct {
	Logger      logger.Logger
	RoleService app.RoleService
}

func New(deps Dependencies) *RoleHandler {
	return &RoleHandler{
		logger:      deps.Logger,
		roleService: deps.RoleService,
	}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
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

	input := command.CreateRoleInput{
		BaseInput:   baseCmd.NewInput(tenant.Domain(), branch),
		Name:        req.Name,
		Description: req.Description,
	}

	_, err = h.roleService.CreateRole(ctx, &input)
	if err != nil {
		h.logger.Error(ctx, err, "Failed to create role", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role " + err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	roleName := c.Param("name")
	ctx := c.Request.Context()

	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.RoleIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		RoleID:    roleName,
	}
	role, err := h.roleService.GetRole(c.Request.Context(), &input)
	if err != nil {
		h.logger.Error(c.Request.Context(), err, "Failed to get role", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role"})
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, RoleResponse{
		Name:        role.Name,
		Description: role.Description,
	})
}

func (h *RoleHandler) ListRoles(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	inputBase := baseCmd.NewInput(tenant.Domain(), branch)
	roles, err := h.roleService.ListRoles(ctx, &inputBase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToRoleResponses(roles),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	ctx := c.Request.Context()
	roleName := c.Param("name")
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.RoleIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		RoleID:    roleName,
	}
	err = h.roleService.DeleteRole(c.Request.Context(), &input)
	if err != nil {
		h.logger.Error(c.Request.Context(), err, "Failed to delete role", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.Status(http.StatusNoContent)
}
