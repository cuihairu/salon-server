package data

import (
	"fmt"
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

func (c *CategoryRepository) Create(category *model.Category) error {
	return c.ExecuteInTransaction(func(tx *gorm.DB) error {
		var existing model.Category
		// 查找时忽略软删除标志
		if err := tx.Unscoped().Where("name = ?", category.Name).First(&existing).Error; err == nil {
			if existing.ID != 0 && existing.DeletedAt.Valid {
				// 恢复已软删除的记录
				existing.DeletedAt = gorm.DeletedAt{}
				existing.Description = category.Description
				return tx.Save(&existing).Error
			}
			return fmt.Errorf("record with name %s already exists", category.Name)
		}
		return tx.Create(&category).Error
	})
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
