package accounts

import (
	"fmt"

	"exemplo.com/pkg/domain/entities"
)

type AccountUseCase struct {
	numberOfAccounts int
	storage          map[int]entities.Account
}

func NewAccountUseCase(numberOfAccounts int, storage map[int]entities.Account) AccountUseCase {
	return AccountUseCase{
		numberOfAccounts: numberOfAccounts,
		storage:          storage,
	}
}

func (au AccountUseCase) CreateAccount(account entities.Account) (entities.Account, error) {

	for _, storedAccount := range au.storage {
		if storedAccount.Cpf == account.Cpf {
			return entities.Account{}, fmt.Errorf("account with cpf %s already exists", account.Cpf)
		}
	}

	account, err := entities.NewAccount(account.Id, account.Name, account.Cpf, account.Secret, account.Balance)

	if err != nil {
		return entities.Account{}, fmt.Errorf("err to create an new account")
	}

	account.Id = au.numberOfAccounts
	au.storage[au.numberOfAccounts] = account
	au.numberOfAccounts++

	return account, nil
}

func (au AccountUseCase) GetBalanceById(id int) (float64, error) {
	for key, value := range au.storage {
		if value.Id == id {
			balance := au.storage[key].Balance
			return balance, nil
		}
	}
	return 0, fmt.Errorf("no id %d found", id)
}

func (au AccountUseCase) GetAccounts() []entities.Account {
	var AccountsList []entities.Account

	for _, value := range au.storage {
		AccountsList = append(AccountsList, value)
	}

	return AccountsList
}

func (au AccountUseCase) CheckAccounts(id ...int) error {
	for _, v := range id {
		if _, ok := au.storage[v]; !ok {
			return fmt.Errorf("no account found with id %d", id)
		}
	}
	return nil
}

func (au AccountUseCase) UpdateAccountBalance(id int, balance float64) error {
	account, err := au.GetAccountByID(id)
	if err != nil {
		return err
	}
	if balance < 0 {
		return fmt.Errorf("can't update account with value %f below zero", balance)
	}

	account.Balance = balance
	au.storage[id] = account

	return nil
}

func (au AccountUseCase) GetAccountByID(id int) (entities.Account, error) {
	account, ok := au.storage[id]
	if !ok {
		return entities.Account{}, fmt.Errorf("account with id %d not found", id)
	}

	return account, nil
}
