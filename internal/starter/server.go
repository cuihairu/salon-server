package starter

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/controller"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/fvbock/endless"
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

func NewApiRouter(config *config.Config, router *gin.Engine, db *gorm.DB, logger *zap.Logger) (*gin.RouterGroup, error) {
	apiGroup := router.Group("/api")
	//data
	data, err := data.NewData(db, config, logger)
	if err != nil {
		return nil, err
	}
	// services
	jwtConfig := config.GetJwtConfig()
	jwtService := utils.NewJWT(jwtConfig.SecretKey, jwtConfig.Expire)
	// users
	userBiz := biz.NewUserBiz(data.UserRepo, logger)
	userApi := controller.NewUserAPI(userBiz, logger)
	userApi.RegisterRoutes(apiGroup)
	// auth
	authBiz := biz.NewAuth(config, jwtService, data.UserRepo, logger)
	authApi := controller.NewAuthAPI(config, userBiz, authBiz, logger)
	authApi.RegisterRoutes(apiGroup)
	return apiGroup, nil
}

func NewAdminRouter(config *config.Config, router *gin.Engine, db *gorm.DB, logger *zap.Logger) (*gin.RouterGroup, error) {
	apiGroup := router.Group("/admin")
	userRepo := data.NewUserRepository(db)
	userBiz := biz.NewUserBiz(userRepo, logger)
	userApi := controller.NewUserAPI(userBiz, logger)
	userApi.RegisterRoutes(apiGroup)
	return apiGroup, nil
}
func NewRouter(config *config.Config, db *gorm.DB, logger *zap.Logger) (*gin.Engine, error) {
	if config.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else if config.IsTest() {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// api
	router := gin.Default()
	_, err := NewApiRouter(config, router, db, logger)
	if err != nil {
		return nil, err
	}
	// admin
	_, err = NewAdminRouter(config, router, db, logger)
	if err != nil {
		return nil, err
	}
	return router, nil
}

func NewHttpServer(config *config.Config, db *gorm.DB, logger *zap.Logger) (*HttpServer, error) {
	serverConfig, err := config.GetServerConfig()
	if err != nil {
		return nil, err
	}
	router, err := NewRouter(config, db, logger)
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
