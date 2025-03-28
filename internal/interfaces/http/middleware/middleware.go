package middleware

import (
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/keycloak"
	"backend/shared/logger"
	"errors"
	"net/http"
	"regexp"
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
		if match, _ := regexp.MatchString(`/api/[^/]+/[^/]+/ws$`, c.Request.URL.Path); match {
			c.Next()
			return
		}

		if strings.Contains(c.Request.URL.String(), "tenant") { //by pass for now to allow create tenant without token
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
		userID, err := m.getUserID(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			return
		}
		ctx := httputil.SetToken(c.Request.Context(), token)
		ctx = httputil.SetUserID(ctx, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (m *MiddlewareManager) getUserID(token *jwt.Token) (userID string, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}
	userID, ok = claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token sub")
	}

	return
}
func (m *MiddlewareManager) Recovery() gin.HandlerFunc {
	return httputil.RecoveryMiddleware(m.logger)
}
