package rolehandler

import (
	authApp "ddd/shared/auth/app"
	"ddd/shared/auth/domain/command"
	authCmd "ddd/shared/auth/domain/command"
	baseCmd "ddd/shared/base/command"
	ginhelp "ddd/shared/http/gin"
	"ddd/shared/http/httputil"
	"ddd/shared/logger"
	"ddd/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	logger      logger.Logger
	authService *authApp.Service
}

type Dependencies struct {
	Logger      logger.Logger
	AuthService *authApp.Service
}

func New(deps Dependencies) *RoleHandler {
	return &RoleHandler{
		logger:      deps.Logger,
		authService: deps.AuthService,
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

	input := authCmd.CreateRoleInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		Name:      req.Name,
		Description: req.Description,
	}

	_, err = h.authService.CreateRole(ctx, &input)
	if err != nil {
		h.logger.Error(ctx, err, "Failed to create role", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
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
	role, err := h.authService.GetRole(c.Request.Context(), &input)
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
	roles, err := h.authService.GetRoles(ctx, &inputBase)
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
	err = h.authService.DeleteRole(c.Request.Context(), &input)
	if err != nil {
		h.logger.Error(c.Request.Context(), err, "Failed to delete role", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Additional handlers for role assignments if needed
func (h *RoleHandler) AssignRole(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("userId")
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := authCmd.AssignRolesInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		UserID:    userID,
		Roles:     []string{req.Name},
	}
	err = h.authService.AssignRoles(c.Request.Context(), &input)
	if err != nil {
		h.logger.Error(c.Request.Context(), err, "Failed to assign role", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	c.Status(http.StatusNoContent)
}
