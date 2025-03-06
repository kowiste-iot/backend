package websocket_handler

import (
	"backend/shared/http/httputil"
	"backend/shared/websocket/app"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
) // GenerateToken generates a token for WebSocket authentication
func (h *Handler) GenerateToken(c *gin.Context) {
	ctx := c.Request.Context()
	tenant, _, err := httputil.GetBase(ctx)
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
	token, err := h.tokenService.GenerateWebSocketToken(ctx, tenantIDStr, userID)
	if err != nil {
		h.base.Logger.Error(ctx, err, "Failed to generate WebSocket token",
			map[string]interface{}{
				"tenantID": tenantIDStr,
				"userID":   userID,
			})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token to the client
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
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
	tenant, ok := httputil.GetTenant(ctx)
	if !ok {
		h.base.Logger.Error(ctx, errors.New("cannot get tenant"), "Failed to upgrade WebSocket connection",
			map[string]interface{}{})
		return
	}
	// Validate the token
	err := h.tokenService.ValidateToken(ctx, token)
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
	client := app.NewClient(h.hub, conn, tenant.Domain(), "pablo")
	h.hub.RegisterClient(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}
