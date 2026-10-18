// DTO para representar os dados de uma transação no sistema usado para input e output de requisições HTTP.
package dtos

// TransactionDTO representa os dados de uma transação.
type TransactionDTO struct {
	ID            int     `json:"id"`
	AccountID     int     `json:"account_id"`
	OperationType int     `json:"operation_type_id"`
	Amount        float64 `json:"amount"`
	EventDate     string  `json:"event_date"`
}