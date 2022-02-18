package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	accounts_usecase "github.com/daniel1sender/Desafio-API/pkg/domain/accounts/usecases"
	transfers_usecase "github.com/daniel1sender/Desafio-API/pkg/domain/transfers"
	accounts_handler "github.com/daniel1sender/Desafio-API/pkg/gateways/http/accounts"
	transfers_handler "github.com/daniel1sender/Desafio-API/pkg/gateways/http/transfers"
	postgres "github.com/daniel1sender/Desafio-API/pkg/gateways/store/postgres"
	"github.com/daniel1sender/Desafio-API/pkg/gateways/store/postgres/accounts"
	"github.com/daniel1sender/Desafio-API/pkg/gateways/store/postgres/transfers"
)

type Config struct {
	DatabaseURL string `envconfig:"DB_URL" required:"true"`
	Port        string `envconfig:"API_PORT" required:"true" default:":3000"`
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.NewEntry(logger)

	var apiConfig Config
	err := envconfig.Process("", &apiConfig)
	if err != nil {
		log.WithError(err).Fatal("error while processing environment variables")
	}
	log.WithFields(logrus.Fields{"DB_URL": apiConfig.DatabaseURL, "API_PORT": apiConfig.Port}).Info("environment variables processed sucessfully")

	err = postgres.RunMigrations(apiConfig.DatabaseURL)
	if err != nil {
		log.WithError(err).Fatal("error while running migrations")
	}
	log.Info("migrations executed successfully ")

	dbPool, err := pgxpool.Connect(context.Background(), apiConfig.DatabaseURL)
	if err != nil {
		log.WithError(err).Fatal("error while connecting with the database")
	}
	log.Info("connection with the database established successfully")

	defer dbPool.Close()

	accountRepository := accounts.NewStorage(dbPool)
	accountUseCase := accounts_usecase.NewUseCase(accountRepository)
	accountHandler := accounts_handler.NewHandler(accountUseCase, log)

	transferStorage := transfers.NewStorage(dbPool)
	transferUseCase := transfers_usecase.NewUseCase(transferStorage, accountRepository)
	transferHandler := transfers_handler.NewHandler(transferUseCase, log)

	r := mux.NewRouter()
	r.HandleFunc("/accounts", accountHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/accounts", accountHandler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/accounts/{id}/balance", accountHandler.GetBalanceByID).Methods(http.MethodGet)
	r.HandleFunc("/transfers", transferHandler.Make).Methods(http.MethodPost)

	const writeTime = 60 * time.Second
	const readTime = 60 * time.Second

	server := &http.Server{
		Handler:      r,
		WriteTimeout: writeTime,
		ReadTimeout:  readTime,
		Addr:         apiConfig.Port,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), writeTime)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.WithError(err).Error("error while shut down the application")
		}
		done <- true
	}()
	log.Infof("server is running on port %s", apiConfig.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Fatal("failed to listen and serve")
	}
	<-done
}
