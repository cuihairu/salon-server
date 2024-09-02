package biz

import (
	"fmt"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
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

func (a *AdminBiz) RefreshJwt(id uint) (string, error) {
	admin, err := a.adminRepo.FindByID(id)
	if err != nil {
		return "", err
	}
	if admin == nil {
		return "", fmt.Errorf("admin not found")
	}
	token, err := a.jwtService.GenerateTokenWithGroup(admin.ID, "admin")
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *AdminBiz) createDefaultAdmin(password string) (*model.Admin, error) {
	// create admin
	passwordHash, saltHash, err := utils.GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	admin := &model.Admin{
		Name:     "admin",
		Password: passwordHash,
		Salt:     saltHash,
	}
	err = a.adminRepo.Create(admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (a *AdminBiz) Auth(username string, password string) (string, *model.Admin, error) {
	a.logger.Info("admin login", zap.String("username", username))
	if len(username) == 0 {
		return "", nil, fmt.Errorf("username cannot be empty")
	}
	if len(password) == 0 {
		return "", nil, fmt.Errorf("password cannot be empty")
	}
	admin, err := a.adminRepo.FindByName(username)
	if err != nil {
		return "", nil, err
	}
	if admin != nil {
		if !utils.VerifyPassword(admin.Password, []byte(password), []byte(admin.Salt)) {
			return "", nil, fmt.Errorf("password not match")
		}
	} else if username == "admin" {
		admin, err = a.createDefaultAdmin(password)
		if err != nil {
			return "", nil, err
		}
	} else {
		return "", nil, fmt.Errorf("admin not found")
	}

	token, err := a.jwtService.GenerateTokenWithGroup(admin.ID, "admin")
	if err != nil {
		return "", nil, err
	}
	return token, admin, nil
}

func (a *AdminBiz) UpdatePassword(id uint, newPassword string) error {
	a.logger.Info("admin login", zap.Uint("user id", id))
	admin, err := a.adminRepo.FindByID(id)
	if err != nil {
		return err
	}
	if admin == nil {
		return fmt.Errorf("admin not found")
	}
	passwordHash, saltHash, err := utils.GeneratePasswordHash(newPassword)
	if err != nil {
		return err
	}
	admin.Password = passwordHash
	admin.Salt = saltHash
	return a.adminRepo.Update(admin)
}

func (a *AdminBiz) GetAllAdmins() ([]model.Admin, error) {
	return a.adminRepo.FindAll()
}
