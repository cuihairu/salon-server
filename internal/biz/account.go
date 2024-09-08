package biz

import (
	"fmt"
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/model"
	"go.uber.org/zap"
)

type AccountBiz struct {
	userRepo    *data.UserRepository
	accountRepo *data.AccountRepository
	logger      *zap.Logger
}

func NewAccountBiz(dataStore *data.DataStore, logger *zap.Logger) *AccountBiz {
	return &AccountBiz{
		accountRepo: dataStore.AccountRepo,
		logger:      logger,
	}
}

func (a *AccountBiz) GetAccountInfo(accountId uint) (*model.Account, error) {
	acc, err := a.accountRepo.FindByID(accountId)
	if err != nil {
		return nil, err
	}
	if acc != nil {
		return acc, nil
	}
	user, err := a.userRepo.FindByID(accountId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return a.accountRepo.FindByID(accountId)
}

func (a *AccountBiz) GetAllAccounts() ([]model.Account, error) {
	return a.accountRepo.FindAll()
}

func (a *AccountBiz) CreateAccount(account *model.Account) error {
	if account.Balance < 0 {
		return fmt.Errorf("balance cannot be negative")
	}
	if account.Consumed < 0 {
		return fmt.Errorf("consumed cannot be negative")
	}
	user, err := a.userRepo.FindByID(account.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	return a.accountRepo.Create(account)
}

func (a *AccountBiz) UpdateAccount(id uint, account *model.Account) error {
	if id == 0 {
		return fmt.Errorf("account id cannot be zero")
	}
	if account.Balance < 0 {
		return fmt.Errorf("balance cannot be negative")
	}
	if account.Consumed < 0 {
		return fmt.Errorf("consumed cannot be negative")
	}
	user, err := a.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	account.ID = id
	return a.accountRepo.Update(account)
}

func (a *AccountBiz) DeleteAccount(accountId uint) error {
	return a.accountRepo.Delete(accountId)
}
