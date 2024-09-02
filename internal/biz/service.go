package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type ServiceBiz struct {
	serviceRepo *data.ServiceRepository
	logger      *zap.Logger
}

func NewServiceBiz(serviceRepo *data.ServiceRepository, logger *zap.Logger) *ServiceBiz {
	return &ServiceBiz{
		serviceRepo: serviceRepo,
		logger:      logger,
	}
}

func (biz *ServiceBiz) GetAllServices() ([]model.Service, error) {
	return biz.serviceRepo.FindAll()
}

func (biz *ServiceBiz) GetServiceByID(id uint) (*model.Service, error) {
	return biz.serviceRepo.FindByID(id)
}

func (biz *ServiceBiz) CreateService(service *model.Service) error {
	return biz.serviceRepo.Create(service)
}
