package starter

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/controller"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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

func NewRouter(config *config.Config, db *gorm.DB, redisStore redis.Store, logger *zap.Logger) (*gin.Engine, error) {
	// services
	jwtConfig := config.GetJwtConfig()
	jwtService := utils.NewJWT(jwtConfig.SecretKey, jwtConfig.Expire)
	// api
	router := gin.Default()
	router.Use(sessions.Sessions("session", redisStore), middleware.TokenRequired(jwtService, logger))
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
