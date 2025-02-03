package http

import (
	ginhelp "backend/shared/http/gin"
	"backend/shared/http/httputil"
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
				// Measure routes
				measures := apiBranch.Group("measures")
				{
					measures.POST("", s.measureHandler.CreateMeasure)
					measures.GET("", s.measureHandler.ListMeasures)
					measures.GET(":id", s.measureHandler.GetMeasure)
					measures.PUT(":id", s.measureHandler.UpdateMeasure)
					measures.DELETE(":id", s.measureHandler.DeleteMeasure)
				}
				// Dashboard routes
				dashboards := apiBranch.Group("dashboards")
				{
					dashboards.POST("", s.dashboardHandler.CreateDashboard)
					dashboards.GET("", s.dashboardHandler.ListDashboards)
					dashboards.GET(":id", s.dashboardHandler.GetDashboard)
					dashboards.PUT(":id", s.dashboardHandler.UpdateDashboard)
					dashboards.DELETE(":id", s.dashboardHandler.DeleteDashboard)
					widgets := dashboards.Group(":id/widgets")
					{
						widgets.POST("", s.widgetHandler.CreateWidget)
						widgets.GET("", s.widgetHandler.ListWidgets)
						widgets.GET(":wid", s.widgetHandler.GetWidget)
						widgets.PUT(":wid", s.widgetHandler.UpdateWidget)
						widgets.DELETE(":wid", s.widgetHandler.DeleteWidget)
					}
				}
				// Device routes
				devices := apiBranch.Group("devices")
				{
					devices.POST("", s.deviceHandler.CreateDevice)
					devices.GET("", s.deviceHandler.ListDevices)
					devices.GET(":id", s.deviceHandler.GetDevice)
					devices.PUT(":id", s.deviceHandler.UpdateDevice)
					devices.DELETE(":id", s.deviceHandler.DeleteDevice)
				}
				// Action routes
				actions := apiBranch.Group("actions")
				{
					actions.POST("", s.actionHandler.CreateAction)
					actions.GET("", s.actionHandler.ListActions)
					actions.GET(":id", s.actionHandler.GetAction)
					actions.PUT(":id", s.actionHandler.UpdateAction)
					actions.DELETE(":id", s.actionHandler.DeleteAction)
				}
				// Alert routes
				alerts := apiBranch.Group("alerts")
				{
					alerts.POST("", s.alertHandler.CreateAlert)
					alerts.GET("", s.alertHandler.ListAlerts)
					alerts.GET(":id", s.alertHandler.GetAlert)
					alerts.PUT(":id", s.alertHandler.UpdateAlert)
					alerts.DELETE(":id", s.alertHandler.DeleteAlert)
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
					resource.PUT(":id", s.resourceHandler.UpdateResource)
				}
				scopes := apiBranch.Group("scopes")
				{
					scopes.GET("", s.scopesHandler.ListRoles)
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
