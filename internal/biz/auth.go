package biz

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/utils"
	"go.uber.org/zap"
)

type AuthBiz struct {
	userRepo   *data.UserRepository
	logger     *zap.Logger
	config     *config.Config
	jwtService *utils.JWT
}

func NewAuthBiz(config *config.Config, jwtService *utils.JWT, dataStore *data.DataStore, logger *zap.Logger) *AuthBiz {
	return &AuthBiz{
		config:     config,
		userRepo:   dataStore.UserRepo,
		logger:     logger,
		jwtService: jwtService,
	}
}
