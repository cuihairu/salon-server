package internal

import (
	"fmt"
	"github.com/fvbock/endless"
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
	db     *gorm.DB
	logger *zap.Logger
	lock   sync.RWMutex
	server Server
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func NewApp(db *gorm.DB, logger *zap.Logger) *App {
	app := &App{
		db:     db,
		logger: logger,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	server := endless.NewServer(":8080", mux)
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 20
	app.server = server
	//gracefulShutdown := utils.NewGracefulShutdown(30 * time.Second)
	////app.graceful = gracefulShutdown
	//gracefulShutdown.AddShutdownHook(func(ctx context.Context) error {
	//	return app.gracefulShutdown()
	//})
	return app
}

func (a *App) gracefulShutdown() error {
	return nil
}

func (a *App) Start() {
	//a.graceful.StartListening()
	//a.graceful.Wait()
	if err := a.server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
