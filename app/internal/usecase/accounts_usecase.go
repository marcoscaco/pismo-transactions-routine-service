// Use case para manipulação de contas no sistema - lugar centralizado para a lógica de negócio relacionada a contas.
package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
	// Start a child span for the use‑case operation
	var span trace.Span
	ctx, span = otel.Tracer("pismo-transactions-routine-service").Start(ctx, "usecase.CreateAccount")
	defer span.End()

	// validar se o documento é válido
	if document == "" {
		span.RecordError(domain.ErrInvalidDocument)
		return nil, domain.ErrInvalidDocument
	}

	// validar que nao sao apenas espaços em branco
	if strings.TrimSpace(document) == "" {
		span.RecordError(domain.ErrInvalidDocument)
		return nil, domain.ErrInvalidDocument
	}

	account := &entity.Account{
		Document:  document,
		CreatedAt: time.Now(),
	}

	// implementa a criação da conta no repositório
	err := uc.accountRepository.Create(ctx, account)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return account, nil
}

// GetAccountByID retorna uma conta pelo seu ID.
func (uc *accountUseCase) GetAccountByID(ctx context.Context, id int64) (*entity.Account, error) {
	var span trace.Span
	ctx, span = otel.Tracer("pismo-transactions-routine-service").Start(ctx, "usecase.GetAccountByID")
	defer span.End()

	// implementa a busca da conta pelo ID no repositório

	// validar se o ID é válido
	if id <= 0 {
		span.RecordError(domain.ErrInvalidAccountID)
		return nil, domain.ErrInvalidAccountID
	}

	account, err := uc.accountRepository.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	if account == nil {
		span.RecordError(domain.ErrAccountNotFound)
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}
