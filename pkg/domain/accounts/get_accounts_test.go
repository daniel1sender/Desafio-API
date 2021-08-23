package accounts

import (
	"testing"

	"exemplo.com/pkg/domain/entities"
	"exemplo.com/pkg/store/accounts"
)

func TestAccountUseCase_GetAccounts(t *testing.T) {

	t.Run("should return a full list of accounts", func(t *testing.T) {

		storage := accounts.NewAccountStorage()
		AccountUseCase := NewAccountUseCase(storage)
		name := "John Doe"
		cpf := "11111111030"
		secret := "123"
		balance := 10

		account, err := entities.NewAccount(name, cpf, secret, balance)
		if err != nil {
			t.Errorf("Expected nil error to create a new account but got '%s'", err)
		}

		storage.UpdateStorage(account.ID, account)

		getAccounts := AccountUseCase.GetAccounts()

		if len(getAccounts) == 0 {
			t.Error("expected a full list of accounts")
		}

	})

	t.Run("should return an empty list", func(t *testing.T) {

		storage := accounts.NewAccountStorage()
		AccountUseCase := NewAccountUseCase(storage)

		getAccounts := AccountUseCase.GetAccounts()

		if len(getAccounts) != 0 {
			t.Error("expected an empty list")
		}

	})

}