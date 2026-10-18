package domain

import "errors"

var (
	ErrInvalidDocument            = errors.New("invalid document")
	ErrInvalidAccountID           = errors.New("invalid account ID")
	ErrInvalidTransactionID       = errors.New("invalid transaction ID")
	ErrInvalidTransactionAmount   = errors.New("transaction amount must be greater than zero")
	ErrInvalidOperationTypeID     = errors.New("invalid operation type ID")
	ErrInvalidOperationMultiplier = errors.New("operation type has invalid amount multiplier")
	ErrAccountNotFound            = errors.New("account not found")
	ErrTransactionNotFound        = errors.New("transaction not found")
)
