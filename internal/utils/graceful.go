package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type GracefulShutdown struct {
	timeout time.Duration
	hooks   []func(context.Context) error
	done    chan struct{}
}

// NewGracefulShutdown creates a new GracefulShutdown instance
func NewGracefulShutdown(timeout time.Duration) *GracefulShutdown {
	return &GracefulShutdown{
		timeout: timeout,
		done:    make(chan struct{}),
	}
}

// AddShutdownHook adds a shutdown hook that will be called during shutdown
func (g *GracefulShutdown) AddShutdownHook(hook func(context.Context) error) {
	g.hooks = append(g.hooks, hook)
}

// StartListening starts listening for termination signals
func (g *GracefulShutdown) StartListening() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		g.shutdown()
	}()
}

// Shutdown handles the actual shutdown process
func (g *GracefulShutdown) shutdown() {
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	var wg sync.WaitGroup

	for _, hook := range g.hooks {
		wg.Add(1)
		go func(h func(context.Context) error) {
			defer wg.Done()
			if err := h(ctx); err != nil {
				log.Printf("Shutdown hook error: %v", err)
			}
		}(hook)
	}

	wg.Wait()
	close(g.done)
	log.Println("Shutdown complete.")
}

// Wait blocks until the shutdown process is complete
func (g *GracefulShutdown) Wait() {
	<-g.done
}
