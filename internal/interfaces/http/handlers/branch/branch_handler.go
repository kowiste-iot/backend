package branchhandler

import (
	"ddd/internal/features/tenant/app"
	"ddd/internal/features/tenant/domain"
	"ddd/internal/features/tenant/domain/command"
	baseCmd "ddd/shared/base/command"
	ginhelp "ddd/shared/http/gin"
	"ddd/shared/http/httputil"
	"ddd/shared/logger"
	"ddd/shared/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	logger        logger.Logger
	branchService app.BranchService
}

type Dependencies struct {
	Logger        logger.Logger
	BranchService app.BranchService
}

func New(deps Dependencies) *BranchHandler {
	return &BranchHandler{
		logger:        deps.Logger,
		branchService: deps.BranchService,
	}
}

// @Summary Create a new branch
// @Description Create a new branch in the tenant
// @Tags branches
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID"
// @Param branch body CreateBranchRequest true "Branch creation request"
// @Success 201 {object} domain.Branch
// @Failure 400,401,500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/branches [post]
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req CreateBranchRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, err, "Failed to bind JSON request", nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant context not found"})
		return
	}

	input := command.CreateBranchInput{
		TenantDomain: tenant.Domain(),
		Name:         req.Name,
		Description:  req.Description,
	}

	result, err := h.branchService.CreateBranch(ctx, &input)
	if err != nil {
		h.logger.Error(ctx, err, "Failed to create branch", nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create branch"})
		return
	}

	c.JSON(http.StatusCreated, ToBranchResponse(result))
}

// @Summary Get a branch by ID
// @Description Get a branch by its ID
// @Tags branches
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID"
// @Param id path string true "Branch ID"
// @Success 200 {object} domain.Branch
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/branches/{id} [get]
func (h *BranchHandler) GetBranch(c *gin.Context) {
	ctx := c.Request.Context()
	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant"})
		return
	}
	branch, ok := httputil.GetBranch(ctx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch"})
		return
	}
	input := baseCmd.NewInput(tenant.Domain(), branch)
	result, err := h.branchService.GetBranch(ctx, &input)
	if err != nil {
		if err == domain.ErrBranchNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch"})
		return
	}

	c.JSON(http.StatusOK, ToBranchResponse(result))
}

// @Summary List all branches
// @Description List all branches in the tenant
// @Tags branches
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID"
// @Success 200 {array} domain.Branch
// @Failure 400,401,500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/branches [get]
func (h *BranchHandler) ListBranches(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ginhelp.SetPaginationGin(ctx, c)
	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant"})
		return
	}
	results, err := h.branchService.ListBranches(ctx, tenant.Domain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list branches"})
		return
	}

	pg, _ := pagination.GetPagination(ctx)
	response := pagination.PaginatedResponse{
		Data:       ToBranchResponses(results),
		Pagination: *pg,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update a branch
// @Description Update an existing branch
// @Tags branches
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID"
// @Param id path string true "Branch ID"
// @Param branch body UpdateBranchRequest true "Branch update request"
// @Success 200 {object} domain.Branch
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/branches/{id} [put]
func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	ctx := c.Request.Context()
	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant"})
		return
	}
	branch, ok := httputil.GetBranch(ctx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch"})
		return
	}
	var req UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := command.UpdateBranchInput{
		ID:branch,
		TenantDomain: tenant.Domain(),
		Name:        req.Name,
		Description: req.Description,
	}

	result, err := h.branchService.UpdateBranch(ctx, &input)
	if err != nil {
		if err == domain.ErrBranchNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch"})
		return
	}

	c.JSON(http.StatusOK, ToBranchResponse(result))
}

// @Summary Delete a branch
// @Description Delete a branch by ID
// @Tags branches
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param tenantid path string true "Tenant ID"
// @Param id path string true "Branch ID"
// @Success 204 "No Content"
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /api/v1/{tenantid}/branches/{id} [delete]
func (h *BranchHandler) DeleteBranch(c *gin.Context) {

	ctx := c.Request.Context()
	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant"})
		return
	}
	branch, ok := httputil.GetBranch(ctx)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch"})
		return
	}
	input := baseCmd.NewInput(tenant.Domain(), branch)
	err := h.branchService.DeleteBranch(ctx, &input)
	if err != nil {
		if err == domain.ErrBranchNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
		return
	}

	c.Status(http.StatusNoContent)
}
