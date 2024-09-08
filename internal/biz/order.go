package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type OrderBiz struct {
	orderRepo *data.OrderRepository
	logger    *zap.Logger
}

func NewOrderBiz(dataStore *data.DataStore, logger *zap.Logger) *OrderBiz {
	return &OrderBiz{
		orderRepo: dataStore.OrderRepo,
		logger:    logger,
	}
}

func (b *OrderBiz) GetAllOrders() ([]model.Order, error) {
	return b.orderRepo.FindAll()
}

func (b *OrderBiz) GetOrderByID(id uint) (*model.Order, error) {
	return b.orderRepo.FindByID(id)
}

func (b *OrderBiz) CreateOrder(order *model.Order) error {
	return b.orderRepo.Create(order)
}

func (b *OrderBiz) UpdateOrder(order *model.Order) error {
	return b.orderRepo.Update(order)
}

func (b *OrderBiz) DeleteOrder(id uint) error {
	return b.orderRepo.Delete(id)
}
