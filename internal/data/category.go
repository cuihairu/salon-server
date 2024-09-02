package data

import (
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	utils.Repository[model.Category]
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		Repository: utils.NewRepository[model.Category](db),
	}
}
