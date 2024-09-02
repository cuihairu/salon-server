package controller

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AdminAPI struct {
	adminBiz *biz.AdminBiz
	logger   *zap.Logger
	config   *config.Config
}

func NewAdminAPI(config *config.Config, adminBiz *biz.AdminBiz, logger *zap.Logger) *AdminAPI {
	return &AdminAPI{
		adminBiz: adminBiz,
		logger:   logger,
		config:   config,
	}
}

func (a *AdminAPI) RegisterRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/jwt/refresh", a.RefreshJwt)
		adminGroup.POST("/login", a.Login)
		adminGroup.POST("/logout", a.Logout)
		adminGroup.POST("/password", a.UpdatePassword)
	}
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AdminAPI) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, admin, err := a.adminBiz.Auth(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	session := sessions.Default(c)
	key := fmt.Sprintf("admin:%d", admin.ID)
	session.Set(key, token)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header(utils.AuthorizationKey, utils.AuthorizationPrefix+token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *AdminAPI) Logout(c *gin.Context) {

}

func (a *AdminAPI) RefreshJwt(c *gin.Context) {

}

func (a *AdminAPI) UpdatePassword(c *gin.Context) {

}
