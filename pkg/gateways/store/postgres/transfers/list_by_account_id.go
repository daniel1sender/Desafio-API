package transfers

import (
	"context"
	"fmt"

	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
	"github.com/daniel1sender/Desafio-API/pkg/domain/transfers"
)

const listByIDStatement = "SELECT id, account_origin_id, account_destination_id, amount, created_at FROM transfers WHERE account_origin_id = $1"

func (tr TransfersRepository) ListByAccountID(ctx context.Context, accountID string) ([]entities.Transfer, error) {
	var transfersList []entities.Transfer
	var transfer entities.Transfer

	rows, err := tr.Query(ctx, listByIDStatement, accountID)
	if err != nil {
		return []entities.Transfer{}, err
	}

	for rows.Next() {
		err := rows.Scan(&transfer.ID, &transfer.AccountOriginID, &transfer.AccountDestinationID, &transfer.Amount, &transfer.CreatedAt)
		if err != nil {
			return []entities.Transfer{}, err
		}
		transfersList = append(transfersList, transfer)
	}

	if len(transfersList) == 0 {
		return []entities.Transfer{}, fmt.Errorf("error while listing transfers: %w", transfers.ErrEmptyList)
	}

	return transfersList, nil
}
