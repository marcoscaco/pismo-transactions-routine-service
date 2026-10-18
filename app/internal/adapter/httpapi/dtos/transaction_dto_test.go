// teste unitario para o DTO de transacao
package dtos

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Teste unitário para a função NewTransactionDTO.
func TestNewTransactionDTO(t *testing.T) {
	transaction := TransactionDTO{
		ID:            1,
		AccountID:     1,
		OperationType: 1,
		Amount:        decimal.NewFromFloat(100.0),
		EventDate:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	assert.Equal(t, int64(1), transaction.ID)
	assert.Equal(t, int64(1), transaction.AccountID)
	assert.Equal(t, int64(1), transaction.OperationType)
	assert.Equal(t, decimal.NewFromFloat(100.0), transaction.Amount)
	assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), transaction.EventDate)
}
