package main

import (
	"fmt"

	"exemplo.com/pkg/domain/entities"
)

func main() {

	//criando duas contas
	fmt.Println(entities.NewAccount("John Doe", "11111111030", "123", 0))
	//account2 := entities.Account{1, "joão", "12345678910", 0}

	/* 	number := new(int)
	   	*number = 0
	   	sto := new(map[int]entities.Account)
	   	*sto = make(map[int]entities.Account)
	   	useCase := accounts.NewAccountUseCase(number, sto)
	   	fmt.Println(*number, *sto)

	   	x, err := useCase.CreateAccount(account1)

	   	fmt.Println(x, err)

	   	fmt.Println(*number, *sto) */

	// Criando a transferência

	//transfer := entities.Transfer{1, 0, 1, 2}
	//amount := 2.0

	//Checando a existência das contas:

	//fmt.Println(accounts.CheckAccounts(transfer.Account_origin_id, transfer.Account_destinantion_id))

	//Criando a transferência

	//fmt.Println(transfers.MakeTransfer(transfer))

	//Atualizando o saldo

	//accounts.UpdateAccountBalance(transfer.Account_origin_id, account1.Balance-amount)

	//accounts.UpdateAccountBalance(transfer.Account_destinantion_id, account2.Balance-amount)

}
