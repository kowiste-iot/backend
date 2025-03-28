// // internal/interfaces/http/handlers/websocket/token_handler.go
 package websocket

// import (
// 	"backend/shared/http/httputil"
// 	"backend/shared/token/domain"

// 	"github.com/gin-gonic/gin"
// )

// type TokenHandler struct {
// 	tokenService domain.TokenService
// }

// func NewTokenHandler(tokenService domain.TokenService) *TokenHandler {
// 	return &TokenHandler{
// 		tokenService: tokenService,
// 	}
// }

// func (h *TokenHandler) GenerateToken(c *gin.Context) {
// 	tenant, ok := httputil.GetTenant(c.Request.Context())
// 	if !ok {
// 		c.JSON(400, gin.H{"error": "tenant id not found"})
// 		return
// 	}

// 	userID, ok := httputil.GetUserID(c.Request.Context())
// 	if !ok {
// 		c.JSON(400, gin.H{"error": "user id not found"})
// 		return
// 	}

// 	token, err := h.tokenService.GenerateWSToken(tenant.ID(), userID)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "could not generate token"})
// 		return
// 	}

// 	c.JSON(200, gin.H{"token": token})
// }
