// Use case para manipulação de contas no sistema - lugar centralizado para a lógica de negócio relacionada a contas.
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/repository"
)

// AccountUseCase é a interface que define os métodos para manipulação de contas.
type AccountUseCase interface {
	CreateAccount(ctx context.Context, document string) (*entity.Account, error)
	GetAccountByID(ctx context.Context, id int64) (*entity.Account, error)
}

// accountUseCase é a implementação da interface AccountUseCase.
type accountUseCase struct {
	accountRepository repository.AccountsRepository
}

// NewAccountUseCase cria uma nova instância de accountUseCase.
func NewAccountUseCase(accountRepository repository.AccountsRepository) AccountUseCase {
	return &accountUseCase{
		accountRepository: accountRepository,
	}
}

// CreateAccount cria uma nova conta com o documento fornecido.
func (uc *accountUseCase) CreateAccount(ctx context.Context, document string) (*entity.Account, error) {

	// validar se o documento é válido
	if document == "" {
		return nil, errors.New("document cannot be empty") // Documento vazio
	}

	account := &entity.Account{
		Document:  document,
		CreatedAt: time.Now(),
	}

	// implementa a criação da conta no repositório
	err := uc.accountRepository.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByID retorna uma conta pelo seu ID.
func (uc *accountUseCase) GetAccountByID(ctx context.Context, id int64) (*entity.Account, error) {
	// implementa a busca da conta pelo ID no repositório

	// validar se o ID é válido
	if id <= 0 {
		return nil, errors.New("invalid account ID") // ID inválido
	}

	account, err := uc.accountRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found") // Conta não encontrada
	}

	return account, nil
}
