package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type UserBiz struct {
	userRepo *data.UserRepository
	logger   *zap.Logger
}

func NewUserBiz(dataStore *data.DataStore, logger *zap.Logger) *UserBiz {
	return &UserBiz{
		userRepo: dataStore.UserRepo,
		logger:   logger,
	}
}

func (b *UserBiz) GetUserByID(id uint) (*model.User, error) {
	return b.userRepo.FindByID(id)
}

func (b *UserBiz) GetAllUsers() ([]model.User, error) {
	return b.userRepo.FindAll()
}

func (b *UserBiz) UpdateUser(id uint, user *model.User) error {
	existingUser, err := b.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Update user fields
	existingUser.Name = user.Name
	existingUser.Phone = user.Phone
	existingUser.Address = user.Address
	existingUser.Birthday = user.Birthday

	return b.userRepo.Update(existingUser)
}

func (b *UserBiz) DeleteUser(id uint) error {
	return b.userRepo.Delete(id)
}
