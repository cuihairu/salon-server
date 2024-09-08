package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"go.uber.org/zap"
)

type OperationLogBiz struct {
	OperationLogRepo *data.OperationLogRepository
	logger           *zap.Logger
}

func NewOperationLogBiz(dataStore *data.DataStore, logger *zap.Logger) *OperationLogBiz {
	return &OperationLogBiz{
		OperationLogRepo: dataStore.OperationLogRepo,
		logger:           logger,
	}
}

func (l *OperationLogBiz) Log(username string, role string, ip string, location string, agent string, table string, action string, content string, status int, err string) error {
	return l.OperationLogRepo.Log(username, role, ip, location, agent, table, action, content, status, err)
}
