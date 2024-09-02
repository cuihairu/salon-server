package data

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Data struct {
	db           *gorm.DB
	UserRepo     *UserRepository
	AccountRepo  *AccountRepository
	AdminRepo    *AdminRepository
	MemberRepo   *MemberRepository
	OrderRepo    *OrderRepository
	ServiceRepo  *ServiceRepository
	CategoryRepo *CategoryRepository
	config       *config.Config
	logger       *zap.Logger
}

func (d *Data) AutoMigrate() error {
	return d.db.AutoMigrate(&model.User{}, &model.Account{}, &model.Member{}, &model.Order{}, &model.Category{}, &model.Service{}, &model.Admin{})
}

func NewData(db *gorm.DB, conf *config.Config, logger *zap.Logger) (*Data, error) {
	data := &Data{
		db:           db,
		UserRepo:     NewUserRepository(db),
		AccountRepo:  NewAccountRepository(db),
		AdminRepo:    NewAdminRepository(db),
		MemberRepo:   NewMemberRepository(db),
		OrderRepo:    NewOrderRepository(db),
		ServiceRepo:  NewServiceRepository(db),
		CategoryRepo: NewCategoryRepository(db),
		config:       conf,
		logger:       logger,
	}
	dbConfig, err := conf.GetDbConfig()
	if err != nil {
		return nil, err
	}
	if dbConfig.AutoMigrate {
		err = data.AutoMigrate()
	}
	return data, err
}
