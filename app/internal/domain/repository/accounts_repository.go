// Interface repository para gerenciar contas no sistema.
package repository

import (
	"context"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

// AccountsRepository define a interface para gerenciar contas no sistema.
type AccountsRepository interface {
	// Create cria uma nova conta no sistema.
	Create(ctx context.Context, account *entity.Account) error
	// GetByID retorna uma conta pelo seu ID.
	GetByID(ctx context.Context, id int) (*entity.Account, error)
}
