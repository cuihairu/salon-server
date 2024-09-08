package biz

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/utils"
	"go.uber.org/zap"
)

type BizStore struct {
	logger          *zap.Logger
	data            *data.DataStore
	AdminBiz        *AdminBiz
	OperationLogBiz *OperationLogBiz
	AccountBiz      *AccountBiz
	UserBiz         *UserBiz
	AuthBiz         *AuthBiz
	CategoryBiz     *CategoryBiz
	ServiceBiz      *ServiceBiz
	OrderBiz        *OrderBiz
	MemberBiz       *MemberBiz
}

func NewBizStore(config *config.Config, dataStore *data.DataStore, tokenService *utils.JWT, logger *zap.Logger) *BizStore {
	return &BizStore{
		logger:          logger,
		data:            dataStore,
		AdminBiz:        NewAdminBiz(config, tokenService, dataStore, logger),
		OperationLogBiz: NewOperationLogBiz(dataStore, logger),
		AccountBiz:      NewAccountBiz(dataStore, logger),
		UserBiz:         NewUserBiz(dataStore, logger),
		AuthBiz:         NewAuthBiz(config, tokenService, dataStore, logger),
		CategoryBiz:     NewCategoryBiz(dataStore, logger),
		ServiceBiz:      NewServiceBiz(dataStore, logger),
		OrderBiz:        NewOrderBiz(dataStore, logger),
		MemberBiz:       NewMemberBiz(dataStore, logger),
	}
}
