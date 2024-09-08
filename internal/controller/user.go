package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserAPI struct {
	userBiz *biz.UserBiz
	logger  *zap.Logger
}

func (api *UserAPI) Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	api.userBiz = bizStore.UserBiz
	api.logger = logger
	return nil
}

func (api *UserAPI) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", middleware.RequiredRole(middleware.Admin), api.CreateUser)
		userGroup.GET("/:id", middleware.RequiredRole(middleware.Admin), api.GetUserByID)
		userGroup.GET("/", middleware.RequiredRole(middleware.Admin), api.GetAllUsers)
		userGroup.PUT("/:id", middleware.RequiredRole(middleware.Admin), api.UpdateUser)
		userGroup.DELETE("/:id", middleware.RequiredRole(middleware.Admin), api.DeleteUser)
	}
}

// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.User true "User"
// @Success 201 {object} model.User
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/ [post]
func (api *UserAPI) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//if err := api.userBiz.CreateUser(&user); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(http.StatusCreated, user)
}

func (api *UserAPI) GetUserByID(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	user, err := api.userBiz.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (api *UserAPI) GetAllUsers(c *gin.Context) {
	users, err := api.userBiz.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (api *UserAPI) UpdateUser(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.userBiz.UpdateUser(id, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (api *UserAPI) DeleteUser(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	if err := api.userBiz.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
