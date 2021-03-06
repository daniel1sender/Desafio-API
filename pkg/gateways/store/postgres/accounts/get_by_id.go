package accounts

import (
	"context"

	"github.com/daniel1sender/Desafio-API/pkg/domain/accounts"
	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
	"github.com/jackc/pgx/v4"
)

const getByIDStatement = "SELECT id, name, cpf, secret, balance, created_at FROM accounts WHERE id = $1"

func (ar AccountRepository) GetByID(ctx context.Context, id string) (entities.Account, error) {
	account := entities.Account{}
	err := ar.QueryRow(ctx, getByIDStatement, id).Scan(&account.ID, &account.Name, &account.CPF, &account.Secret, &account.Balance, &account.CreatedAt)
	if err == pgx.ErrNoRows {
		return entities.Account{}, accounts.ErrAccountNotFound
	} else if err != nil {
		return entities.Account{}, err
	}

	return account, nil
}
