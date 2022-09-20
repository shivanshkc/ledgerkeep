package core

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// transactionManager implements the TransactionManager interface.
type transactionManager struct {
	collection *mongo.Collection
}

// NewTransactionManager creates a new instance of the TransactionManager.
func NewTransactionManager(collection *mongo.Collection) TransactionManager {
	return &transactionManager{collection: collection}
}

func (t *transactionManager) Create(ctx context.Context, params *ParamsTransactionCreate) (string, error) {
	panic("implement me")
}

func (t *transactionManager) Update(ctx context.Context, params *ParamsTransactionUpdate) error {
	panic("implement me")
}

func (t *transactionManager) Delete(ctx context.Context, params *ParamsTransactionDelete) error {
	panic("implement me")
}

func (t *transactionManager) Get(ctx context.Context, params *ParamsTransactionGet) (*TransactionDoc, error) {
	panic("implement me")
}

func (t *transactionManager) List(ctx context.Context, params *ParamsTransactionList) ([]*TransactionDoc, error) {
	panic("implement me")
}
