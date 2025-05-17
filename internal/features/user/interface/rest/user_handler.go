package userhandler

import (
	"backend/internal/features/user/domain"
	"backend/internal/features/user/domain/command"
	baseCmd "backend/shared/base/command"

	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
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
	input := command.CreateUserInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Roles:     req.Roles,
	}

	result, err := h.userService.CreateUser(c.Request.Context(), &input)
	if err != nil {
		h.base.Logger.Error(c.Request.Context(), err, "Failed to create user", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.UserIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		UserID:    userID,
	}
	result, err := h.userService.GetUser(ctx, &input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.base.Logger.Error(ctx, err, "Failed to get user", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := baseCmd.BaseInput{
		TenantDomain: tenant.Domain(),
		BranchName:   branch,
	}
	results, err := h.userService.ListUsers(ctx, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       results,
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	ctx := c.Request.Context()

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.UpdateUserInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		ID:        userID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	result, err := h.userService.UpdateUser(ctx, &input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.base.Logger.Error(ctx, err, "Failed to update user", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	ctx := c.Request.Context()
	tenant, branch, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get base: " + err.Error()})
		return
	}
	input := command.UserIDInput{
		BaseInput: baseCmd.NewInput(tenant.Domain(), branch),
		UserID:    userID,
	}
	err = h.userService.DeleteUser(ctx, &input)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		h.base.Logger.Error(ctx, err, "Failed to delete user", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
