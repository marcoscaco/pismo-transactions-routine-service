// Dominio de negocio para representar uma transação no sistema.
package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction representa uma transação no sistema.
type Transaction struct {
	ID            int64
	AccountID     int64
	OperationType int64
	Amount        decimal.Decimal
	EventDate     time.Time
}
