package starter

import (
	"fmt"
	config2 "github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type App struct {
	config     *config2.Config
	db         *gorm.DB
	logger     *zap.Logger
	lock       sync.RWMutex
	httpServer *HttpServer
}

func NewApp(v *viper.Viper) (*App, error) {
	conf, err := config2.New(v)
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
	redis, err := NewRedis(redisConfig)
	if err != nil {
		return nil, err
	}
	app := &App{
		config: conf,
		logger: logger,
		db:     database,
	}
	app.httpServer, err = NewHttpServer(conf, database, redis, logger)
	return app, err
}

func (a *App) Run() error {
	if err := a.httpServer.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		return err
	}
	return nil
}
