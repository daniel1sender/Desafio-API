package usecases

import (
	"context"
	"testing"

	accounts_usecase "github.com/daniel1sender/Desafio-API/pkg/domain/accounts"
	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
	accounts_repository "github.com/daniel1sender/Desafio-API/pkg/gateways/store/postgres/accounts"
	"github.com/daniel1sender/Desafio-API/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAccountUseCase_GetBalanceByID(t *testing.T) {
	repository := accounts_repository.NewStorage(Db)
	accountUsecase := NewUseCase(repository)
	ctx := context.Background()

	t.Run("should return an account balance when the id is found", func(t *testing.T) {
		name := "John Doe"
		cpf := "11111111030"
		secret := "123"
		balance := 10

		account, _ := entities.NewAccount(name, cpf, secret, balance)
		repository.Upsert(ctx, account)

		getBalance, err := accountUsecase.GetBalanceByID(ctx, account.ID)
		assert.Nil(t, err)
		assert.NotNil(t, getBalance)
		assert.Equal(t, getBalance, account.Balance)
	})

	t.Run("should return an account balance equal to zero and an error when the id wasn't found", func(t *testing.T) {

		id := "247eade0-53b2-4ee5-9a3b-0da9afba5700"
		tests.DeleteAllAccounts(Db)

		getBalance, err := accountUsecase.GetBalanceByID(ctx, id)

		assert.NotNil(t, getBalance)
		assert.Equal(t, err, accounts_usecase.ErrAccountNotFound)
	})
}
