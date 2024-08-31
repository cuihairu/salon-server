package biz

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/data"
	"go.uber.org/zap"
)

type Auth struct {
	userRepo *data.UserRepository
	logger   *zap.Logger
	config   *config.Config
}

func NewAuth(config *config.Config, userRepo *data.UserRepository, logger *zap.Logger) *Auth {
	return &Auth{
		config:   config,
		userRepo: userRepo,
		logger:   logger,
	}
}
