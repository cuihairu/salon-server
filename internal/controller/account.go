package controller

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AccountAPI struct {
	accountBiz *biz.AccountBiz
	logger     *zap.Logger
}

func (a *AccountAPI) Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	a.accountBiz = bizStore.AccountBiz
	a.logger = logger
	return nil
}

func (a *AccountAPI) RegisterRoutes(router *gin.RouterGroup) {
	accountGroup := router.Group("/account")
	{
		accountGroup.GET("/:id", middleware.RequiredRole(middleware.User), a.GetAccountInfo)
		accountGroup.GET("/", middleware.RequiredRole(middleware.User), a.GetAllAccounts)
		accountGroup.PUT("/:id", middleware.RequiredRole(middleware.User), a.UpdateAccount)
		accountGroup.DELETE("/:id", middleware.RequiredRole(middleware.User), a.DeleteAccount)
	}
}

func (a *AccountAPI) GetAccountInfo(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	ctx := utils.NewContext(c)
	account, err := a.accountBiz.GetAccountInfo(id)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	if account == nil {
		ctx.NotFound(fmt.Errorf("account not found"))
		return
	}
	ctx.Success(account)
}

func (a *AccountAPI) GetAllAccounts(c *gin.Context) {
	ctx := utils.NewContext(c)
	accounts, err := a.accountBiz.GetAllAccounts()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(accounts)
}

func (a *AccountAPI) UpdateAccount(c *gin.Context) {
	ctx := utils.NewContext(c)
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	var acc model.Account
	if err = c.ShouldBindJSON(&acc); err != nil {
		ctx.BadRequest(err)
		return
	}
	err = a.accountBiz.UpdateAccount(id, &acc)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(acc)
}

func (a *AccountAPI) DeleteAccount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
