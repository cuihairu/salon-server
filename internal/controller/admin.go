package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"go.uber.org/zap"
)

type AdminAPI struct {
	adminBiz *biz.AdminBiz
	logger   *zap.Logger
	config   *config.Config
}
