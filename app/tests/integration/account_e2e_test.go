//go:build e2e

// esse arquivo deve implementar os testes de integração para manutençao de conta, usando o endpoint /accounts da API.
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/platform/env"
)

func TestCreateAccountE2E(t *testing.T) {
	baseURL := env.OrDefault("BASE_URL", "http://localhost:8080")
	t.Logf("Testing create account against: %s/accounts", baseURL)

	body, err := json.Marshal(map[string]string{"document_number": "12345678900"})
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/accounts", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("send create account request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestCreateAccountInvalidPayloadE2E(t *testing.T) {
	baseURL := env.OrDefault("BASE_URL", "http://localhost:8080")

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "malformed JSON",
			body:       `{"document_number":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       ``,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "document_number has an incompatible type",
			body:       `{"document_number":12345678900}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "document_number is an empty string",
			body:       `{"document_number":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "document_number is whitespace only",
			body:       `{"document_number":"   "}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "document_number is missing",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+"/accounts", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("send create account request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d for invalid payload, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestQueryAccountE2E(t *testing.T) {
	baseURL := env.OrDefault("BASE_URL", "http://localhost:8080")

	// Create an account first so the query test does not depend on seeded data.
	body, err := json.Marshal(map[string]string{"document_number": "98765432100"})
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}

	createReq, err := http.NewRequest(http.MethodPost, baseURL+"/accounts", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Fatalf("send create account request: %v", err)
	}
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createResp.StatusCode)
	}

	var created struct {
		AccountID int64 `json:"account_id"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create account response: %v", err)
	}
	if created.AccountID == 0 {
		t.Fatal("create account response returned an empty account_id")
	}

	queryReq, err := http.NewRequest(http.MethodGet, baseURL+"/accounts/"+strconv.FormatInt(created.AccountID, 10), nil)
	if err != nil {
		t.Fatalf("create query request: %v", err)
	}
	queryResp, err := http.DefaultClient.Do(queryReq)
	if err != nil {
		t.Fatalf("send query account request: %v", err)
	}
	defer queryResp.Body.Close()

	if queryResp.StatusCode != http.StatusOK {
		t.Fatalf("expected query status %d, got %d", http.StatusOK, queryResp.StatusCode)
	}

	var account struct {
		AccountID int64 `json:"account_id"`
	}
	if err := json.NewDecoder(queryResp.Body).Decode(&account); err != nil {
		t.Fatalf("decode query account response: %v", err)
	}
	if account.AccountID != created.AccountID {
		t.Fatalf("expected account_id %d, got %d", created.AccountID, account.AccountID)
	}
}

func TestQueryAccountInvalidE2E(t *testing.T) {
	baseURL := env.OrDefault("BASE_URL", "http://localhost:8080")

	tests := []struct {
		name       string
		accountID  string
		wantStatus int
	}{
		{
			name:       "non-numeric ID",
			accountID:  "abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "zero ID",
			accountID:  "0",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative ID",
			accountID:  "-1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ID overflows int64",
			accountID:  "99999999999999999999",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ID does not exist",
			accountID:  "999999999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, baseURL+"/accounts/"+tt.accountID, nil)
			if err != nil {
				t.Fatalf("create query request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("send query account request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d for account ID %q, got %d", tt.wantStatus, tt.accountID, resp.StatusCode)
			}
		})
	}
}
