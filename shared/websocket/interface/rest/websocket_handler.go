package websocket_handler

import (
	"backend/shared/base"
	"backend/shared/http/httputil"
	"backend/shared/websocket/app"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Handler handles WebSocket HTTP requests
type Handler struct {
	base     *base.BaseService
	hub      *app.Hub
	service  *app.WebSocketStreamService
	upgrader websocket.Upgrader
}

// New creates a new WebSocket handler
func New(base *base.BaseService, hub *app.Hub, service *app.WebSocketStreamService) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development; customize for production
		},
	}

	return &Handler{
		base:     base,
		hub:      hub,
		service:  service,
		upgrader: upgrader,
	}
}

// Init initializes the WebSocket routes
func (h *Handler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	ws := rg.Group("ws")
	{
		ws.GET("", h.HandleWebSocket)
	}
	return ws
}

// HandleWebSocket handles WebSocket connection requests
func (h *Handler) HandleWebSocket(c *gin.Context) {
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

	// Upgrade the HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.base.Logger.Error(ctx, err, "Failed to upgrade WebSocket connection",
			map[string]interface{}{
				"tenantID": tenantIDStr,
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
	client := app.NewClient(h.hub, conn, tenantIDStr, userID)
	h.hub.RegisterClient(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}
