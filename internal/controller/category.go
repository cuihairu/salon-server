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
)

type CategoryAPI struct {
	categoryBiz *biz.CategoryBiz
	logger      *zap.Logger
}

func (api *CategoryAPI) Initialize(conf *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	api.categoryBiz = bizStore.CategoryBiz
	api.logger = logger
	return nil
}

func (api *CategoryAPI) RegisterRoutes(router *gin.RouterGroup) {
	categoryGroup := router.Group("/category")
	{
		categoryGroup.GET("/", middleware.RequiredRole(middleware.User), api.GetAllCategories)
		categoryGroup.GET("/:id", middleware.RequiredRole(middleware.User), api.GetCategoryByID)
		categoryGroup.POST("/", middleware.RequiredRole(middleware.Admin), api.CreateCategory)
		categoryGroup.PUT("/:id", middleware.RequiredRole(middleware.Admin), api.UpdateCategory)
		categoryGroup.DELETE("/:id", middleware.RequiredRole(middleware.Admin), api.DeleteCategory)
	}
}

type CategoryView struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func (api *CategoryAPI) GetAllCategories(c *gin.Context) {
	ctx := utils.NewContext(c)
	categories, err := api.categoryBiz.GetAllCategories()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	categoryList := make([]CategoryView, 0)
	for _, category := range categories {
		categoryList = append(categoryList, CategoryView{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.Model.CreatedAt.UnixMilli(),
			UpdatedAt:   category.Model.UpdatedAt.UnixMilli(),
		})
	}
	ctx.Paginated(categoryList, int64(len(categoryList)))
}

func (api *CategoryAPI) GetCategoryByID(c *gin.Context) {
	ctx := utils.NewContext(c)
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	category, err := api.categoryBiz.GetCategoryByID(id)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	categoryView := CategoryView{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.Model.CreatedAt.UnixMilli(),
		UpdatedAt:   category.Model.UpdatedAt.UnixMilli(),
	}
	ctx.Success(categoryView)
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
	categoryView := CategoryView{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.Model.CreatedAt.UnixMilli(),
		UpdatedAt:   category.Model.UpdatedAt.UnixMilli(),
	}
	ctx.Success(categoryView)
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
		ctx.ServerError(err)
		return
	}
	categoryView := CategoryView{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.Model.CreatedAt.UnixMilli(),
		UpdatedAt:   category.Model.UpdatedAt.UnixMilli(),
	}
	ctx.Success(categoryView)
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
