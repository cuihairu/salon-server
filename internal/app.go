package internal

import (
	"context"
	"github.com/cuihairu/salon/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

type App struct {
	db       *gorm.DB
	logger   *zap.Logger
	lock     sync.RWMutex
	graceful *utils.GracefulShutdown
}

func NewApp(db *gorm.DB, logger *zap.Logger) *App {
	app := &App{
		db:     db,
		logger: logger,
	}
	gracefulShutdown := utils.NewGracefulShutdown(30 * time.Second)
	app.graceful = gracefulShutdown
	gracefulShutdown.AddShutdownHook(func(ctx context.Context) error {
		return app.gracefulShutdown()
	})
	return app
}

func (a *App) gracefulShutdown() error {
	return nil
}

func (a *App) Start() {
	a.graceful.StartListening()
	a.graceful.Wait()
}
