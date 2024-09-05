package middleware

import (
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TokenRequired(jwtService *utils.JWT, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetHeaderToken(c)
		if token == "" {
			c.Next()
			return
		}
		claims, err := jwtService.ParseToken(token)
		if err != nil {
			logger.Error("token err", zap.Error(err))
			c.Next()
			return
		}
		utils.SetClaimsToContext(c, claims)
		c.Next()
	}
}
