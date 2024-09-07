package controller

import (
	"github.com/cuihairu/salon/internal/biz"
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
		serviceGroup.GET("/:id", middleware.RequiredRole(middleware.Admin), s.GetServicesByID)
		serviceGroup.POST("/", middleware.RequiredRole(middleware.Admin), s.CreateService)
		serviceGroup.PUT("/:id", middleware.RequiredRole(middleware.Admin), s.UpdateService)
		serviceGroup.DELETE("/:id", middleware.RequiredRole(middleware.Admin), s.DeleteService)

	}
}

type ServiceView struct {
	Name         string   `json:"name"`
	CategoryId   uint     `json:"category_id"`
	CategoryName string   `json:"category_name"`
	Intro        string   `json:"intro"`
	Cover        string   `json:"cover"`
	Images       []string `json:"images"`
	Description  string   `json:"description"`
	Duration     int      `json:"duration"`
	Price        float64  `json:"price"`
	Amount       float64  `json:"amount"`
	Recommend    bool     `json:"recommend"`
}

func (s *ServiceAPI) GetAllServices(c *gin.Context) {
	ctx := utils.NewContext(c)
	services, err := s.serviceBiz.GetAllServices()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(services)
}

func (s *ServiceAPI) GetServicesByID(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	services, err := s.serviceBiz.GetServiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
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
	ctx.Success(service)
}

func (s *ServiceAPI) UpdateService(c *gin.Context) {
	var service model.Service
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.serviceBiz.UpdateService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
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
