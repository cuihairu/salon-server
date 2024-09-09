package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ServiceAPI struct {
	serviceBiz *biz.ServiceBiz
	logger     *zap.Logger
}

func (s *ServiceAPI) Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	s.serviceBiz = bizStore.ServiceBiz
	s.logger = logger
	return nil
}

func NewServiceAPI(serviceBiz *biz.ServiceBiz, logger *zap.Logger) *ServiceAPI {
	return &ServiceAPI{
		serviceBiz: serviceBiz,
		logger:     logger,
	}
}

func (s *ServiceAPI) RegisterRoutes(router *gin.RouterGroup) {
	serviceGroup := router.Group("/services")
	{
		serviceGroup.GET("/", middleware.RequiredRole(middleware.Admin), s.GetAllServices)
		serviceGroup.GET("/paging", middleware.RequiredRole(middleware.User), s.GetServicesByPaging)
		serviceGroup.GET("/:id", middleware.RequiredRole(middleware.Admin), s.GetServicesByID)
		serviceGroup.POST("/", middleware.RequiredRole(middleware.Admin), s.CreateService)
		serviceGroup.GET("/category/:id", middleware.RequiredRole(middleware.Admin), s.GetServicesByCategory)
		serviceGroup.PUT("/:id", middleware.RequiredRole(middleware.Admin), s.UpdateService)
		serviceGroup.DELETE("/:id", middleware.RequiredRole(middleware.Admin), s.DeleteService)

	}
}

type ServiceView struct {
	Id           uint    `json:"id"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
	Name         string  `json:"name"`
	CategoryId   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Intro        string  `json:"intro"`
	Cover        string  `json:"cover"`
	Content      string  `json:"content"`
	Description  string  `json:"description"`
	Duration     int     `json:"duration"`
	Price        float64 `json:"price"`
	Amount       float64 `json:"amount"`
	Recommend    bool    `json:"recommend"`
}

func serviceToView(service *model.Service) *ServiceView {
	return &ServiceView{
		Id:           service.ID,
		CreatedAt:    service.CreatedAt.UnixMilli(),
		UpdatedAt:    service.UpdatedAt.UnixMilli(),
		Name:         service.Name,
		CategoryId:   service.CategoryId,
		CategoryName: service.CategoryName,
		Intro:        service.Intro,
		Cover:        service.Cover,
		Content:      service.Content,
		Duration:     service.Duration,
		Price:        service.Price,
		Amount:       service.Amount,
		Recommend:    service.Recommend,
	}
}

func servicesToViews(services []model.Service) []*ServiceView {
	views := make([]*ServiceView, 0, len(services))
	for _, v := range services {
		views = append(views, serviceToView(&v))
	}
	return views
}

func (s *ServiceAPI) GetAllServices(c *gin.Context) {
	ctx := utils.NewContext(c)
	services, err := s.serviceBiz.GetAllServices()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(servicesToViews(services))
}

func (s *ServiceAPI) GetServicesByID(c *gin.Context) {
	ctx := utils.NewContext(c)
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	services, err := s.serviceBiz.GetServiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Success(serviceToView(services))
}

func (s *ServiceAPI) CreateService(c *gin.Context) {
	ctx := utils.NewContext(c)
	var service model.Service
	if err := c.BindJSON(&service); err != nil {
		ctx.BadRequest(err)
		return
	}
	if err := s.serviceBiz.CreateService(&service); err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(serviceToView(&service))
}

func (s *ServiceAPI) UpdateService(c *gin.Context) {
	ctx := utils.NewContext(c)
	var service model.Service
	if err := c.BindJSON(&service); err != nil {
		ctx.BadRequest(err)
		return
	}
	if err := s.serviceBiz.UpdateService(&service); err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(serviceToView(&service))
}

func (s *ServiceAPI) DeleteService(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	if err := s.serviceBiz.DeleteService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *ServiceAPI) GetServicesByCategory(context *gin.Context) {
	id, err := utils.ParseUintParam[uint](context, "id")
	if err != nil {
		return
	}
	services, err := s.serviceBiz.GetServicesByCategory(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, servicesToViews(services))
}

func (s *ServiceAPI) GetServicesByPaging(context *gin.Context) {
	ctx := utils.NewContext(context)
	services, total, err := s.serviceBiz.GetServicesByPaging(ctx.Paging())
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Paginated(servicesToViews(services), total)
}
