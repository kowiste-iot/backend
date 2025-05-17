package websocket_handler

import (
	"backend/shared/base/command"
	"backend/shared/http/httputil"
	"backend/shared/pagination"
	"backend/shared/websocket/app"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
) // GenerateToken generates a token for WebSocket authentication
func (h *Handler) GenerateToken(c *gin.Context) {
	ctx := c.Request.Context()
	tenant, branchID, err := httputil.GetBase(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tenant: " + err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := httputil.GetUserID(ctx)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	tenantIDStr := tenant.Domain()

	// Generate a WebSocket token
	token, err := h.tokenService.GenerateWebSocketToken(ctx, tenantIDStr, branchID, userID)
	if err != nil {
		h.base.Logger.Error(ctx, err, "Failed to generate WebSocket token",
			map[string]interface{}{
				"tenantID": tenantIDStr,
				"branchID": branchID,
				"userID":   userID,
			})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	host := c.Request.Host
	uri := c.Request.URL.String()
	uriTrim := strings.TrimSuffix(uri, "/token")

	serverURL := fmt.Sprintf("ws://%s%s?token=%s", host, uriTrim, token)
	response := pagination.PaginatedResponse{
		Data: gin.H{
			"url": serverURL,
		},
	}
	// Return the token to the client
	c.JSON(http.StatusOK, response)
}

// HandleWebSocket handles WebSocket connection requests
func (h *Handler) HandleWebSocket(c *gin.Context) {
	ctx := c.Request.Context()

	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}
	//TODO:Should check tenant exist?
	tenant := c.Query("tenant")
	if tenant == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant is required"})
		return
	}
	branch := c.Query("branch")
	if branch == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Branch is required"})
		return
	}

	// Validate the token
	userID, err := h.tokenService.ValidateToken(ctx,
		command.BaseInput{TenantDomain: tenant, BranchName: branch},
		token)
	if err != nil {
		h.base.Logger.Error(ctx, errors.New("not valid token"), "Invalid WebSocket token",
			map[string]interface{}{
				"token": token,
			})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Upgrade the HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.base.Logger.Error(ctx, err, "Failed to upgrade WebSocket connection",
			map[string]interface{}{
				"tenantID": tenant,
				"userID":   userID,
			})
		return
	}

	// Configure connection
	conn.SetReadLimit(32768) // 32KB max message size
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Create and register client
	client := app.NewClient(h.hub, conn, tenant+branch, userID)
	h.hub.RegisterClient(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}
