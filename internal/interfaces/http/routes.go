package http

import (
	ginhelp "ddd/shared/http/gin"
	"ddd/shared/http/httputil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) setupRoutes() {
	api := s.router.Group("api")
	{
		s.router.GET("/ws/notifications", s.wsNotifyHandler.HandleWebSocket)
		tenant := api.Group("tenant")
		{
			tenant.Use(s.AuthMiddleware())
			tenant.POST("", s.tenantHandler.CreateTenant)
			tenant.GET("", s.tenantHandler.ListTenants)
			tenant.GET(":id", s.tenantHandler.GetTenant)
			tenant.PUT(":id", s.tenantHandler.UpdateTenant)
			tenant.DELETE(":id", s.tenantHandler.DeleteTenant)
		}

		apiTenant := api.Group(":tenantid")
		{
			apiBranch := apiTenant.Group(":branchid")
			{
				apiBranch.Use(httputil.RecoveryMiddleware(s.logger))
				apiBranch.Use(s.tenantHandler.TenantIDMiddleware(), s.AuthMiddleware())
				notifications := apiBranch.Group("notifications")
				{
					notifications.GET("/ws-token", s.tokenHandler.GenerateToken)
				}
				// Asset routes
				assets := apiBranch.Group("assets")
				{
					assets.POST("", s.assetHandler.CreateAsset)
					assets.GET("", s.assetHandler.ListAssets)
					assets.GET(":id", s.assetHandler.GetAsset)
					assets.PUT(":id", s.assetHandler.UpdateAsset)
					assets.DELETE(":id", s.assetHandler.DeleteAsset)
				}
				// User routes
				users := apiBranch.Group("users")
				{
					users.POST("", s.userHandler.CreateUser)
					users.GET("", s.userHandler.ListUsers)
					users.GET(":id", s.userHandler.GetUser)
					users.PUT(":id", s.userHandler.UpdateUser)
					users.DELETE(":id", s.userHandler.DeleteUser)
				}
				// Roles routes
				roles := apiBranch.Group("roles")
				{
					roles.POST("", s.rolesHandler.CreateRole)
					roles.GET("", s.rolesHandler.ListRoles)
					roles.GET(":name", s.rolesHandler.GetRole)
					roles.DELETE(":name", s.rolesHandler.DeleteRole)
				}
				// Resource routes
				resource := apiBranch.Group("resources")
				{
					resource.GET("", s.resourceHandler.ListResources)
				}

				// 	// Measure routes
				// 	measures := v1.Group("/measures")
				// 	{
				// 		measures.POST("", s.measureHandler.Create)
				// 		measures.GET("", s.measureHandler.List)
				// 		measures.GET("/:id", s.measureHandler.Get)
				// 	}

				// 	// Dashboard routes
				// 	dashboards := v1.Group("/dashboards")
				// 	{
				// 		dashboards.POST("", s.dashboardHandler.Create)
				// 		dashboards.GET("", s.dashboardHandler.List)
				// 		dashboards.GET("/:id", s.dashboardHandler.Get)
				// 		dashboards.PUT("/:id", s.dashboardHandler.Update)
				// 		dashboards.DELETE("/:id", s.dashboardHandler.Delete)
				// 	}

				// 	// Widget routes
				// 	widgets := v1.Group("/widgets")
				// 	{
				// 		widgets.POST("", s.widgetHandler.Create)
				// 		widgets.GET("", s.widgetHandler.List)
				// 		widgets.GET("/:id", s.widgetHandler.Get)
				// 		widgets.PUT("/:id", s.widgetHandler.Update)
				// 		widgets.DELETE("/:id", s.widgetHandler.Delete)
				// 	}
			}

		}
	}

}

func (m *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.String(), "tenant") {
			//TODO:remove only for testing
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

		// Validate token
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
