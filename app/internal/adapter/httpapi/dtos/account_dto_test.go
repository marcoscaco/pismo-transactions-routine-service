// Arquivo com os testes unitários para o DTO de conta.
package dtos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Teste unitário para a função NewAccountDTO.
func TestNewAccountDTO(t *testing.T) {
	account := AccountDTO{
		ID:        1,
		Document:  "12345678900",
		CreatedAt: "2023-01-01T00:00:00Z",
	}

	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, "12345678900", account.Document)
	assert.Equal(t, "2023-01-01T00:00:00Z", account.CreatedAt)
}
