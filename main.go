package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	transfers_usecase "github.com/daniel1sender/Desafio-API/pkg/domain/transfers"
	accounts_handler "github.com/daniel1sender/Desafio-API/pkg/gateways/http/accounts"
	transfers_handler "github.com/daniel1sender/Desafio-API/pkg/gateways/http/transfers"
)

func main() {

	accountStorage := accounts_storage.NewStorage()
	accountUseCase := accounts_usecase.NewUseCase(accountStorage)
	accountHandler := accounts_handler.NewHandler(accountUseCase)

	transferStorage := transfers_storage.NewStorage()
	transferUseCase := transfers_usecase.NewUseCase(transferStorage, accountStorage)
	transferHandler := transfers_handler.NewHandler(transferUseCase)

	r := mux.NewRouter()
	r.HandleFunc("/accounts", accountHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/accounts", accountHandler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/accounts/{id}/balance", accountHandler.GetBalanceByID).Methods(http.MethodGet)
	r.HandleFunc("/transfers", transferHandler.Make).Methods(http.MethodPost)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("failed to listen and serve: %s", err)
	}
}
