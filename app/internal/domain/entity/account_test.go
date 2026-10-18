// Testes unitários para a entidade Account
package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Teste unitário para a criação de uma nova conta.
func TestNewAccount(t *testing.T) {
	account := Account{
		ID:        int64(1),
		Document:  "12345678900",
		CreatedAt: time.Now(),
	}

	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, "12345678900", account.Document)
	assert.WithinDuration(t, time.Now(), account.CreatedAt, time.Second)
}
