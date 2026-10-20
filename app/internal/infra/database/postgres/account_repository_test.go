package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepositoryCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	createdAt := time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (document_number) VALUES ($1) RETURNING id, created_at")).
		WithArgs("12345678900").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(int64(1), createdAt))

	repo := NewAccountRepository(db)
	account := &entity.Account{Document: "12345678900"}

	err = repo.Create(context.Background(), account)

	require.NoError(t, err)
	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, createdAt, account.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepositoryCreate_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbErr := errors.New("connection refused")
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO accounts (document_number) VALUES ($1) RETURNING id, created_at")).
		WithArgs("12345678900").
		WillReturnError(dbErr)

	repo := NewAccountRepository(db)
	account := &entity.Account{Document: "12345678900"}

	err = repo.Create(context.Background(), account)

	require.ErrorIs(t, err, dbErr)
	assert.Zero(t, account.ID)
	assert.True(t, account.CreatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepositoryGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	createdAt := time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, document_number, created_at FROM accounts WHERE id = $1")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "document_number", "created_at"}).
			AddRow(int64(1), "12345678900", createdAt))

	repo := NewAccountRepository(db)

	account, err := repo.GetByID(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, account)
	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, "12345678900", account.Document)
	assert.Equal(t, createdAt, account.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepositoryGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, document_number, created_at FROM accounts WHERE id = $1")).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "document_number", "created_at"}))

	repo := NewAccountRepository(db)

	account, err := repo.GetByID(context.Background(), 99)

	require.ErrorIs(t, err, domain.ErrAccountNotFound)
	assert.Nil(t, account)
	assert.NoError(t, mock.ExpectationsWereMet())
}
