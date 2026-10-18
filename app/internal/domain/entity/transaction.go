// Dominio de negocio para representar uma transação no sistema.
package entity

import (
	"time"
	"github.com/shopspring/decimal"
)

// Transaction representa uma transação no sistema.
type Transaction struct {
	ID            int
	AccountID     int64
	OperationType int
	Amount        decimal.Decimal
	EventDate     time.Time
}
