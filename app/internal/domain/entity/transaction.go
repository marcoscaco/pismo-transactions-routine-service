// Dominio de negocio para representar uma transação no sistema.
package entity

import (
	"time"
)

// Transaction representa uma transação no sistema.
type Transaction struct {
	ID            int
	AccountID     int
	OperationType int
	Amount        float64
	EventDate     time.Time
}
