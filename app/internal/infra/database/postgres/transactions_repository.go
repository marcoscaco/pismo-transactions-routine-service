// implementação concreta do repositório de transações usando PostgreSQL
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

// TransactionRepository é a implementação concreta do repositório de transações usando PostgreSQL.
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository cria uma nova instância de TransactionRepository.
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// Create insere uma nova transação no banco de dados.

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) error {
	query := "INSERT INTO transactions_operations (account_id, operation_type_id, amount) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, transaction.AccountID, transaction.OperationType, transaction.Amount).Scan(&transaction.ID)
	if err != nil {
		return mapPostgresError(err)
	}
	return nil
}

// GetTransactionByID recupera uma transação pelo seu ID.
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, transactionID int64) (*entity.Transaction, error) {
	query := "SELECT id, account_id, operation_type_id, amount, event_date FROM transactions_operations WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, transactionID)

	transaction := &entity.Transaction{}
	err := row.Scan(&transaction.ID, &transaction.AccountID, &transaction.OperationType, &transaction.Amount, &transaction.EventDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTransactionNotFound
		}
		return nil, err
	}
	return transaction, nil
}

// GetTransactionsByAccountID recupera todas as transações associadas a uma conta específica.
func (r *TransactionRepository) GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]*entity.Transaction, error) {
	query := "SELECT id, account_id, operation_type_id, amount, event_date FROM transactions_operations WHERE account_id = $1"
	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]*entity.Transaction, 0)
	for rows.Next() {
		transaction := &entity.Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.OperationType, &transaction.Amount, &transaction.EventDate)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
