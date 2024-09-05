package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthAPI struct {
	userBiz *biz.UserBiz
	authBiz *biz.AuthBiz
	logger  *zap.Logger
	config  *config.Config
}

func NewAuthAPI(config *config.Config, userBiz *biz.UserBiz, authBiz *biz.AuthBiz, logger *zap.Logger) *AuthAPI {
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
		userGroup.POST("/login", middleware.RequiredRole(middleware.Anonymous), api.Login)
		userGroup.POST("/logout", middleware.RequiredRole(middleware.Admin), api.Logout)
	}
}

func (api *AuthAPI) Login(c *gin.Context) {
	ctx := utils.NewContext(c)
	//api.logger.Info("login", zap.String("type", loginType))
	ctx.OK()
}

func (api *AuthAPI) Logout(c *gin.Context) {
	ctx := utils.NewContext(c)
	ctx.Session().Clear()
	ctx.OK()
}
