package testutil

import (
	"context"
	"sync"

	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain/entity"
)

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
	r.accounts[account.ID] = &entity.Account{
		ID:        account.ID,
		Document:  account.Document,
		CreatedAt: account.CreatedAt,
	}
	r.nextID++
	return nil
}

func (r *FakeAccountRepository) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	if !exists {
		return nil, domain.ErrAccountNotFound
	}
	return account, nil
}

type FakeTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[int64]*entity.Transaction
	nextID       int64
}

func NewFakeTransactionRepository() *FakeTransactionRepository {
	return &FakeTransactionRepository{
		transactions: make(map[int64]*entity.Transaction),
		nextID:       1,
	}
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
		return nil, domain.ErrTransactionNotFound
	}
	return transaction, nil
}

func (r *FakeTransactionRepository) GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]*entity.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions := make([]*entity.Transaction, 0)
	for _, transaction := range r.transactions {
		if transaction.AccountID == accountID {
			transactions = append(transactions, transaction)
		}
	}
	return transactions, nil
}

type FakeOperationTypeRepository struct {
	mu             sync.RWMutex
	operationTypes map[int64]*entity.OperationType
}

func NewFakeOperationTypeRepository(operationTypes ...*entity.OperationType) *FakeOperationTypeRepository {
	repository := &FakeOperationTypeRepository{
		operationTypes: make(map[int64]*entity.OperationType, len(operationTypes)),
	}
	for _, operationType := range operationTypes {
		repository.operationTypes[operationType.ID] = operationType
	}
	return repository
}

func (r *FakeOperationTypeRepository) GetByID(ctx context.Context, id int64) (*entity.OperationType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	operationType, exists := r.operationTypes[id]
	if !exists {
		return nil, domain.ErrInvalidOperationTypeID
	}
	return operationType, nil
}
