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

type OrderAPI struct {
	orderBiz *biz.OrderBiz
	logger   *zap.Logger
}

func NewOrderAPI(orderBiz *biz.OrderBiz, logger *zap.Logger) *OrderAPI {
	return &OrderAPI{
		orderBiz: orderBiz,
		logger:   logger,
	}
}

func (o *OrderAPI) RegisterRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.GET("/", middleware.RequiredRole(middleware.User), o.GetAllOrders)
		orderGroup.GET("/:id", middleware.RequiredRole(middleware.User), o.GetOrderByID)
		orderGroup.POST("/", middleware.RequiredRole(middleware.User), o.CreateOrder)
		orderGroup.PUT("/:id", middleware.RequiredRole(middleware.User), o.UpdateOrder)
		orderGroup.DELETE("/:id", middleware.RequiredRole(middleware.User), o.DeleteOrder)
	}
}

func (o *OrderAPI) GetAllOrders(c *gin.Context) {
	orders, err := o.orderBiz.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (o *OrderAPI) GetOrderByID(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	order, err := o.orderBiz.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (o *OrderAPI) CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := o.orderBiz.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (o *OrderAPI) UpdateOrder(c *gin.Context) {
	var order model.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := o.orderBiz.UpdateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (o *OrderAPI) DeleteOrder(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	if err := o.orderBiz.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
