package biz

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/utils"
	"go.uber.org/zap"
)

type Auth struct {
	userRepo   *data.UserRepository
	logger     *zap.Logger
	config     *config.Config
	jwtService *utils.JWT
}

func NewAuth(config *config.Config, jwtService *utils.JWT, userRepo *data.UserRepository, logger *zap.Logger) *Auth {
	return &Auth{
		config:     config,
		userRepo:   userRepo,
		logger:     logger,
		jwtService: jwtService,
	}
}
