// Testes unitários para o use case de transações
package usecase

import (
	"context"
	"testing"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/testutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Teste unitário para a função CreateTransaction.
func TestCreateTransaction(t *testing.T) {
	fakeRepository := testutil.NewFakeTransactionRepository()
	operationTypeRepository := testutil.NewFakeOperationTypeRepository(
		&entity.OperationType{ID: 40, AmountMultiplier: 1},
		&entity.OperationType{ID: 10, AmountMultiplier: -1},
	)
	useCase := NewTransactionUseCase(fakeRepository, operationTypeRepository)
	accountID := int64(1)
	amount := decimal.NewFromFloat(100.0)

	//operaçao 4 aumenta o saldo entao trabalha positivo
	operationTypeId := int64(40)

	transaction, err := useCase.CreateTransaction(context.Background(), accountID, amount, operationTypeId)

	assert.NoError(t, err)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, operationTypeId, transaction.OperationType)

	// operacoes 1 2 3 devem diminuir o saldo, entao devem reportar o valor negativo
	operationTypeId = int64(10)
	transaction, err = useCase.CreateTransaction(context.Background(), accountID, amount, operationTypeId)

	assert.NoError(t, err)
	assert.Equal(t, accountID, transaction.AccountID)
	assert.Equal(t, amount.Neg(), transaction.Amount) // Verifica se o valor é negativo
	assert.Equal(t, operationTypeId, transaction.OperationType)

}

func TestGetTransactionByID(t *testing.T) {
	fakeRepository := testutil.NewFakeTransactionRepository()
	operationTypeRepository := testutil.NewFakeOperationTypeRepository(
		&entity.OperationType{ID: 10, AmountMultiplier: -1},
	)
	useCase := NewTransactionUseCase(fakeRepository, operationTypeRepository)

	transaction, err := useCase.CreateTransaction(context.Background(), 1, decimal.NewFromInt(100), 10)
	assert.NoError(t, err)

	got, err := useCase.GetTransactionByID(context.Background(), transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, transaction, got)
}

func TestCreateTransactionSentinelErrors(t *testing.T) {
	useCase := NewTransactionUseCase(
		testutil.NewFakeTransactionRepository(),
		testutil.NewFakeOperationTypeRepository(&entity.OperationType{ID: 10, AmountMultiplier: -1}),
	)
	amount := decimal.NewFromInt(10)

	_, err := useCase.CreateTransaction(context.Background(), 0, amount, 1)
	assert.ErrorIs(t, err, domain.ErrInvalidAccountID)

	_, err = useCase.CreateTransaction(context.Background(), 1, decimal.Zero, 1)
	assert.ErrorIs(t, err, domain.ErrInvalidTransactionAmount)

	_, err = useCase.CreateTransaction(context.Background(), 1, amount, 0)
	assert.ErrorIs(t, err, domain.ErrInvalidOperationTypeID)
}

func TestGetTransactionByIDNotFound(t *testing.T) {
	useCase := NewTransactionUseCase(
		testutil.NewFakeTransactionRepository(),
		testutil.NewFakeOperationTypeRepository(),
	)

	_, err := useCase.GetTransactionByID(context.Background(), 1)
	assert.ErrorIs(t, err, domain.ErrTransactionNotFound)

	_, err = useCase.GetTransactionByID(context.Background(), 0)
	assert.ErrorIs(t, err, domain.ErrInvalidTransactionID)
}
