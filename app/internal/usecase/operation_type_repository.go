package usecase

import (
	"context"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

type OperationTypeRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.OperationType, error)
}
