package postgres

import (
	"errors"

	"github.com/lib/pq"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
)

func mapPostgresError(err error) error {
	var postgresError *pq.Error
	if !errors.As(err, &postgresError) || postgresError.Code != "23503" {
		return err
	}

	switch postgresError.Constraint {
	case "fk_operation_type":
		return domain.ErrInvalidOperationTypeID
	case "fk_account":
		return domain.ErrAccountNotFound
	default:
		return err
	}
}
