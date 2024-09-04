package middleware

import (
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func AuthRequired(noAuthRoutes map[string]map[string]bool, adminRoutes map[string]map[string]string, jwtService *utils.JWT, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if paths, methodExists := noAuthRoutes[c.Request.Method]; methodExists {
			if skipAuth := paths[c.Request.URL.Path]; skipAuth {
				c.Next()
				return
			}
		}
		token := utils.GetHeaderToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := jwtService.ParseToken(token)
		if err != nil || claims == nil || claims.UserID == 0 || claims.Group == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized"})
			c.Abort()
			return
		}
		if claims.IsAdmin() {
			if paths, methodExists := adminRoutes[c.Request.Method]; methodExists {
				if _, needAdmin := paths[c.Request.URL.Path]; !needAdmin {
					logger.Error("unauthorized try admin", zap.String("path", c.Request.URL.Path))
					c.JSON(http.StatusForbidden, gin.H{"errorMessage": "forbidden"})
					c.Abort()
					return
				}
			}
		}
		session := sessions.Default(c)
		oldToken := session.Get("token")
		if oldToken == nil || oldToken.(string) != token {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized"})
			c.Abort()
			return
		}
		utils.SetClaimsToContext(c, claims)
		c.Next()
	}
}
