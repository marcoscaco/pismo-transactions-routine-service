// DTO para representar os dados de uma conta no sistema, usado para input e output de requisições HTTP.
package dtos

import "time"

// AccountDTO representa os dados de uma conta.
type AccountDTO struct {
	ID        int64     `json:"account_id,omitempty"`
	Document  string    `json:"document_number"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
