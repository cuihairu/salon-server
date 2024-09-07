package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type ServiceBiz struct {
	categoryRepo *data.CategoryRepository
	serviceRepo  *data.ServiceRepository
	logger       *zap.Logger
}

func NewServiceBiz(serviceRepo *data.ServiceRepository, categoryRepo *data.CategoryRepository, logger *zap.Logger) *ServiceBiz {
	return &ServiceBiz{
		serviceRepo:  serviceRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (biz *ServiceBiz) GetAllServices() ([]model.Service, error) {
	all, err := biz.serviceRepo.FindAll()
	if err != nil {
		return nil, err
	}
	categoryNameMap, err := biz.categoryRepo.GetCategoryNameMap()
	if err != nil {
		return nil, err
	}

	for _, v := range all {
		v.CategoryName = categoryNameMap[v.CategoryId]
	}
	return all, nil
}

func (biz *ServiceBiz) GetServiceByID(id uint) (*model.Service, error) {
	return biz.serviceRepo.FindByID(id)
}

func (biz *ServiceBiz) CreateService(service *model.Service) error {
	return biz.serviceRepo.Create(service)
}

func (biz *ServiceBiz) UpdateService(service *model.Service) error {
	return biz.serviceRepo.Update(service)
}

func (biz *ServiceBiz) DeleteService(id uint) error {
	return biz.serviceRepo.Delete(id)
}
