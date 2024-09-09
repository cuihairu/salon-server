package controller

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminAPI struct {
	adminBiz        *biz.AdminBiz
	operationLogBiz *biz.OperationLogBiz
	logger          *zap.Logger
	config          *config.Config
}

func (a *AdminAPI) Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	a.adminBiz = bizStore.AdminBiz
	a.operationLogBiz = bizStore.OperationLogBiz
	a.logger = logger
	a.config = config
	return nil
}

func (a *AdminAPI) RegisterRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/token/refresh", middleware.RequiredRole(middleware.Admin), a.RefreshToken)
		adminGroup.POST("/login", a.Login)
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
	ctx := utils.NewContext(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error("invalid request", zap.String("path", c.Request.URL.Path), zap.Error(err))
		ctx.BadRequest(err)
		return
	}
	token, _, err := a.adminBiz.Auth(req.Username, req.Password)
	if err != nil {
		a.logger.Error("login failed", zap.Error(err))
		ctx.BadRequest(err)
		a.operationLogBiz.Log(req.Username, "admin", ctx.ClientIP(), "", c.Request.UserAgent(), "admin", "login", "", err.Error())
		return
	}
	if err = ctx.SetToken(token); err != nil {
		a.logger.Error("save session failed", zap.Error(err))
		return
	}
	var res LoginRes
	res.Token = token
	res.CurrentAuthority = "admin"
	res.Type = "account"
	res.Status = "ok"
	ctx.Success(res)
	a.operationLogBiz.Log(req.Username, res.CurrentAuthority, ctx.ClientIP(), "", c.Request.UserAgent(), "admin", "login", "", "")
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
	if admin.Tags.HasValue() {
		tags := admin.Tags.Data()
		for _, tag := range *tags {
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
	if admin.Geographic.HasValue() {
		geographic := admin.Geographic.Data()
		res.Geographic.Province.Key = geographic.Province.Key
		res.Geographic.City.Key = geographic.City.Key
		res.Geographic.Province.Label = geographic.Province.Label
		res.Geographic.City.Label = geographic.City.Label
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
	ctx.OK()
}

type RefreshTokenRes struct {
	Token string `json:"token"`
}

func (a *AdminAPI) RefreshToken(c *gin.Context) {
	ctx := utils.NewContext(c)
	claims, ok := ctx.Claims()
	if !ok {
		return
	}
	jwt, err := a.adminBiz.RefreshToken(claims.UserID)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	if err = ctx.SetToken(jwt); err != nil {
		ctx.ServerError(err)
		return
	}
	res := RefreshTokenRes{
		Token: jwt,
	}
	ctx.Success(res)
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
