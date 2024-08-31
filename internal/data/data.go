package data

import (
	"github.com/cuihairu/salon/internal/model"
	"gorm.io/gorm"
)

type Data struct {
	db       *gorm.DB
	UserRepo *UserRepository
}

func (d *Data) AutoMigrate() error {
	return d.db.AutoMigrate(&model.User{}, &model.Account{}, &model.Member{}, &model.Order{}, &model.Service{}, &model.Admin{})
}

func NewData(db *gorm.DB) (*Data, error) {
	data := &Data{
		db:       db,
		UserRepo: NewUserRepository(db),
	}
	return data, nil
}
