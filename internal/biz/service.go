package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"go.uber.org/zap"
)

type ServiceBiz struct {
	categoryRepo *data.CategoryRepository
	serviceRepo  *data.ServiceRepository
	logger       *zap.Logger
}

func NewServiceBiz(dataStore *data.DataStore, logger *zap.Logger) *ServiceBiz {
	return &ServiceBiz{
		serviceRepo:  dataStore.ServiceRepo,
		categoryRepo: dataStore.CategoryRepo,
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
	var ret = make([]model.Service, 0)
	for _, v := range all {
		v.CategoryName = categoryNameMap[v.CategoryId]
		ret = append(ret, v)
	}
	return all, nil
}

func (biz *ServiceBiz) GetServicesByCategory(id uint) ([]model.Service, error) {
	all, err := biz.serviceRepo.FindByField("category_id", id)
	if err != nil {
		return nil, err
	}
	categoryNameMap, err := biz.categoryRepo.GetCategoryNameMap()
	if err != nil {
		return nil, err
	}

	var ret = make([]model.Service, 0)
	for _, v := range all {
		v.CategoryName = categoryNameMap[v.CategoryId]
		ret = append(ret, v)
	}
	return ret, nil
}

func (biz *ServiceBiz) GetServicesByPaging(paging *utils.Paging) ([]model.Service, int64, error) {
	all, err := biz.serviceRepo.FindWithPaging(paging)
	if err != nil {
		return nil, 0, err
	}
	if all == nil {
		return make([]model.Service, 0), 0, nil
	}
	categoryNameMap, err := biz.categoryRepo.GetCategoryNameMap()
	if err != nil {
		return nil, 0, err
	}
	var ret = make([]model.Service, 0)
	for _, v := range all.List {
		v.CategoryName = categoryNameMap[v.CategoryId]
		ret = append(ret, v)
	}
	return ret, all.Total, nil
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
