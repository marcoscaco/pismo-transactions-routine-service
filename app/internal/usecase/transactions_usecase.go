// Usecase de manipulação de transações no sistema - lugar centralizado para a lógica de negócio relacionada a transações.
package usecase

import (
	"context"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
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
		return nil, domain.ErrInvalidAccountID
	}

	// valida amount
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, domain.ErrInvalidTransactionAmount
	}

	// valida operationTypeId
	if operationTypeId <= 0 {
		return nil, domain.ErrInvalidOperationTypeID
	}

	// transacoes do tipo purchase, purchacestalments, widrawal devem ser negativas, mas o valor passado deve ser positivo, e o sistema deve tratar a inversão do sinal.
	amount = calculateAmount(operationTypeId, amount)

	transaction := &entity.Transaction{
		AccountID:     accountID,
		Amount:        amount,
		EventDate:     time.Now(),
		OperationType: operationTypeId,
	}

	err := uc.transactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func calculateAmount(operationTypeId int64, amount decimal.Decimal) decimal.Decimal {
	if operationTypeId == 1 || operationTypeId == 2 || operationTypeId == 3 {
		amount = amount.Neg()
	}
	return amount
}
func (uc *transactionUseCase) GetTransactionByID(ctx context.Context, transactionID int64) (*entity.Transaction, error) {
	if transactionID <= 0 {
		return nil, domain.ErrInvalidTransactionID
	}
	return uc.transactionRepository.GetTransactionByID(ctx, transactionID)
}
