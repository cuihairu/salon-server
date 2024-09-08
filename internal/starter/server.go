package starter

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	app    *App
}

func NewHttpServer(app *App) (*HttpServer, error) {
	serverConfig, err := app.Config.GetServerConfig()
	if err != nil {
		return nil, err
	}
	server := endless.NewServer(serverConfig.Address, app.Engine)
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 23
	return &HttpServer{
		app:    app,
		server: server,
		logger: app.Logger,
		config: app.Config,
	}, nil
}

func (s *HttpServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}
