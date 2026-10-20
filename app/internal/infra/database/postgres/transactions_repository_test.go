package postgres

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepositoryCreateTransaction_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	amount := decimal.RequireFromString("-50.00")
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO transactions_operations (account_id, operation_type_id, amount) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs(int64(1), int64(4), amount).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(10)))

	repo := NewTransactionRepository(db)
	transaction := &entity.Transaction{
		AccountID:     1,
		OperationType: 4,
		Amount:        amount,
	}

	err = repo.CreateTransaction(context.Background(), transaction)

	require.NoError(t, err)
	assert.Equal(t, int64(10), transaction.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
