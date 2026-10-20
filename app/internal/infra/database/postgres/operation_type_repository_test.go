package postgres

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationTypeRepositoryGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, amount_multiplier FROM operation_types WHERE id = $1")).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "description", "amount_multiplier"}).
			AddRow(int64(1), "Normal Purchase", int32(-1)))

	repo := NewOperationTypeRepository(db)

	operationType, err := repo.GetByID(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, operationType)
	assert.Equal(t, int64(1), operationType.ID)
	assert.Equal(t, "Normal Purchase", operationType.Description)
	assert.Equal(t, int32(-1), operationType.AmountMultiplier)
	assert.NoError(t, mock.ExpectationsWereMet())
}
