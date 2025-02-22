package middleware

import (
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/keycloak"
	"backend/shared/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MiddlewareManager struct {
	logger logger.Logger
	auth   *keycloak.Keycloak
}

func NewMiddlewareManager(logger logger.Logger, auth *keycloak.Keycloak) *MiddlewareManager {
	return &MiddlewareManager{
		logger: logger,
		auth:   auth,
	}
}

func (m *MiddlewareManager) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.String(), "tenant") {
			c.Next()
			return
		}
		token, err := ginhelp.GetAuhtHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		jwtToken, err := m.auth.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		userID, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		ctx := httputil.SetToken(c.Request.Context(), token)
		ctx = httputil.SetUserID(ctx, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (m *MiddlewareManager) Recovery() gin.HandlerFunc {
	return httputil.RecoveryMiddleware(m.logger)
}
