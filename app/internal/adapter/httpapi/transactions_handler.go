// handler para a API HTTP de transações.
package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/adapter/httpapi/dtos"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/usecase"
)

// TransactionHandler é o handler para a API HTTP de transações.
type TransactionHandler struct {
	transactionUseCase usecase.TransactionUseCase
}

// NewTransactionHandler cria uma nova instância de TransactionHandler.
func NewTransactionHandler(transactionUseCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

// CreateTransaction trata a requisição HTTP para criar uma nova transação.
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var request dtos.TransactionDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionUseCase.CreateTransaction(c.Request.Context(), request.AccountID, request.Amount, request.OperationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dtos.TransactionDTO{
		AccountID:     transaction.AccountID,
		Amount:        transaction.Amount,
		OperationType: transaction.OperationType,
	}
	c.JSON(http.StatusCreated, resp)
}
