package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthAPI struct {
	userBiz *biz.UserBiz
	authBiz *biz.Auth
	logger  *zap.Logger
	config  *config.Config
}

func NewAuthAPI(config *config.Config, userBiz *biz.UserBiz, authBiz *biz.Auth, logger *zap.Logger) *AuthAPI {
	return &AuthAPI{
		userBiz: userBiz,
		authBiz: authBiz,
		logger:  logger,
		config:  config,
	}
}

func (api *AuthAPI) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/auth")
	{
		userGroup.POST("/login", api.Login)
		userGroup.POST("/logout", api.Logout)
	}
}

func (api *AuthAPI) Login(c *gin.Context) {

}

func (api *AuthAPI) Logout(c *gin.Context) {

}
