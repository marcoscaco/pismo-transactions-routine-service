// Usecase de manipulação de transações no sistema - lugar centralizado para a lógica de negócio relacionada a transações.
package usecase

import (
	"context"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"time"
	"github.com/shopspring/decimal"
)

// TransactionUseCase é a interface que define os métodos para manipulação de transações.
type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, accountID int64, amount decimal.Decimal, operationTypeId int) (*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, accountID int64) (*entity.Transaction, error)
}

// transactionUseCase é a implementação da interface TransactionUseCase.
type transactionUseCase struct {
	// #TODO: Adicionar dependências.
}

// NewTransactionUseCase cria uma nova instância de transactionUseCase.
func NewTransactionUseCase() TransactionUseCase {
	return &transactionUseCase{}
}

// CreateTransaction cria uma nova transação com os parâmetros fornecidos.
func (uc *transactionUseCase) CreateTransaction(ctx context.Context, accountID int64, amount decimal.Decimal, operationTypeId int) (*entity.Transaction, error) {
	transaction := &entity.Transaction{
		AccountID: accountID,
		Amount:    amount,
		EventDate: time.Now(),
		OperationType: operationTypeId,
	}
	// #TODO: Implementar lógica para persistir a transação no repositório.
	return transaction, nil
}
func (uc *transactionUseCase) GetTransactionByID(ctx context.Context, accountID int64) (*entity.Transaction, error) {
	// #TODO: Implementar lógica para buscar a transação pelo ID no repositório.
	return nil, nil
}