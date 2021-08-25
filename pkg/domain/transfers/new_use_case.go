package transfers

import (
	"github.com/daniel1sender/Desafio-API/pkg/store/transfers"
)

type TransferUseCase struct {
	storage transfers.TransferStorage
}

func NewTransferUseCase(storage transfers.TransferStorage) TransferUseCase {
	return TransferUseCase{
		storage: storage,
	}
}