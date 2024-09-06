package controller

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
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
		adminGroup.GET("/token/refresh", middleware.RequiredRole(middleware.Admin), a.RefreshJwt)
		adminGroup.POST("/login", middleware.RequiredRole(middleware.Anonymous), a.Login)
		adminGroup.POST("/logout", middleware.RequiredRole(middleware.Admin), a.Logout)
		adminGroup.POST("/password", middleware.RequiredRole(middleware.Admin), a.UpdatePassword)
		adminGroup.GET("/current", middleware.RequiredRole(middleware.Admin), a.Current)
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
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	token, _, err := a.adminBiz.Auth(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errorMessage": err.Error()})
		return
	}
	session := sessions.Default(c)
	session.Set("token", token)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
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

type Tag struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}
type Province struct {
	Label string `json:"label"`
	Key   string `json:"key"`
}
type City struct {
	Label string `json:"label"`
	Key   string `json:"key"`
}
type Geographic struct {
	Province Province `json:"province"`
	City     City     `json:"city"`
}
type CurrentRes struct {
	Name        string     `json:"name"`
	Avatar      string     `json:"avatar"`
	Email       string     `json:"email"`
	UserId      string     `json:"userId"`
	Signature   string     `json:"signature"`
	Title       string     `json:"title"`
	Group       string     `json:"group"`
	Tags        []Tag      `json:"tags"`
	NotifyCount int        `json:"notifyCount"`
	UnreadCount int        `json:"unreadCount"`
	Country     string     `json:"country"`
	Access      string     `json:"access"`
	Geographic  Geographic `json:"geographic"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone"`
}

func (a *AdminAPI) Current(c *gin.Context) {
	ctx := utils.NewContext(c)
	claims, ok := ctx.Claims()
	if !ok {
		return
	}
	admin, err := a.adminBiz.GetAdmin(claims.UserID)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	if admin == nil {
		ctx.NotFound(fmt.Errorf("admin not found"))
		return
	}
	var res CurrentRes
	res.Name = admin.Name
	res.UserId = fmt.Sprintf("%d", claims.UserID)
	res.Avatar = admin.Avatar
	if admin.Country != nil {
		res.Country = *admin.Country
	}
	res.Signature = admin.Signature
	res.Title = admin.Title
	res.Access = admin.Role
	res.Tags = make([]Tag, 0)
	if admin.Tags != nil {
		for _, tag := range *admin.Tags {
			res.Tags = append(res.Tags, Tag{Key: tag.Key, Label: tag.Label})
		}
	}
	res.Group = admin.Group
	if admin.Phone != nil {
		res.Phone = *admin.Phone
	} else {
		res.Phone = "13850000000"
	}
	res.Address = admin.Address
	res.Geographic = Geographic{
		Province: Province{
			Key:   "浙江",
			Label: "浙江",
		},
		City: City{
			Key:   "杭州",
			Label: "杭州",
		},
	}
	if admin.Geographic != nil {
		res.Geographic.Province.Key = admin.Geographic.Province.Key
		res.Geographic.City.Key = admin.Geographic.City.Key
		res.Geographic.Province.Label = admin.Geographic.Province.Label
		res.Geographic.City.Label = admin.Geographic.City.Label
	}
	res.Email = "admin@example.com"
	if admin.Email != nil {
		res.Email = *admin.Email
	}
	res.NotifyCount = 10
	res.UnreadCount = 2
	ctx.Success(res)
}
func (a *AdminAPI) Logout(c *gin.Context) {
	ctx := utils.NewContext(c)
	ctx.Session().Clear()
	ctx.Success(nil)
	ctx.OK()
}

func (a *AdminAPI) RefreshJwt(c *gin.Context) {
	ctx := utils.NewContext(c)
	claims, ok := ctx.Claims()
	if !ok {
		return
	}
	jwt, err := a.adminBiz.RefreshJwt(claims.UserID)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	session := sessions.Default(c)
	session.Set("token", jwt)
	err = session.Save()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(gin.H{"token": jwt})
}

type UpdatePasswordReq struct {
	Password string `json:"password"`
}

func (a *AdminAPI) UpdatePassword(c *gin.Context) {
	ctx := utils.NewContext(c)
	var req UpdatePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	claims, ok := ctx.Claims()
	if !ok {
		return
	}
	err = a.adminBiz.UpdatePassword(claims.UserID, req.Password)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.OK()
	return
}

func (a *AdminAPI) GetAllAdmins(c *gin.Context) {
	ctx := utils.NewContext(c)
	admins, err := a.adminBiz.GetAllAdmins()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(admins)
}
