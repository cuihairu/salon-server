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

func NewAccountBiz(accountRepo *data.AccountRepository, logger *zap.Logger) *AccountBiz {
	return &AccountBiz{
		accountRepo: accountRepo,
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

	acc = &model.Account{user.ID, 0, 0}
	err = a.accountRepo.Create(acc)
	if err != nil {
		return nil, err
	}
	return acc, err
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

func (a *AccountBiz) UpdateAccount(account *model.Account) error {
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
	return a.accountRepo.Update(account)
}

func (a *AccountBiz) DeleteAccount(accountId uint) error {
	return a.accountRepo.Delete(accountId)
}
