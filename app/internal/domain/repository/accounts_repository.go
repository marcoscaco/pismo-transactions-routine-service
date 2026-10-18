// Interface repository para gerenciar contas no sistema.
package repository

import (
	"app/internal/domain/entity"
)

// AccountsRepository define a interface para gerenciar contas no sistema.
type AccountsRepository interface {
	// Create cria uma nova conta no sistema.
	Create(account *entity.Account) error
	// GetByID retorna uma conta pelo seu ID.
	GetByID(id int) (*entity.Account, error)
}