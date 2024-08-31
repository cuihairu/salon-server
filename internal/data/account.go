package data

import (
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type AccountRepository struct {
	utils.Repository[model.Account]
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		Repository: utils.NewRepository[model.Account](db),
	}
}
