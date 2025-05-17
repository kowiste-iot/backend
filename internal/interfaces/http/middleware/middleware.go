package middleware

import (
	"backend/internal/features/tenant/app"
	"backend/shared/errors"
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
	"backend/shared/keycloak"
	"backend/shared/logger"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MiddlewareManager struct {
	logger        logger.Logger
	tenantService app.TenantService

	auth *keycloak.Keycloak
}

func NewMiddlewareManager(logger logger.Logger, auth *keycloak.Keycloak, tenantService app.TenantService) *MiddlewareManager {
	return &MiddlewareManager{
		logger:        logger,
		tenantService: tenantService,
		auth:          auth,
	}
}

func (m *MiddlewareManager) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.String(), "tenant") { //by pass for now to allow create tenant without token
			c.Next()
			return
		}
		if c.Request.URL.Path == "/api/ws" { //by pass for websocket
			c.Next()
			return
		}

		tenantHeader := c.GetHeader("X-Tenant-ID")
		if tenantHeader == "" {
			httputil.NewErrorResponse(c, errors.NewBadRequest("Branch ID is required", nil))
			c.Abort()
			return
		}
		branchID := c.GetHeader("X-Branch-ID")
		if branchID == "" {
			httputil.NewErrorResponse(c, errors.NewBadRequest("Branch ID is required", nil))
			c.Abort()
			return
		}
		tenant, err := m.tenantService.GetTenant(c.Request.Context(), tenantHeader)
		if err != nil {
			httputil.NewErrorResponse(c, errors.NewBadRequest("Tenant ID not found", err))
			c.Abort()
			return
		}
		ctx := httputil.SetTenant(c.Request.Context(), tenant)
		ctx = httputil.SetBranch(ctx, branchID)
		if match := c.Request.URL.Path == "/api/ws"; match {
			c.Request = c.Request.WithContext(ctx)
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

		jwtToken, err := m.auth.ValidateToken(ctx, tenantHeader, token)
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
		ctx = httputil.SetToken(ctx, token)
		ctx = httputil.SetUserID(ctx, userID)
		t, err := m.getTenant(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			return
		}
		if t != tenantHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "tenant not match",
			})
			return
		}
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func (m *MiddlewareManager) getUserID(token *jwt.Token) (userID string, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}
	userID, ok = claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token sub")
	}

	return
}
func (m *MiddlewareManager) getTenant(token *jwt.Token) (tenant string, err error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// Get issuer
	issuer, ok := claims["iss"].(string)
	if !ok {
		return "", fmt.Errorf("issuer claim not found")
	}

	// Extract realm from issuer URL
	// Format is typically: https://keycloak-server.com/auth/realms/{realm-name}
	parts := strings.Split(issuer, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid issuer format")
	}

	realm := parts[len(parts)-1]
	return realm, nil
}

func (m *MiddlewareManager) Recovery() gin.HandlerFunc {
	return httputil.RecoveryMiddleware(m.logger)
}
