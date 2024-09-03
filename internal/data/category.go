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

func (c *CategoryRepository) GetCategoryNameMap() (map[uint]string, error) {
	all, err := c.FindAll()
	if err != nil {
		return nil, err
	}
	var result = make(map[uint]string)
	for _, v := range all {
		result[v.ID] = v.Name
	}
	return result, nil
}
