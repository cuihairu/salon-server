package starter

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/controller"
	"github.com/cuihairu/salon/internal/data"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type ServerConfig struct {
	Address string `mapstructure:"address" yaml:"address" json:"address"`
}

type Server interface {
	ListenAndServe() error
}

type Controller interface {
	RegisterRoutes(router *gin.Engine)
}

type HttpServer struct {
	server   Server
	logger   *zap.Logger
	serveMux *http.ServeMux
	config   *Config
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func NewHttpServer(config *Config, db *gorm.DB, logger *zap.Logger) (*HttpServer, error) {
	serverConfig, err := config.GetServerConfig()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	if config.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else if config.IsTest() {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	apiRouter := gin.New()
	userRepo := data.NewUserRepository(db)
	userBiz := biz.NewUserBiz(userRepo)
	userApi := controller.NewUserAPI(userBiz)
	userApi.RegisterRoutes(apiRouter)
	mux.HandleFunc("/status", helloHandler)
	server := endless.NewServer(serverConfig.Address, mux)
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 23
	return &HttpServer{
		server:   server,
		serveMux: mux,
		logger:   logger,
		config:   config,
	}, nil
}

func (s *HttpServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}
