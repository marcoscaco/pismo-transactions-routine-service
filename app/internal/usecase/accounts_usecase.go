// Use case para manipulação de contas no sistema - lugar centralizado para a lógica de negócio relacionada a contas.
package usecase

import (
	"context"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"time"
)

// AccountUseCase é a interface que define os métodos para manipulação de contas.
type AccountUseCase interface {
	CreateAccount(ctx context.Context, document string) (*entity.Account, error)
	GetAccountByID(ctx context.Context, id int) (*entity.Account, error)
}

// accountUseCase é a implementação da interface AccountUseCase.
type accountUseCase struct {
	// #TODO: Adicionar dependências.
}

// NewAccountUseCase cria uma nova instância de accountUseCase.
func NewAccountUseCase() AccountUseCase {
	return &accountUseCase{}
}

// CreateAccount cria uma nova conta com o documento fornecido.
func (uc *accountUseCase) CreateAccount(ctx context.Context, document string) (*entity.Account, error) {
	account := &entity.Account{
		Document:  document,
		CreatedAt: time.Now(),
	}
	// #TODO: Implementar lógica para persistir a conta no repositório.
	return account, nil
}

// GetAccountByID retorna uma conta pelo seu ID.
func (uc *accountUseCase) GetAccountByID(ctx context.Context, id int) (*entity.Account, error) {
	// #TODO: Implementar lógica para buscar a conta pelo ID no repositório.
	return nil, nil
}