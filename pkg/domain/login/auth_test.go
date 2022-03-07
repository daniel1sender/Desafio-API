package login

import (
	"context"
	"errors"
	"testing"

	accounts_usecases "github.com/daniel1sender/Desafio-API/pkg/domain/accounts"
	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
	"github.com/daniel1sender/Desafio-API/pkg/gateways/store/postgres/accounts"
	"github.com/daniel1sender/Desafio-API/pkg/gateways/store/postgres/login"
)

func TestLoginUseCase_Auth(t *testing.T) {
	ctx := context.Background()
	accountRepository := accounts.NewStorage(Db)
	loginRepository := login.NewStorage(Db)
	tokenSecret := "123"
	useCase := LoginUseCase{loginRepository, accountRepository, tokenSecret}
	t.Run("should return a signed token", func(t *testing.T) {
		name := "Jonh Doe"
		cpf := "01481623559"
		secret := "123"
		balance := 10
		account, err := entities.NewAccount(name, cpf, secret, balance)
		if err != nil {
			t.Errorf("expected no error while creating a new account but got '%s'", err)
		}
		accountRepository.Upsert(ctx, account)
		tokenString, accountID, err := useCase.Auth(ctx, account.CPF, secret)
		if err != nil {
			t.Errorf("expected no error but got '%s'", err.Error())
		}
		if len(tokenString) == 0 {
			t.Error("got empty token")
		}
		err = ValidateToken(tokenString, accountID, tokenSecret)
		if err != nil {
			t.Errorf("expected no error while validanting token but got %v", err)
		}
	})

	t.Run("should return an empty token and an error when account is not found", func(t *testing.T) {
		name := "Jonh Doe"
		cpf := "01481623550"
		secret := "123"
		balance := 10
		account, err := entities.NewAccount(name, cpf, secret, balance)
		if err != nil {
			t.Errorf("expected no error while creating a new account but got '%s'", err)
		}
		tokenString, _, err := useCase.Auth(ctx, account.CPF, secret)
		if !errors.Is(err, accounts_usecases.ErrAccountNotFound) {
			t.Errorf("expected no error but got '%v'", err)
		}
		if len(tokenString) != 0 {
			t.Error("got empty token")
		}
	})
}

/* func validateToken(t *testing.T, tokenString string, accountID string, tokenSecret string) {
	t.Helper()
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		t.Errorf("expected no error but got '%v'", err)
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	if claims.Subject != accountID {
		t.Errorf("expected '%s' but got '%s'", accountID, claims.Subject)
	}
	if claims.ID == "" {
		t.Error("expected not empty id")
	}
	if !claims.VerifyExpiresAt(time.Now(), true) {
		t.Error("expected non-zero 'expires at' time")
	}
	if !claims.VerifyIssuedAt(time.Now(), true) {
		t.Error("expected non-zero 'issued at' time")
	}
} */
