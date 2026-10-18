// Arquivo com os testes unitários para o DTO de conta.
package dtos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Teste unitário para a função NewAccountDTO.
func TestNewAccountDTO(t *testing.T) {
	account := AccountDTO{
		ID:        1,
		Document:  "12345678900",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, "12345678900", account.Document)
	assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), account.CreatedAt)
}
