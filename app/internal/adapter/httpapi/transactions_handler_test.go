// testes unitários para o DTO de conta
package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/adapter/httpapi/dtos"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionUseCase := usecase.NewTransactionUseCase()
	transactionHandler := NewTransactionHandler(transactionUseCase)

	router.POST("/transactions", transactionHandler.CreateTransaction)

	requestBody := dtos.TransactionDTO{
		AccountID:     1,
		OperationType: 1,
		Amount:        decimal.NewFromFloat(100.0),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dtos.TransactionDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.AccountID)
	assert.Equal(t, int64(1), response.OperationType)
	assert.True(t, response.Amount.Equal(decimal.NewFromFloat(100.0)), "expected amount 100.0, got %s", response.Amount)
}
