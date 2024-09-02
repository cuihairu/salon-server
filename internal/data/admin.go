package data

import (
	"fmt"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type AdminRepository struct {
	utils.Repository[model.Admin]
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		Repository: utils.NewRepository[model.Admin](db),
	}
}

func (a *AdminRepository) FindByName(name string) (*model.Admin, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("name cannot be empty")
	}
	admins, err := a.FindByField("name", name)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, nil // NotFound
	}
	return &admins[0], nil
}

func (a *AdminRepository) FindByPhone(phone string) (*model.Admin, error) {
	if len(phone) == 0 {
		return nil, fmt.Errorf("phoneqq cannot be empty")
	}
	admins, err := a.FindByField("phone", phone)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, nil // NotFound
	}
	return &admins[0], nil
}
