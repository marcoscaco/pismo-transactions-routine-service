//go:build e2e

//esse arquivo deve contar os testes funcionais para as rotas de transação

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/platform/env"
	"github.com/shopspring/decimal"
)

func TestCreateTransactionInvalidPayloadE2E(t *testing.T) {
	baseURL := env.OrDefault("BASE_URL", "http://localhost:8080")

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "malformed JSON",
			body:       `{"account_id":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       ``,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "account_id has an incompatible type",
			body:       `{"account_id":"1","operation_type_id":1,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "account_id is zero",
			body:       `{"account_id":0,"operation_type_id":1,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "account_id is negative",
			body:       `{"account_id":-1,"operation_type_id":1,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "account_id is missing",
			body:       `{"operation_type_id":1,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "operation_type_id is zero",
			body:       `{"account_id":1,"operation_type_id":0,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "operation_type_id is negative",
			body:       `{"account_id":1,"operation_type_id":-1,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "operation_type_id is missing",
			body:       `{"account_id":1,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "operation_type_id does not exist",
			body:       `{"account_id":1,"operation_type_id":999999,"amount":"100"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "amount is zero",
			body:       `{"account_id":1,"operation_type_id":1,"amount":"0"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "amount is negative",
			body:       `{"account_id":1,"operation_type_id":1,"amount":"-10"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "amount is missing",
			body:       `{"account_id":1,"operation_type_id":1}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "amount is not a number",
			body:       `{"account_id":1,"operation_type_id":1,"amount":"abc"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "account does not exist",
			body:       `{"account_id":999999999,"operation_type_id":1,"amount":"100"}`,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+"/transactions", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("send create transaction request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d for invalid payload, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestCreateTransactionE2E(t *testing.T) {
	baseURL := env.OrDefault("BASE_URL", "http://localhost:8080")

	// Cria uma conta com documento único para não colidir com execuções anteriores.
	document := fmt.Sprintf("%011d", time.Now().UnixNano()%100000000000)
	accountBody, err := json.Marshal(map[string]string{"document_number": document})
	if err != nil {
		t.Fatalf("marshal account request body: %v", err)
	}

	accountResp, err := http.Post(baseURL+"/accounts", "application/json", bytes.NewReader(accountBody))
	if err != nil {
		t.Fatalf("send create account request: %v", err)
	}
	defer accountResp.Body.Close()
	if accountResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected create account status %d, got %d", http.StatusCreated, accountResp.StatusCode)
	}

	var account struct {
		AccountID int64 `json:"account_id"`
	}
	if err := json.NewDecoder(accountResp.Body).Decode(&account); err != nil {
		t.Fatalf("decode create account response: %v", err)
	}

	// operation_type_id 1 é "Normal Purchase", com multiplicador -1: o valor gravado deve ser negativo.
	txBody := fmt.Sprintf(`{"account_id":%d,"operation_type_id":1,"amount":"100"}`, account.AccountID)
	txResp, err := http.Post(baseURL+"/transactions", "application/json", bytes.NewBufferString(txBody))
	if err != nil {
		t.Fatalf("send create transaction request: %v", err)
	}
	defer txResp.Body.Close()
	if txResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, txResp.StatusCode)
	}

	var transaction struct {
		AccountID     int64           `json:"account_id"`
		OperationType int64           `json:"operation_type_id"`
		Amount        decimal.Decimal `json:"amount"`
	}
	if err := json.NewDecoder(txResp.Body).Decode(&transaction); err != nil {
		t.Fatalf("decode create transaction response: %v", err)
	}
	if transaction.AccountID != account.AccountID {
		t.Fatalf("expected account_id %d, got %d", account.AccountID, transaction.AccountID)
	}
	if transaction.OperationType != 1 {
		t.Fatalf("expected operation_type_id 1, got %d", transaction.OperationType)
	}
	if want := decimal.NewFromInt(-100); !transaction.Amount.Equal(want) {
		t.Fatalf("expected amount %s, got %s", want, transaction.Amount)
	}
}
