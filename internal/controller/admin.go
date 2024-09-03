package controller

import (
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
		adminGroup.GET("/token/refresh", a.RefreshJwt)
		adminGroup.POST("/login", a.Login)
		adminGroup.POST("/logout", a.Logout)
		adminGroup.POST("/password", a.UpdatePassword)
	}
}

type LoginReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AutoLogin bool   `json:"autoLogin"`
	Type      string `json:"type"`
}

type LoginRes struct {
	Status           string `json:"status"`
	Type             string `json:"type"`
	CurrentAuthority string `json:"currentAuthority"`
	Token            string `json:"token"`
}

func (a *AdminAPI) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, _, err := a.adminBiz.Auth(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	session := sessions.Default(c)
	session.Set("token", token)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.SetHeaderToken(c, token)
	var res LoginRes
	res.Token = token
	res.CurrentAuthority = "admin"
	res.Type = "account"
	res.Status = "ok"
	c.JSON(http.StatusOK, res)
}

func (a *AdminAPI) Current(c *gin.Context) {
	claims, ok := utils.MustGetClaimsFormContext(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, claims)
}
func (a *AdminAPI) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.JSON(http.StatusOK, gin.H{})
}

func (a *AdminAPI) RefreshJwt(c *gin.Context) {
	claims, ok := utils.MustGetClaimsFormContext(c)
	if !ok {
		return
	}
	jwt, err := a.adminBiz.RefreshJwt(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	session := sessions.Default(c)
	session.Set("token", jwt)
	err = session.Save()
	c.JSON(http.StatusOK, gin.H{"token": jwt})
}

type UpdatePasswordReq struct {
	Password string `json:"password"`
}

func (a *AdminAPI) UpdatePassword(c *gin.Context) {
	var req UpdatePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims, ok := utils.MustGetClaimsFormContext(c)
	if !ok {
		return
	}
	err = a.adminBiz.UpdatePassword(claims.UserID, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
	return
}

func (a *AdminAPI) GetAllAdmins(c *gin.Context) {
	admins, err := a.adminBiz.GetAllAdmins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, admins)
}
