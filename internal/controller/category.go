package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CategoryAPI struct {
	categoryBiz *biz.CategoryBiz
	logger      *zap.Logger
}

func NewCategoryAPI(categoryBiz *biz.CategoryBiz, logger *zap.Logger) *CategoryAPI {
	return &CategoryAPI{
		categoryBiz: categoryBiz,
		logger:      logger,
	}
}

func (api *CategoryAPI) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/categories", api.GetAllCategories)
	group.GET("/categories/:id", api.GetCategoryByID)
	group.POST("/categories", api.CreateCategory)
	group.PUT("/categories/:id", api.UpdateCategory)
	group.DELETE("/categories/:id", api.DeleteCategory)
}

func (api *CategoryAPI) GetAllCategories(c *gin.Context) {
	categories, err := api.categoryBiz.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (api *CategoryAPI) GetCategoryByID(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := api.categoryBiz.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (api *CategoryAPI) CreateCategory(c *gin.Context) {

	var category model.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := api.categoryBiz.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (api *CategoryAPI) UpdateCategory(c *gin.Context) {

	var category model.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := api.categoryBiz.UpdateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (api *CategoryAPI) DeleteCategory(c *gin.Context) {

	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := api.categoryBiz.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
