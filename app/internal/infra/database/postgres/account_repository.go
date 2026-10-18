// implementacao concreta do repositorio de contas usando PostgreSQL
package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

// AccountRepository é a implementação concreta do repositório de contas usando PostgreSQL.
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository cria uma nova instância de AccountRepository.
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// Create insere uma nova conta no banco de dados.
func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) error {

	query := "INSERT INTO accounts (document_number, created_at) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, account.Document, time.Now()).Scan(&account.ID)
	if err != nil {
		return err
	}
	account.CreatedAt = time.Now()
	return nil
}

// GetByID busca uma conta pelo seu ID no banco de dados.
func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	query := "SELECT id, document_number, created_at FROM accounts WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	account := &entity.Account{}
	err := row.Scan(&account.ID, &account.Document, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Conta não encontrada
		}
		return nil, err
	}
	return account, nil
}
