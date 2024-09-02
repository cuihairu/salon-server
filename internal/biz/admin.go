package biz

import (
	"fmt"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/utils"
	"go.uber.org/zap"
)

type AdminBiz struct {
	adminRepo  *data.AdminRepository
	logger     *zap.Logger
	config     *config.Config
	jwtService *utils.JWT
}

func NewAdminBiz(config *config.Config, jwtService *utils.JWT, adminRepo *data.AdminRepository, logger *zap.Logger) *AdminBiz {
	return &AdminBiz{
		adminRepo:  adminRepo,
		logger:     logger,
		config:     config,
		jwtService: jwtService,
	}
}

func (a *AdminBiz) Auth(username string, password string) (string, error) {
	a.logger.Info("admin login", zap.String("username", username))
	admin, err := a.adminRepo.FindByName(username)
	if err != nil {
		return "", err
	}
	if admin == nil {
		return "", fmt.Errorf("admin not found")
	}
	if !utils.VerifyPassword([]byte(admin.Password), []byte(password), []byte(admin.Salt)) {
		return "", fmt.Errorf("password not match")
	}
	token, err := a.jwtService.GenerateTokenWithGroup(admin.ID, "admin")
	if err != nil {
		return "", err
	}
	return token, nil
}
