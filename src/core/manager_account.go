package core

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// accountManager implements the AccountManager interface.
type accountManager struct{}

// NewAccountManager creates an instance of the AccountManager.
func NewAccountManager(client *mongo.Client) AccountManager {
	return &accountManager{}
}

func (a *accountManager) Create(ctx context.Context, params *ParamsAccountCreate) (string, error) {
	panic("implement me")
}

func (a *accountManager) Update(ctx context.Context, params *ParamsAccountUpdate) error {
	panic("implement me")
}

func (a *accountManager) Delete(ctx context.Context, params *ParamsAccountDelete) error {
	panic("implement me")
}

func (a *accountManager) List(ctx context.Context, params *ParamsAccountList) ([]*AccountDoc, error) {
	panic("implement me")
}
