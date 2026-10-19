// implementacao concreta do repositorio de contas usando PostgreSQL
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("pismo-transactions-routine-service").Start(ctx, "repository.CreateAccount")
	defer span.End()

	query := "INSERT INTO accounts (document_number, created_at) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, account.Document, time.Now()).Scan(&account.ID)
	if err != nil {
		span.RecordError(err)
		return err
	}
	account.CreatedAt = time.Now()
	return nil
}

// GetByID busca uma conta pelo seu ID no banco de dados.
func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	ctx, span := otel.Tracer("pismo-transactions-routine-service").Start(ctx, "repository.GetAccountByID")
	defer span.End()
	query := "SELECT id, document_number, created_at FROM accounts WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	account := &entity.Account{}
	err := row.Scan(&account.ID, &account.Document, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		span.RecordError(err)
		return nil, err
	}
	return account, nil
}
