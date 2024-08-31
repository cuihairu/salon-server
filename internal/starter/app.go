package starter

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type App struct {
	config     *Config
	db         *gorm.DB
	logger     *zap.Logger
	lock       sync.RWMutex
	httpServer *HttpServer
}

func NewApp(v *viper.Viper) (*App, error) {
	conf, err := New(v)
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
	dbConf, err := conf.GetDbConfig()
	if err != nil {
		return nil, err
	}
	database, err := NewDb(dbConf)
	if err != nil {
		return nil, err
	}
	app := &App{
		config: conf,
		logger: logger,
		db:     database,
	}
	app.httpServer, err = NewHttpServer(conf, database, logger)
	return app, err
}

func (a *App) Run() error {
	if err := a.httpServer.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		return err
	}
	return nil
}
