package data

import (
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"gorm.io/gorm"
)

type OperationLogRepository struct {
	utils.Repository[model.OperationLog]
}

func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{
		Repository: utils.NewRepository[model.OperationLog](db),
	}
}

func (l *OperationLogRepository) Log(username string, role string, ip string, location string, agent string, table string, action string, content string, err string) error {
	loginLog := model.OperationLog{
		Username: username,
		Role:     role,
		Ip:       ip,
		Location: location,
		Agent:    agent,
		Table:    table,
		Action:   action,
		Content:  content,
		Err:      err,
	}
	return l.Create(&loginLog)
}
