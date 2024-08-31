package data

import (
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	utils.Repository[model.Service]
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{
		Repository: utils.NewRepository[model.Service](db),
	}
}
