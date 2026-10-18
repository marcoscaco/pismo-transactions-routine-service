// Usecase de manipulação de transações no sistema - lugar centralizado para a lógica de negócio relacionada a transações.
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/repository"
	"github.com/shopspring/decimal"
)

// TransactionUseCase é a interface que define os métodos para manipulação de transações.
type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, accountID int64, amount decimal.Decimal, operationTypeId int64) (*entity.Transaction, error)
	GetTransactionByID(ctx context.Context, transactionID int64) (*entity.Transaction, error)
}

// transactionUseCase é a implementação da interface TransactionUseCase.
type transactionUseCase struct {
	// #TODO: Adicionar dependências.
	transactionRepository repository.TransactionsRepository
}

// NewTransactionUseCase cria uma nova instância de transactionUseCase.
func NewTransactionUseCase(transactionRepository repository.TransactionsRepository) TransactionUseCase {
	return &transactionUseCase{
		transactionRepository: transactionRepository,
	}
}

// CreateTransaction cria uma nova transação com os parâmetros fornecidos.
func (uc *transactionUseCase) CreateTransaction(ctx context.Context, accountID int64, amount decimal.Decimal, operationTypeId int64) (*entity.Transaction, error) {

	// valida accountID
	if accountID <= 0 {
		return nil, errors.New("invalid account ID") // ID inválido
	}

	// valida amount
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than zero") // Valor inválido
	}

	// valida operationTypeId
	if operationTypeId <= 0 {
		return nil, errors.New("invalid operation type ID") // ID de tipo de operação inválido
	}

	transaction := &entity.Transaction{
		AccountID:     accountID,
		Amount:        amount,
		EventDate:     time.Now(),
		OperationType: operationTypeId,
	}
	// #TODO: Implementar lógica para persistir a transação no repositório.
	err := uc.transactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
func (uc *transactionUseCase) GetTransactionByID(ctx context.Context, transactionID int64) (*entity.Transaction, error) {
	// #TODO: Implementar lógica para buscar a transação pelo ID no repositório.
	return uc.transactionRepository.GetTransactionByID(ctx, transactionID)
}
