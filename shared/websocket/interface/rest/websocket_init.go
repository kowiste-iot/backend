package websocket_handler

import (
	"backend/shared/base"
	"backend/shared/token/app"
	wsapp "backend/shared/websocket/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Handler handles WebSocket HTTP requests
type Handler struct {
	base         *base.BaseService
	hub          *wsapp.Hub
	service      *wsapp.WebSocketStreamService
	tokenService *app.TokenService
	upgrader     websocket.Upgrader
}

// New creates a new WebSocket handler
func New(base *base.BaseService, hub *wsapp.Hub, service *wsapp.WebSocketStreamService, tokenService *app.TokenService) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development; customize for production
		},
	}

	return &Handler{
		base:         base,
		hub:          hub,
		service:      service,
		tokenService: tokenService,
		upgrader:     upgrader,
	}
}

// Init initializes the WebSocket routes
func (h *Handler) Init(rg *gin.RouterGroup) *gin.RouterGroup {
	ws := rg.Group("ws")
	{
		ws.POST("/token", h.GenerateToken)
		ws.GET("", h.HandleWebSocket)
	}
	return ws
}
