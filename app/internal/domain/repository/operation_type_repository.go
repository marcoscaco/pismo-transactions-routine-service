// Interface repository para gerenciar tipos de operação no sistema.
package repository

import (
    "context"

    "github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

// OperationTypeRepository define a interface para gerenciar tipos de operação.
type OperationTypeRepository interface {
    // GetByID retorna o tipo de operação pelo seu ID.
    GetByID(ctx context.Context, id int64) (*entity.OperationType, error)
}
