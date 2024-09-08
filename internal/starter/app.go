package starter

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	config "github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type App struct {
	Config       *config.Config
	DB           *gorm.DB
	DataStore    *data.DataStore
	BizStore     *biz.BizStore
	RedisStore   redis.Store
	routers      []APIRouter
	Logger       *zap.Logger
	TokenService *utils.JWT
	Engine       *gin.Engine
	lock         sync.RWMutex
	httpServer   *HttpServer
}

func NewApp(v *viper.Viper) (*App, error) {
	conf, err := config.New(v)
	if err != nil {
		return nil, err
	}
	zapConfig, err := conf.GetZapConfig()
	if err != nil {
		return nil, err
	}
	logger, err := NewZapLogger(zapConfig)
	if err != nil {
		return nil, err
	}
	middleware.SetLogger(logger)
	dbConf, err := conf.GetDbConfig()
	if err != nil {
		return nil, err
	}
	database, err := NewDb(dbConf)
	if err != nil {
		return nil, err
	}
	redisConfig, err := conf.GetRedisConfig()
	if err != nil {
		return nil, err
	}
	redisStore, err := NewRedis(redisConfig)
	if err != nil {
		return nil, err
	}
	dataRepo, err := data.NewDataStore(database, conf, logger)
	if err != nil {
		return nil, err
	}
	jwtConfig := conf.GetJwtConfig()
	jwtService := utils.NewJWT(jwtConfig.SecretKey, jwtConfig.Expire)
	bizStore := biz.NewBizStore(conf, dataRepo, jwtService, logger)

	app := &App{
		Config:       conf,
		Logger:       logger,
		DB:           database,
		DataStore:    dataRepo,
		BizStore:     bizStore,
		RedisStore:   redisStore,
		TokenService: jwtService,
	}
	_, err = NewEngine(app)
	if err != nil {
		return nil, err
	}
	app.httpServer, err = NewHttpServer(app)
	return app, err
}

func (a *App) AddRouters(routers ...APIRouter) {
	a.routers = append(a.routers, routers...)
}

func (a *App) Run() error {
	if err := a.httpServer.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		return err
	}
	return nil
}
