package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
)

func writeUseCaseError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, domain.ErrInvalidDocument),
		errors.Is(err, domain.ErrInvalidAccountID),
		errors.Is(err, domain.ErrInvalidTransactionID),
		errors.Is(err, domain.ErrInvalidTransactionAmount),
		errors.Is(err, domain.ErrInvalidOperationTypeID):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, domain.ErrAccountNotFound),
		errors.Is(err, domain.ErrTransactionNotFound):
		status = http.StatusNotFound
		message = err.Error()
	}

	c.JSON(status, gin.H{"error": message})
}
