// DTO para representar os dados de uma transação no sistema usado para input e output de requisições HTTP.
package dtos

import (
	"github.com/shopspring/decimal"
)

// TransactionDTO representa os dados de uma transação.
type TransactionDTO struct {
	ID            int64           `json:"id,omitempty"`
	AccountID     int64           `json:"account_id"`
	OperationType int64           `json:"operation_type_id"`
	Amount        decimal.Decimal `json:"amount"`
}
