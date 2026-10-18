package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

// FakeAccountRepository é uma implementação em memória de repository.AccountRepository para testes.
type FakeAccountRepository struct {
	mu       sync.RWMutex
	accounts map[int64]*entity.Account
	nextID   int64
}

func NewFakeAccountRepository() *FakeAccountRepository {
	return &FakeAccountRepository{
		accounts: make(map[int64]*entity.Account),
		nextID:   1,
	}
}

func (r *FakeAccountRepository) Create(ctx context.Context, account *entity.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	created := &entity.Account{
		ID:        r.nextID,
		Document:  account.Document,
		CreatedAt: account.CreatedAt,
	}
	r.accounts[r.nextID] = created
	r.nextID++

	return nil
}

func (r *FakeAccountRepository) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	if !exists {
		return nil, errors.New("account not found")
	}

	return account, nil
}
