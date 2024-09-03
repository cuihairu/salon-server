package controller

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AccountAPI struct {
	accountBiz *biz.AccountBiz
	logger     *zap.Logger
}

func NewAccountAPI(accountBiz *biz.AccountBiz, logger *zap.Logger) *AccountAPI {
	return &AccountAPI{
		accountBiz: accountBiz,
		logger:     logger,
	}
}

func (a *AccountAPI) RegisterRoutes(router *gin.RouterGroup) {
	accountGroup := router.Group("/account")
	{
		accountGroup.GET("/:id", a.GetAccountInfo)
		accountGroup.GET("/", a.GetAllAccounts)
		accountGroup.PUT("/:id", a.UpdateAccount)
		accountGroup.DELETE("/:id", a.DeleteAccount)
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
	accounts, err := a.accountBiz.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (a *AccountAPI) UpdateAccount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (a *AccountAPI) DeleteAccount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
