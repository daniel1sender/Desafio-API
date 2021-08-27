package accounts

import (
	"errors"
	"testing"

	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
	"github.com/daniel1sender/Desafio-API/pkg/store/accounts"
)

func TestAccountUseCase_UpdateBalance(t *testing.T) {

	t.Run("should return nil when account was updated", func(t *testing.T) {

		storage := accounts.NewStorage()
		AccountUseCase := NewUseCase(storage)
		name := "John Doe"
		cpf := "11111111030"
		secret := "123"
		balance := 10

		account, err := entities.NewAccount(name, cpf, secret, balance)
		if err != nil {
			t.Errorf("expected no error to create a new account but got '%s'", err)
		}

		storage.Upsert(account.ID, account)

		UpdateAccountError := AccountUseCase.UpdateBalance(account.ID, 20.0)

		if UpdateAccountError != nil {
			t.Errorf("expected no error but got '%s'", UpdateAccountError)
		}

	})

	t.Run("should return an error massage when account don't exists", func(t *testing.T) {

		storage := accounts.NewStorage()
		AccountUseCase := NewUseCase(storage)

		//passando qualquer id, sem criar a conta
		err := AccountUseCase.UpdateBalance("1", 20.0)

		if err != accounts.ErrIDNotFound {
			t.Errorf("expected '%s' but got '%s'", accounts.ErrIDNotFound, err)
		}

	})

	t.Run("should return an error message when balance account is less than zero", func(t *testing.T) {

		storage := accounts.NewStorage()
		AccountUseCase := NewUseCase(storage)

		name := "John Doe"
		cpf := "11111111030"
		secret := "123"
		balance := 10

		account, err := entities.NewAccount(name, cpf, secret, balance)
		if err != nil {
			t.Errorf("expected no error to create a new account but got '%s'", err)
		}

		storage.Upsert(account.ID, account)

		err = AccountUseCase.UpdateBalance(account.ID, -10)

		if !errors.Is(err, ErrBalanceLessZero) {
			t.Errorf("expected '%s' but got '%s'", ErrBalanceLessZero, err)
		}

	})
}
