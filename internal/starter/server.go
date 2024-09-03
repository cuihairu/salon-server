package starter

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Server interface {
	ListenAndServe() error
}

type Controller interface {
	RegisterRoutes(router *gin.Engine)
}

type HttpServer struct {
	server Server
	logger *zap.Logger
	config *config.Config
}

func NewHttpServer(config *config.Config, db *gorm.DB, redisStore redis.Store, logger *zap.Logger) (*HttpServer, error) {
	serverConfig, err := config.GetServerConfig()
	if err != nil {
		return nil, err
	}
	router, err := NewRouter(config, db, redisStore, logger)
	if err != nil {
		return nil, err
	}
	server := endless.NewServer(serverConfig.Address, router)
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 23
	return &HttpServer{
		server: server,
		logger: logger,
		config: config,
	}, nil
}

func (s *HttpServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}
