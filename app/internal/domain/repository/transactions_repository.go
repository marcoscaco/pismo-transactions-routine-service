// Interface para o repositório de transações, definindo os métodos que devem ser implementados para persistência e recuperação de dados relacionados a transações.
package repository

import (
	"app/internal/domain/entity"
)

// TransactionsRepository define os métodos para interagir com o repositório de transações.
type TransactionsRepository interface {
	// CreateTransaction cria uma nova transação no repositório.
	CreateTransaction(transaction *entity.Transaction) error

	// GetTransactionByID recupera uma transação pelo seu ID.
	GetTransactionByID(id int) (*entity.Transaction, error)

	// GetTransactionsByAccountID recupera todas as transações associadas a uma conta específica.
	GetTransactionsByAccountID(accountID int) ([]*entity.Transaction, error)
}