package websocket

// import (
// 	domainToken "backend/shared/token/domain"
// 	appToken "backend/shared/token/app"
// 	appWS "backend/shared/websocket/app"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// type NotificationHandler struct {
// 	hub *appWS.Hub
// 	// subscriber   streamingToken.Subscriber
// 	tokenService appToken.TokenService
// 	upgrader     websocket.Upgrader
// }

// func NewNotificationHandler(hub *appWS.Hub, subscriber streamingToken.Subscriber, tokenService domainToken.TokenService) *NotificationHandler {
// 	return &NotificationHandler{
// 		hub: hub,
// 		// subscriber:   subscriber,
// 		tokenService: tokenService,
// 		upgrader: websocket.Upgrader{
// 			ReadBufferSize:  1024,
// 			WriteBufferSize: 1024,
// 			CheckOrigin: func(r *http.Request) bool {
// 				// Implement your origin check logic
// 				return true
// 			},
// 		},
// 	}
// }

// func (h *NotificationHandler) HandleWebSocket(c *gin.Context) {
// 	token := c.Query("token")
// 	if token == "" {
// 		token = c.GetHeader("X-WS-Token")
// 	}

// 	claims, err := h.tokenService.ValidateWSToken(token) // Change this
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
// 		return
// 	}

// 	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, gin.H{"error": "cannot upgrade" + err.Error()})
// 		return
// 	}
// 	// ctx := c.Request.Context()
// 	client := appWS.NewClient(h.hub, conn, claims.TenantID, claims.UserID)
// 	// subGen := subjects.NewNotificationSubject()
// 	// wsHandler := appWS.NewWebSocketHandler(h.hub)
// 	// tenant, ok := httputil.GetTenant(ctx)
// 	// if !ok {
// 	// 	c.JSON(http.StatusBadGateway, gin.H{"error": "not tenant id"})
// 	// 	return
// 	// }
// 	// userID, ok := httputil.GetUserID(ctx)
// 	// if !ok {
// 	// 	c.JSON(http.StatusBadGateway, gin.H{"error": "not user id"})
// 	// 	return
// 	// }
// 	// // Subscribe to user's notifications
// 	// err = h.subscriber.Subscribe(ctx, subGen, wsHandler, tenant.ID(), userID)
// 	// if err != nil {
// 	// 	conn.Close()
// 	// 	return
// 	// }

// 	h.hub.RegisterClient(client) // Using the exported method instead

// 	// Start client goroutines
// 	go client.WritePump()
// 	go client.ReadPump()
// }
