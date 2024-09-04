package controller

import (
	"fmt"
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

func (api *CategoryAPI) RegisterRoutes(router *gin.RouterGroup) {
	categoryGroup := router.Group("/category")
	{
		//categoryGroup.GET("/", api.GetAllCategories)
		//categoryGroup.GET("/:id", api.GetCategoryByID)
		categoryGroup.POST("/", api.CreateCategory)
		//categoryGroup.PUT("/:id", api.UpdateCategory)
		//categoryGroup.DELETE("/:id", api.DeleteCategory)
	}
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
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := api.categoryBiz.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

type CategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (api *CategoryAPI) CreateCategory(c *gin.Context) {
	var categoryParams CategoryParams
	ctx := utils.NewContext(c)
	if err := c.BindJSON(&categoryParams); err != nil {
		ctx.BadRequest(err)
		return
	}
	if len(categoryParams.Name) == 0 || len(categoryParams.Description) == 0 {
		ctx.BadRequest(fmt.Errorf("name or desc is nil"))
		return
	}
	category := model.Category{
		Name:        categoryParams.Name,
		Description: categoryParams.Description,
	}
	if err := api.categoryBiz.CreateCategory(&category); err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(category)
}

func (api *CategoryAPI) UpdateCategory(c *gin.Context) {
	var categoryParams CategoryParams
	ctx := utils.NewContext(c)
	if err := c.BindJSON(&categoryParams); err != nil {
		ctx.BadRequest(err)
		return
	}
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	if len(categoryParams.Name) == 0 && len(categoryParams.Description) == 0 {
		ctx.BadRequest(fmt.Errorf("name or desc is nil"))
		return
	}
	category := model.Category{
		Name:        categoryParams.Name,
		Description: categoryParams.Description,
	}
	if err := api.categoryBiz.UpdateCategory(id, &category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (api *CategoryAPI) DeleteCategory(c *gin.Context) {
	ctx := utils.NewContext(c)
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	if err := api.categoryBiz.DeleteCategory(id); err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.OK()
}
