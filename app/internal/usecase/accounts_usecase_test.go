// testes para o usecase de contas
package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

// Teste unitário para a função CreateAccount.
func TestCreateAccount(t *testing.T) {
	fakeRepository := repository.NewFakeAccountRepository()

	useCase := NewAccountUseCase(fakeRepository)

	account, err := useCase.CreateAccount(context.Background(), "12345678900")
	assert.NoError(t, err)
	assert.Equal(t, "12345678900", account.Document)
	assert.WithinDuration(t, time.Now(), account.CreatedAt, time.Second)
}

// Teste unitário para a função GetAccountByID.
func TestGetAccountByID(t *testing.T) {

	fakeRepository := repository.NewFakeAccountRepository()
	useCase := NewAccountUseCase(fakeRepository)

	// Criar uma conta para testar a recuperação
	account, _ := useCase.CreateAccount(context.Background(), "12345678900")

	_, err := useCase.GetAccountByID(context.Background(), account.ID)
	assert.NoError(t, err)
	// #TODO: Implementar lógica real para buscar a conta pelo ID no repositório.
}
