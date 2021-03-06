package transfers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daniel1sender/Desafio-API/pkg/domain/accounts"
	"github.com/daniel1sender/Desafio-API/pkg/domain/entities"
	"github.com/daniel1sender/Desafio-API/pkg/domain/transfers"
	"github.com/daniel1sender/Desafio-API/pkg/domain/transfers/usecases"
	server_http "github.com/daniel1sender/Desafio-API/pkg/gateways/http"
	"github.com/sirupsen/logrus"
)

func TestHandlerMake(t *testing.T) {
	log := logrus.NewEntry(logrus.New())

	t.Run("should return 201 and a transfer when it's been sucessfully created", func(t *testing.T) {
		transfer := entities.Transfer{AccountOriginID: "1", AccountDestinationID: "0", Amount: 10}
		useCase := transfers.UseCaseMock{Transfer: transfer}
		h := NewHandler(&useCase, log)
		createRequest := TransferRequest{transfer.AccountDestinationID, transfer.Amount}
		request, _ := json.Marshal(createRequest)
		newRequest, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(request))
		ctx := context.WithValue(newRequest.Context(), server_http.ContextAccountID, transfer.AccountOriginID)
		newResponse := httptest.NewRecorder()
		h.Make(newResponse, newRequest.WithContext(ctx))
		ExpectedCreateAt := transfer.CreatedAt.Format(server_http.DateLayout)

		var response TransferResponse
		_ = json.Unmarshal(newResponse.Body.Bytes(), &response)

		if newResponse.Code != http.StatusCreated {
			t.Errorf("expected '%d' but got '%d'", http.StatusCreated, newResponse.Code)
		}

		headerType := newResponse.Header().Get("content-type")
		if headerType != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, headerType)
		}

		if response.AccountOriginID != transfer.AccountOriginID {
			t.Errorf("expected '%s' but got '%s'", transfer.AccountOriginID, response.AccountOriginID)
		}

		if response.AccountDestinationID != transfer.AccountDestinationID {
			t.Errorf("expected '%s' but got '%s'", transfer.AccountDestinationID, response.AccountDestinationID)
		}

		if response.Amount != transfer.Amount {
			t.Errorf("expected '%d' but got '%d'", transfer.Amount, response.Amount)
		}

		if response.CreatedAt != ExpectedCreateAt {
			t.Errorf("expected '%s' but got '%s'", ExpectedCreateAt, response.CreatedAt)
		}

	})

	t.Run("should return 400 and an error when it failed to decode the request", func(t *testing.T) {

		transfer := entities.Transfer{AccountOriginID: "1", AccountDestinationID: "0", Amount: 10}
		useCase := transfers.UseCaseMock{Transfer: transfer}
		h := NewHandler(&useCase, log)
		b := []byte{}
		newRequest, _ := http.NewRequest(http.MethodPost, "transfers", bytes.NewBuffer(b))
		newResponse := httptest.NewRecorder()
		ctx := context.WithValue(newRequest.Context(), server_http.ContextAccountID, transfer.AccountOriginID)
		h.Make(newResponse, newRequest.WithContext(ctx))

		var responseReason server_http.Error
		json.Unmarshal(newResponse.Body.Bytes(), &responseReason)

		if newResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status '%d' but got '%d'", http.StatusBadRequest, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		expected := "invalid request body"
		if responseReason.Reason != expected {
			t.Errorf("expected '%s' but got '%s'", expected, responseReason.Reason)
		}

	})

	t.Run("should return 400 and an error when amount is less or equal zero", func(t *testing.T) {

		transfer := entities.Transfer{AccountOriginID: "1", AccountDestinationID: "0", Amount: -10}
		useCase := transfers.UseCaseMock{Transfer: transfer, Error: entities.ErrAmountLessOrEqualZero}
		h := NewHandler(&useCase, log)

		createRequest := TransferRequest{transfer.AccountDestinationID, transfer.Amount}
		request, _ := json.Marshal(createRequest)
		newRequest, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(request))
		newResponse := httptest.NewRecorder()
		ctx := context.WithValue(newRequest.Context(), server_http.ContextAccountID, transfer.AccountOriginID)
		h.Make(newResponse, newRequest.WithContext(ctx))

		var responseReason server_http.Error
		json.Unmarshal(newResponse.Body.Bytes(), &responseReason)

		if newResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status '%d' but got '%d'", http.StatusBadRequest, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		if responseReason.Reason != entities.ErrAmountLessOrEqualZero.Error() {
			t.Errorf("expected '%s' but got '%s'", entities.ErrAmountLessOrEqualZero.Error(), responseReason.Reason)
		}

	})

	t.Run("should return 400 and an error when the transfer is to the same account", func(t *testing.T) {

		transfer := entities.Transfer{AccountOriginID: "0", AccountDestinationID: "0", Amount: 10}
		useCase := transfers.UseCaseMock{Transfer: transfer, Error: entities.ErrSameAccountTransfer}
		h := NewHandler(&useCase, log)

		createRequest := TransferRequest{transfer.AccountDestinationID, transfer.Amount}
		request, _ := json.Marshal(createRequest)
		newRequest, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(request))
		newResponse := httptest.NewRecorder()
		ctx := context.WithValue(newRequest.Context(), server_http.ContextAccountID, transfer.AccountOriginID)
		h.Make(newResponse, newRequest.WithContext(ctx))

		var responseReason server_http.Error
		json.Unmarshal(newResponse.Body.Bytes(), &responseReason)

		if newResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status '%d' but got '%d'", http.StatusBadRequest, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		if responseReason.Reason != entities.ErrSameAccountTransfer.Error() {
			t.Errorf("expected '%s' but got '%s'", entities.ErrSameAccountTransfer.Error(), responseReason.Reason)
		}
	})

	t.Run("should return 400 and an error when transfer origin id is not found", func(t *testing.T) {

		transfer := entities.Transfer{AccountOriginID: "0", AccountDestinationID: "1", Amount: 10}
		useCase := transfers.UseCaseMock{Transfer: transfer, Error: usecases.ErrOriginAccountNotFound}
		h := NewHandler(&useCase, log)

		createRequest := TransferRequest{transfer.AccountDestinationID, transfer.Amount}
		request, _ := json.Marshal(createRequest)
		newRequest, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(request))
		newResponse := httptest.NewRecorder()
		ctx := context.WithValue(newRequest.Context(), server_http.ContextAccountID, transfer.AccountOriginID)
		h.Make(newResponse, newRequest.WithContext(ctx))

		var responseReason server_http.Error
		json.Unmarshal(newResponse.Body.Bytes(), &responseReason)

		if newResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status '%d' but got '%d'", http.StatusBadRequest, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		if responseReason.Reason != usecases.ErrOriginAccountNotFound.Error() {
			t.Errorf("expected '%s' but got '%s'", accounts.ErrAccountNotFound.Error(), responseReason.Reason)
		}
	})

	t.Run("should return 400 and an error when transfer destination id is not found", func(t *testing.T) {

		transfer := entities.Transfer{AccountOriginID: "0", AccountDestinationID: "1", Amount: 10}
		useCase := transfers.UseCaseMock{Transfer: transfer, Error: usecases.ErrDestinationAccountNotFound}
		h := NewHandler(&useCase, log)

		createRequest := TransferRequest{transfer.AccountDestinationID, transfer.Amount}
		request, _ := json.Marshal(createRequest)
		newRequest, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(request))
		newResponse := httptest.NewRecorder()
		ctx := context.WithValue(newRequest.Context(), server_http.ContextAccountID, transfer.AccountOriginID)
		h.Make(newResponse, newRequest.WithContext(ctx))

		var responseReason server_http.Error
		json.Unmarshal(newResponse.Body.Bytes(), &responseReason)

		if newResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status '%d' but got '%d'", http.StatusBadRequest, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		if responseReason.Reason != usecases.ErrDestinationAccountNotFound.Error() {
			t.Errorf("expected '%s' but got '%s'", accounts.ErrAccountNotFound.Error(), responseReason.Reason)
		}
	})

	t.Run("should return 400 and error when origin account doesn't have sufficient funds", func(t *testing.T) {

		transfer := entities.Transfer{AccountOriginID: "0", AccountDestinationID: "1", Amount: 10}
		useCase := transfers.UseCaseMock{Transfer: transfer, Error: usecases.ErrInsufficientFunds}
		h := NewHandler(&useCase, log)

		createRequest := TransferRequest{transfer.AccountDestinationID, transfer.Amount}
		request, _ := json.Marshal(createRequest)
		newRequest, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(request))
		newResponse := httptest.NewRecorder()
		ctx := context.WithValue(newRequest.Context(), server_http.ContextAccountID, transfer.AccountOriginID)
		h.Make(newResponse, newRequest.WithContext(ctx))

		var responseReason server_http.Error
		json.Unmarshal(newResponse.Body.Bytes(), &responseReason)

		if newResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status '%d' but got '%d'", http.StatusBadRequest, newResponse.Code)
		}

		if newResponse.Header().Get("content-type") != server_http.JSONContentType {
			t.Errorf("expected '%s' but got '%s'", server_http.JSONContentType, newResponse.Header().Get("content-type"))
		}

		if responseReason.Reason != usecases.ErrInsufficientFunds.Error() {
			t.Errorf("expected '%s' but got '%s'", usecases.ErrInsufficientFunds.Error(), responseReason.Reason)
		}

	})
}
