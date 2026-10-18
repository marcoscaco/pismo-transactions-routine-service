// Dominio de negocio para representar uma transação no sistema.
package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

// Transaction representa uma transação no sistema.
type Transaction struct {
	ID            int
	AccountID     int64
	OperationType int
	Amount        decimal.Decimal
	EventDate     time.Time
}
