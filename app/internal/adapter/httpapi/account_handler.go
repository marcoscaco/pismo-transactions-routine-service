// Implementação do handler para a API HTTP de contas.
package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/adapter/httpapi/dtos"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/usecase"
)

// AccountHandler é o handler para a API HTTP de contas.
type AccountHandler struct {
	accountUseCase usecase.AccountUseCase
}

// NewAccountHandler cria uma nova instância de AccountHandler.
func NewAccountHandler(accountUseCase usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{
		accountUseCase: accountUseCase,
	}
}

// CreateAccountHandler lida com a criação de uma nova conta.
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var request dtos.AccountDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountUseCase.CreateAccount(c.Request.Context(), request.Document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dtos.AccountDTO{
		ID:        account.ID,
		Document:  account.Document,
		CreatedAt: account.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}
