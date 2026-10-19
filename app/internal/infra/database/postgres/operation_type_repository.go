package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"go.opentelemetry.io/otel"
)

type OperationTypeRepository struct {
	db *sql.DB
}

func NewOperationTypeRepository(db *sql.DB) *OperationTypeRepository {
	return &OperationTypeRepository{db: db}
}

func (r *OperationTypeRepository) GetByID(ctx context.Context, id int64) (*entity.OperationType, error) {
	ctx, span := otel.Tracer("pismo-transactions-routine-service").Start(ctx, "repository.GetOperationTypeByID")
	defer span.End()
	const query = "SELECT id, description, amount_multiplier FROM operation_types WHERE id = $1"

	operationType := &entity.OperationType{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&operationType.ID,
		&operationType.Description,
		&operationType.AmountMultiplier,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidOperationTypeID
		}
		span.RecordError(err)
		return nil, err
	}

	return operationType, nil
}
