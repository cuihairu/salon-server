package controller

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StaticAPI struct {
	config *config.Config
	logger *zap.Logger
}

func (s *StaticAPI) Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	s.logger = logger
	s.config = config
	return nil
}

func (s *StaticAPI) RegisterRoutes(router *gin.RouterGroup) {
	staticGroup := router.Group("/static")
	{
		staticGroup.GET("/list", middleware.RequiredRole(middleware.Admin), s.GetStaticList)
		staticGroup.GET("/files", middleware.RequiredRole(middleware.Admin), s.ListFiles)
		staticGroup.POST("/upload", middleware.RequiredRole(middleware.Admin), s.Upload)
	}
}

func (s *StaticAPI) ListFiles(c *gin.Context) {
	ctx := utils.NewContext(c)
	staticConfig, err := s.config.GetStaticConfig()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	files, err := utils.ListFiles(staticConfig.StaticPath)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(files)
}

func (s *StaticAPI) GetStaticList(c *gin.Context) {
	ctx := utils.NewContext(c)
	staticConfig, err := s.config.GetStaticConfig()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	files, err := utils.ListFiles(staticConfig.StaticPath)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	fileUrls := make([]string, len(files))
	for i, file := range files {
		fileUrls[i] = fmt.Sprintf("%s%s/%s", staticConfig.Domain, "/static", file)
	}
	ctx.Success(fileUrls)
}

func (s *StaticAPI) Upload(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
