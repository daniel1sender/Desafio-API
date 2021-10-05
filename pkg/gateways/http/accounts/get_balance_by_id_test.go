package accounts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	accounts_usecase "github.com/daniel1sender/Desafio-API/pkg/domain/accounts"
	server_http "github.com/daniel1sender/Desafio-API/pkg/gateways/http"
	"github.com/daniel1sender/Desafio-API/pkg/gateways/store/accounts"
)

func TestGetBalanceByID(t *testing.T) {
	t.Run("should return 200 and the account balance", func(t *testing.T) {

		expectedBalance := 20
		useCase := accounts_usecase.UseCaseMock{Balance: expectedBalance, Error: nil}
		h := NewHandler(&useCase)

		newRequest, _ := http.NewRequest(http.MethodGet, "accounts/{id}/balance", nil)
		newResponse := httptest.NewRecorder()

		h.GetBalanceByID(newResponse, newRequest)

		var response ByIDresponse

		json.Unmarshal(newResponse.Body.Bytes(), &response)

		if newResponse.Code != http.StatusOK {
			t.Errorf("expected %d but got %d", http.StatusOK, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected %s but got %s", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		if response.Balance != expectedBalance {
			t.Errorf("expected '%d' but got '%d'", expectedBalance, response.Balance)
		}

	})

	t.Run("should return 404 and a error message when account is not found by id", func(t *testing.T) {

		expectedBalance := 0
		useCase := accounts_usecase.UseCaseMock{Balance: expectedBalance,
			Error: accounts.ErrIDNotFound}

		h := NewHandler(&useCase)

		newRequest, _ := http.NewRequest(http.MethodGet, "accounts/{id}/balance", nil)
		newResponse := httptest.NewRecorder()

		h.GetBalanceByID(newResponse, newRequest)

		var response ByIDresponse

		_ = json.Unmarshal(newResponse.Body.Bytes(), &response)

		if newResponse.Code != http.StatusNotFound {
			t.Errorf("expected '%d' but got '%d'", http.StatusNotFound, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		if response.Balance != 0 {
			t.Errorf("expected null balance but got '%d'", response.Balance)
		}

	})

}