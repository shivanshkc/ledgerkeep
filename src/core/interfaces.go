package core

import (
	"context"
)

// AccountManager encapsulates all methods required to manage accounts.
type AccountManager interface {
	// Create a new account.
	Create(ctx context.Context, params *ParamsAccountCreate) (accountID string, err error)
	// Update an existing account.
	Update(ctx context.Context, params *ParamsAccountUpdate) (err error)
	// Delete an account.
	Delete(ctx context.Context, params *ParamsAccountDelete) (err error)
	// List accounts with filters.
	List(ctx context.Context, params *ParamsAccountList) (accounts []*AccountDoc, err error)
}

// TransactionManager encapsulates all methods required to manage transactions.
type TransactionManager interface {
	// Create a new transaction.
	Create(ctx context.Context, params *ParamsTransactionCreate) (transactionID string, err error)
	// Update an existing transaction.
	Update(ctx context.Context, params *ParamsTransactionUpdate) (err error)
	// Delete a transaction.
	Delete(ctx context.Context, params *ParamsTransactionDelete) (err error)
	// Get a single transaction.
	Get(ctx context.Context, params *ParamsTransactionGet) (transaction *TransactionDoc, err error)
	// List transactions with filters.
	List(ctx context.Context, params *ParamsTransactionList) (transactions []*TransactionDoc, err error)
}

// StatisticsManager encapsulates all methods required to show statistics.
type StatisticsManager interface{}
