// Interface para o repositório de transações, definindo os métodos que devem ser implementados para persistência e recuperação de dados relacionados a transações.
package repository

import (
	"context"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

// TransactionsRepository define os métodos para interagir com o repositório de transações.
type TransactionsRepository interface {
	// CreateTransaction cria uma nova transação no repositório.
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) error

	// GetTransactionByID recupera uma transação pelo seu ID.
	GetTransactionByID(ctx context.Context, accountID int64) (*entity.Transaction, error)

	// GetTransactionsByAccountID recupera todas as transações associadas a uma conta específica.
	GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]*entity.Transaction, error)
}