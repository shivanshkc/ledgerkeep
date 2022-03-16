package database

import (
	"context"
	"time"

	"github.com/shivanshkc/ledgerkeep/src/configs"
	"github.com/shivanshkc/ledgerkeep/src/database/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	accountsCollectionName     = "accounts"
	transactionsCollectionName = "transactions"
)

// ListTransactionsParams is the schema of params required by the ListTransactions operation.
type ListTransactionsParams struct {
	// Filter is the search filter for the transactions.
	Filter map[string]interface{}
	// RequiredFields is the list of fields to be included in the response.
	RequiredFields []string
	// PaginationLimit is the max amount of transactions in the response.
	PaginationLimit int
	// PaginationSkip is the initial offset of the list.
	PaginationSkip int
	// SortField is the name of the field by which the transactions will be sorted.
	SortField string
	// SortOrder is the order of sorting.
	SortOrder int
	// ExcludeCount is a flag to control whether the total count of the transaction should also be calculated or not.
	ExcludeCount bool
}

// getAccountsCollection provides the accounts mongoDB collection.
func getAccountsCollection() *mongo.Collection {
	conf := configs.Get()
	return mongodb.GetClient().Database(conf.Mongo.DatabaseName).Collection(accountsCollectionName)
}

// getTransactionsCollection provides the transactions mongoDB collection.
func getTransactionsCollection() *mongo.Collection {
	conf := configs.Get()
	return mongodb.GetClient().Database(conf.Mongo.DatabaseName).Collection(transactionsCollectionName)
}

// getTimeoutContext provides the timeout context for database operations.
func getTimeoutContext(parent context.Context) (context.Context, context.CancelFunc) {
	conf := configs.Get()
	timeoutDuration := time.Duration(conf.Mongo.OperationTimeoutSec) * time.Second
	return context.WithTimeout(parent, timeoutDuration)
}
