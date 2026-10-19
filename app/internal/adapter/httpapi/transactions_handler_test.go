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
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/testutil"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/usecase"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	transactionRepository := testutil.NewFakeTransactionRepository()
	operationTypeRepository := testutil.NewFakeOperationTypeRepository(
		&entity.OperationType{ID: 4, AmountMultiplier: 1},
	)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepository, operationTypeRepository)
	transactionHandler := NewTransactionHandler(transactionUseCase)

	router.POST("/transactions", transactionHandler.CreateTransaction)

	requestBody := dtos.TransactionDTO{
		AccountID:     1,
		OperationType: 4,
		Amount:        decimal.NewFromFloat(100.0),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]json.RawMessage
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotContains(t, response, "id")

	var transactionResponse dtos.TransactionDTO
	err = json.Unmarshal(w.Body.Bytes(), &transactionResponse)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), transactionResponse.AccountID)
	assert.Equal(t, int64(4), transactionResponse.OperationType)
	assert.True(t, transactionResponse.Amount.Equal(decimal.NewFromFloat(100.0)), "expected amount 100.0, got %s", transactionResponse.Amount)
}
