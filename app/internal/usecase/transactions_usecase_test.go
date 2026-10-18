// Testes unitários para o use case de transações
package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Teste unitário para a função CreateTransaction.
func TestCreateTransaction(t *testing.T) {
	fakeRepository := repository.NewFakeTransactionRepository()
	useCase := NewTransactionUseCase(fakeRepository) // Passando o repositório falso para o use case
	accountID := int64(1)
	amount := decimal.NewFromFloat(100.0)

	//operaçao 4 aumenta o saldo entao trabalha positivo
	operationTypeId := int64(4)

	transaction, err := useCase.CreateTransaction(context.Background(), accountID, amount, operationTypeId)

	assert.NoError(t, err)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, operationTypeId, transaction.OperationType)
	assert.WithinDuration(t, time.Now(), transaction.EventDate, time.Second)

	// operacoes 1 2 3 devem diminuir o saldo, entao devem reportar o valor negativo
	operationTypeId = int64(1)
	transaction, err = useCase.CreateTransaction(context.Background(), accountID, amount, operationTypeId)

	assert.NoError(t, err)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, amount.Neg(), transaction.Amount) // Verifica se o valor é negativo
	assert.Equal(t, operationTypeId, transaction.OperationType)
	assert.WithinDuration(t, time.Now(), transaction.EventDate, time.Second)

}

func TestGetTransactionByID(t *testing.T) {
	fakeRepository := repository.NewFakeTransactionRepository()
	useCase := NewTransactionUseCase(fakeRepository)

	transaction, err := useCase.CreateTransaction(context.Background(), 1, decimal.NewFromInt(100), 1)
	assert.NoError(t, err)

	got, err := useCase.GetTransactionByID(context.Background(), transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, transaction, got)
}

func TestCreateTransactionSentinelErrors(t *testing.T) {
	useCase := NewTransactionUseCase(repository.NewFakeTransactionRepository())
	amount := decimal.NewFromInt(10)

	_, err := useCase.CreateTransaction(context.Background(), 0, amount, 1)
	assert.ErrorIs(t, err, domain.ErrInvalidAccountID)

	_, err = useCase.CreateTransaction(context.Background(), 1, decimal.Zero, 1)
	assert.ErrorIs(t, err, domain.ErrInvalidTransactionAmount)

	_, err = useCase.CreateTransaction(context.Background(), 1, amount, 0)
	assert.ErrorIs(t, err, domain.ErrInvalidOperationTypeID)
}

func TestGetTransactionByIDNotFound(t *testing.T) {
	useCase := NewTransactionUseCase(repository.NewFakeTransactionRepository())

	_, err := useCase.GetTransactionByID(context.Background(), 1)
	assert.ErrorIs(t, err, domain.ErrTransactionNotFound)

	_, err = useCase.GetTransactionByID(context.Background(), 0)
	assert.ErrorIs(t, err, domain.ErrInvalidTransactionID)
}
