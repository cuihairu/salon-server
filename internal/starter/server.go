package starter

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/controller"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/sessions"
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

func NewApiRouter(config *config.Config, router *gin.Engine, db *gorm.DB, jwtService *utils.JWT, logger *zap.Logger) (*gin.RouterGroup, error) {
	apiGroup := router.Group("/api")
	//data
	data, err := data.NewData(db, config, logger)
	if err != nil {
		return nil, err
	}
	// users
	userBiz := biz.NewUserBiz(data.UserRepo, logger)
	userApi := controller.NewUserAPI(userBiz, logger)
	userApi.RegisterRoutes(apiGroup)
	// auth
	authBiz := biz.NewAuthBiz(config, jwtService, data.UserRepo, logger)
	authApi := controller.NewAuthAPI(config, userBiz, authBiz, logger)
	authApi.RegisterRoutes(apiGroup)
	// category
	categoryBiz := biz.NewCategoryBiz(data.CategoryRepo, logger)
	categoryApi := controller.NewCategoryAPI(categoryBiz, logger)
	categoryApi.RegisterRoutes(apiGroup)
	// service
	serviceBiz := biz.NewServiceBiz(data.ServiceRepo, logger)
	serviceApi := controller.NewServiceAPI(serviceBiz, logger)
	serviceApi.RegisterRoutes(apiGroup)
	// account
	accountBiz := biz.NewAccountBiz(data.AccountRepo, logger)
	accountApi := controller.NewAccountAPI(accountBiz, logger)
	accountApi.RegisterRoutes(apiGroup)
	// order
	orderBiz := biz.NewOrderBiz(data.OrderRepo, logger)
	orderApi := controller.NewOrderAPI(orderBiz, logger)
	orderApi.RegisterRoutes(apiGroup)
	// member
	memberBiz := biz.NewMemberBiz(data.MemberRepo, logger)
	memberApi := controller.NewMemberAPI(memberBiz, logger)
	memberApi.RegisterRoutes(apiGroup)
	// admin
	adminBiz := biz.NewAdminBiz(config, jwtService, data.AdminRepo, logger)
	adminApi := controller.NewAdminAPI(config, adminBiz, logger)
	adminApi.RegisterRoutes(apiGroup)
	return apiGroup, nil
}

var noAuthRoutes = map[string]map[string]bool{
	"GET": {
		"/about": true,
	},
	"POST": {
		"/api/auth/login":  true,
		"/api/admin/login": true,
	},
}

func NewRouter(config *config.Config, db *gorm.DB, redisStore redis.Store, logger *zap.Logger) (*gin.Engine, error) {
	// services
	jwtConfig := config.GetJwtConfig()
	jwtService := utils.NewJWT(jwtConfig.SecretKey, jwtConfig.Expire)
	// api
	router := gin.Default()
	router.Use(sessions.Sessions("session", redisStore), middleware.AuthRequired(noAuthRoutes, jwtService))
	if config.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else if config.IsTest() {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	_, err := NewApiRouter(config, router, db, jwtService, logger)
	if err != nil {
		return nil, err
	}
	return router, nil
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
