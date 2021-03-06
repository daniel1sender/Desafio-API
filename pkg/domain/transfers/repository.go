package transfers

import (
	"context"

	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
)

type Repository interface {
	Insert(ctx context.Context, transfer entities.Transfer) error
	ListByAccountID(ctx context.Context, accountID string) ([]entities.Transfer, error)
}

type AccountRepository interface {
	GetBalanceByID(ctx context.Context, id string) (int, error)
	GetByID(ctx context.Context, id string) (entities.Account, error)
	Upsert(ctx context.Context, account entities.Account) error
}
