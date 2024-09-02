package middleware

import (
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthRequired(noAuthRoutes map[string]map[string]bool, jwtService *utils.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		if paths, methodExists := noAuthRoutes[c.Request.Method]; methodExists {
			if skipAuth := paths[c.Request.URL.Path]; skipAuth {
				c.Next()
				return
			}
		}
		token := utils.GetHeaderToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := jwtService.ParseToken(token)
		if err != nil || claims == nil || claims.UserID == 0 || claims.Group == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		session := sessions.Default(c)
		oldToken := session.Get("token")
		if oldToken == nil || oldToken.(string) != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		utils.SetClaimsToContext(c, claims)
		c.Next()
	}
}
