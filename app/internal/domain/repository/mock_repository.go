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

	account.ID = r.nextID
	created := &entity.Account{
		ID:        account.ID,
		Document:  account.Document,
		CreatedAt: account.CreatedAt,
	}
	r.accounts[account.ID] = created
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

// FakeTransactionRepository é uma implementação em memória de repository.TransactionsRepository para testes.
func NewFakeTransactionRepository() *FakeTransactionRepository {
	return &FakeTransactionRepository{
		transactions: make(map[int64]*entity.Transaction),
		nextID:       1,
	}
}

// FakeTransactionRepository é uma implementação em memória de repository.TransactionsRepository para testes.
type FakeTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[int64]*entity.Transaction
	nextID       int64
}

func (r *FakeTransactionRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	transaction.ID = r.nextID
	r.transactions[r.nextID] = transaction
	r.nextID++

	return nil
}

func (r *FakeTransactionRepository) GetTransactionByID(ctx context.Context, transactionID int64) (*entity.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transaction, exists := r.transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	return transaction, nil
}

func (r *FakeTransactionRepository) GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]*entity.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*entity.Transaction
	for _, transaction := range r.transactions {
		if transaction.AccountID == accountID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}
