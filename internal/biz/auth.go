package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"go.uber.org/zap"
)

type Auth struct {
	userRepo *data.UserRepository
	logger   *zap.Logger
}

func NewAuth(userRepo *data.UserRepository, logger *zap.Logger) *Auth {
	return &Auth{
		userRepo: userRepo,
		logger:   logger,
	}
}
