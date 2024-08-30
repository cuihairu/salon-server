package starter

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
)

type Server interface {
	ListenAndServe() error
}

type App struct {
	conf       *Config
	db         *gorm.DB
	logger     *zap.Logger
	lock       sync.RWMutex
	httpServer Server
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
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
		conf:   conf,
		logger: logger,
		db:     database,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	server := endless.NewServer(":8080", mux)
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 20
	app.httpServer = server
	return app, nil
}

func (a *App) Run() error {
	if err := a.httpServer.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		return err
	}
	return nil
}
