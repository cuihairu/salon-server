package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StatisticsAPI struct {
	config *config.Config
	logger *zap.Logger
}

func (s *StatisticsAPI) Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	s.config = config
	s.logger = logger
	return nil
}

func (s *StatisticsAPI) RegisterRoutes(router *gin.RouterGroup) {
	statisticsGroup := router.Group("/statistics")
	{
		statisticsGroup.GET("/welcome", middleware.RequiredRole(middleware.Admin), s.GetStatistics)
	}
}

type StatisticsView struct {
	TodayIncome     int   `json:"todayIncome"`
	TodayVisitors   int   `json:"todayVisitors"`
	MonthIncome     int   `json:"monthIncome"`
	YearIncome      int   `json:"yearIncome"`
	DailyVisitors   []int `json:"dailyVisitors"`
	MonthlyVisitors []int `json:"monthlyVisitors"`
}

func (s *StatisticsAPI) GetStatistics(c *gin.Context) {
	ctx := utils.NewContext(c)
	ctx.Success(&StatisticsView{
		TodayIncome:     12000,
		TodayVisitors:   34,
		MonthIncome:     150000,
		YearIncome:      1000000,
		DailyVisitors:   []int{30, 45, 50, 60, 90, 100, 900, 1000},
		MonthlyVisitors: []int{300, 400, 500, 600, 700, 800, 9000, 10000},
	})
}
