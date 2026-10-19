// testes unitarios para o handler de contas
package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/adapter/httpapi/dtos"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/testutil"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

// Teste unitário para o handler CreateAccount.
func TestCreateAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	fakeRepository := testutil.NewFakeAccountRepository()

	accountUseCase := usecase.NewAccountUseCase(fakeRepository)
	accountHandler := NewAccountHandler(accountUseCase)

	router.POST("/accounts", accountHandler.CreateAccount)

	requestBody := dtos.AccountDTO{
		Document: "12345678900",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dtos.AccountDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "12345678900", response.Document)
}
