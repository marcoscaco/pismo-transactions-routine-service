// DTO para representar os dados de uma transação no sistema usado para input e output de requisições HTTP.
package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

// TransactionDTO representa os dados de uma transação.
type TransactionDTO struct {
	ID            int64           `json:"id"`
	AccountID     int64           `json:"account_id"`
	OperationType int64           `json:"operation_type_id"`
	Amount        decimal.Decimal `json:"amount"`
	EventDate     time.Time       `json:"event_date"`
}
