// DTO para representar os dados de uma conta no sistema, usado para input e output de requisições HTTP.
package dtos

// AccountDTO representa os dados de uma conta.
type AccountDTO struct {
	ID       int `json:"id"`
	Document string `json:"document_number"`
	CreatedAt string `json:"created_at"`
}