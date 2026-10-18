// Testes unitários para o use case de transações
package usecase

import (
	"context"
	"testing"
	"time"

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
	operationTypeId := int64(1)

	transaction, err := useCase.CreateTransaction(context.Background(), accountID, amount, operationTypeId)

	assert.NoError(t, err)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, operationTypeId, transaction.OperationType)
	assert.WithinDuration(t, time.Now(), transaction.EventDate, time.Second)
}

func TestGetTransactionByID(t *testing.T) {
	fakeRepository := repository.NewFakeTransactionRepository()
	useCase := NewTransactionUseCase(fakeRepository) // Passando o repositório falso para o use case

	// Criar uma transação para testar a recuperação
	accountID := int64(1)
	amount := decimal.NewFromFloat(100.0)
	operationTypeId := int64(1)

	transaction, _ := useCase.CreateTransaction(context.Background(), accountID, amount, operationTypeId)

	_, err := useCase.GetTransactionByID(context.Background(), transaction.ID)
	assert.NoError(t, err)
	// #TODO: Implementar lógica real para buscar a transação pelo ID no repositório.
}
