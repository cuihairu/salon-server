package biz

import (
	"fmt"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type CategoryBiz struct {
	categoryRepo *data.CategoryRepository
	logger       *zap.Logger
}

func NewCategoryBiz(categoryRepo *data.CategoryRepository, logger *zap.Logger) *CategoryBiz {
	return &CategoryBiz{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (c *CategoryBiz) GetAllCategories() ([]model.Category, error) {
	return c.categoryRepo.FindAll()
}

func (c *CategoryBiz) GetCategoryByID(id uint) (*model.Category, error) {
	return c.categoryRepo.FindByID(id)
}

func (c *CategoryBiz) CreateCategory(category *model.Category) error {
	return c.categoryRepo.Create(category)
}

func (c *CategoryBiz) UpdateCategory(id uint, category *model.Category) error {
	if category == nil {
		return fmt.Errorf("category is nil")
	}
	preCategory, err := c.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if preCategory == nil {
		return fmt.Errorf("category not found")
	}
	category.Model = preCategory.Model
	return c.categoryRepo.Update(category)
}

func (c *CategoryBiz) DeleteCategory(id uint) error {
	return c.categoryRepo.Delete(id)
}
