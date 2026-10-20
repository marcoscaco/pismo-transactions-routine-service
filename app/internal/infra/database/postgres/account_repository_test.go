package postgres

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
