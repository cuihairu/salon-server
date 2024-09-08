package starter

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/controller"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type APIRouter interface {
	RegisterRoutes(router *gin.RouterGroup)
	Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error
}

func NewApiRouter(app *App) (*gin.RouterGroup, error) {
	routers := []APIRouter{
		&controller.AdminAPI{},
		&controller.UserAPI{},
		&controller.AuthAPI{},
		&controller.CategoryAPI{},
		&controller.ServiceAPI{},
		&controller.AccountAPI{},
		&controller.OrderAPI{},
		&controller.MemberAPI{},
	}

	apiGroup := app.Engine.Group("/api")
	app.routers = routers
	for _, router := range routers {
		if err := router.Initialize(app.Config, app.BizStore, app.Logger); err != nil {
			return nil, err
		}
	}
	for _, router := range routers {
		router.RegisterRoutes(apiGroup)
	}
	return apiGroup, nil
}

func NewRouter(app *App) (*gin.Engine, error) {
	// services
	app.Engine = gin.Default()
	app.Engine.Use(sessions.Sessions("session", app.RedisStore), middleware.TokenRequired(app.TokenService, app.Logger))
	if app.Config.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else if app.Config.IsTest() {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	_, err := NewApiRouter(app)
	if err != nil {
		return nil, err
	}
	return app.Engine, nil
}
