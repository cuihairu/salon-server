package biz

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type AccountBiz struct {
	accountRepo *data.AccountRepository
	logger      *zap.Logger
}

func NewAccountBiz(accountRepo *data.AccountRepository, logger *zap.Logger) *AccountBiz {
	return &AccountBiz{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

func (a *AccountBiz) GetAccountInfo(accountId uint) (*model.Account, error) {
	return a.accountRepo.FindByID(accountId)
}

func (a *AccountBiz) GetAllAccounts() ([]model.Account, error) {
	return a.accountRepo.FindAll()
}

func (a *AccountBiz) CreateAccount(account *model.Account) error {
	return a.accountRepo.Create(account)
}

func (a *AccountBiz) UpdateAccount(account *model.Account) error {
	return a.accountRepo.Update(account)
}

func (a *AccountBiz) DeleteAccount(accountId uint) error {
	return a.accountRepo.Delete(accountId)
}
