package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserAPI struct {
	userBiz *biz.UserBiz
}

func NewUserAPI(userBiz *biz.UserBiz) *UserAPI {
	return &UserAPI{
		userBiz: userBiz,
	}
}

func (api *UserAPI) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", api.CreateUser)
		userGroup.GET("/:id", api.GetUserByID)
		userGroup.GET("/", api.GetAllUsers)
		userGroup.PUT("/:id", api.UpdateUser)
		userGroup.DELETE("/:id", api.DeleteUser)
	}
}

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
