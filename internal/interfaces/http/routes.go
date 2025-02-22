package http

// func (s *Server) setupRoutes() {
// 	api := s.router.Group("api")
// 	{
// 		// s.router.GET("/ws/notifications", s.wsNotifyHandler.HandleWebSocket)

// 		apiTenant := api.Group(":tenantid")
// 		{
// 			apiBranch := apiTenant.Group(":branchid")
// 			{
// 				apiBranch.Use(httputil.RecoveryMiddleware(s.logger))
// 				apiBranch.Use(s.tenantHandler.TenantIDMiddleware(), s.AuthMiddleware())
// 				// notifications := apiBranch.Group("notifications")
// 				// {
// 				// 	notifications.GET("/ws-token", s.tokenHandler.GenerateToken)
// 				// }

// 			}

// 		}
// 	}

// }

// func (m *Server) AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if strings.Contains(c.Request.URL.String(), "tenant") {
// 			//TODO:remove only for testing
// 			c.Next()
// 			return
// 		}
// 		token, err := ginhelp.GetAuhtHeader(c)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		// Validate token
// 		jwtToken, err := m.auth.ValidateToken(c.Request.Context(), token)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error": "invalid token",
// 			})
// 			return
// 		}
// 		claims := jwtToken.Claims.(jwt.MapClaims)
// 		userID, ok := claims["sub"].(string)
// 		if !ok {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error": "invalid token claims",
// 			})
// 			return
// 		}

// 		ctx := httputil.SetToken(c.Request.Context(), token)
// 		ctx = httputil.SetUserID(ctx, userID)
// 		c.Request = c.Request.WithContext(ctx)
// 		c.Next()
// 	}
// }
