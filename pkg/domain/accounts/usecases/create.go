package usecases

import (
	"context"
	"fmt"

	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
)

func (au AccountUseCase) Create(ctx context.Context, name, cpf, secret string, balance int) (entities.Account, error) {

	if err := au.storage.CheckCPF(ctx, cpf); err != nil {
		return entities.Account{}, err
	}

	account, err := entities.NewAccount(name, cpf, secret, balance)
	if err != nil {
		return entities.Account{}, fmt.Errorf("error while creating an account: %w", err)
	}

	err = au.storage.Upsert(ctx, account)
	if err != nil {
		return entities.Account{}, err
	}

	return account, nil
}
