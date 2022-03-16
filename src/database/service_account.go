package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/models"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertAccount creates a new account in the database.
func InsertAccount(ctx context.Context, account *models.AccountDTO) error {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	if _, err := getAccountsCollection().InsertOne(callCtx, account); err != nil {
		// Checking if the error is a duplicate key error (already exists error).
		if mongo.IsDuplicateKeyError(err) {
			return errutils.AccountAlreadyExists()
		}
		err = fmt.Errorf("mongodb InsertOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return err
	}

	return nil
}

// IsAccountExists returns true if the account with the provided ID exists.
func IsAccountExists(ctx context.Context, accountID string) (bool, error) {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	if err := getAccountsCollection().FindOne(callCtx, bson.M{"_id": accountID}).Err(); err != nil {
		// Handling the no account found scenario.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		err = fmt.Errorf("mongodb FindOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return false, err
	}

	return true, nil
}

// IsAccountUsed returns true if even a single transaction is using the provided account, otherwise it returns false.
func IsAccountUsed(ctx context.Context, accountID string) (bool, error) {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	count, err := getTransactionsCollection().CountDocuments(callCtx, bson.M{"account_id": accountID})
	if err != nil {
		err = fmt.Errorf("mongodb CountDocuments error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return false, err
	}

	return count != 0, nil
}

// ListAccounts provides a list of all accounts.
func ListAccounts(ctx context.Context) ([]*models.AccountDTO, error) {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	cursor, err := getAccountsCollection().Find(callCtx, bson.M{})
	if err != nil {
		err = fmt.Errorf("mongodb Find error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return nil, err
	}

	var results []*models.AccountDTO
	if err := cursor.All(ctx, &results); err != nil {
		err = fmt.Errorf("mongodb cursor.All error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return nil, err
	}

	return results, nil
}

// GetAccountBalances provides a map of account IDs to their balance.
func GetAccountBalances(ctx context.Context) (map[string]float64, error) {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	// This query aggregates account_id -> balance data.
	groupStage := bson.D{{
		Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: "$account_id"},
			{Key: "balance", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		},
	}}

	// Database call.
	cursor, err := getTransactionsCollection().Aggregate(callCtx, mongo.Pipeline{groupStage})
	if err != nil {
		err = fmt.Errorf("mongodb Aggregate error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return nil, err
	}

	var results []struct {
		ID      string  `json:"id" bson:"_id"`
		Balance float64 `json:"balance" bson:"balance"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		err = fmt.Errorf("mongodb cursor.All error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return nil, err
	}

	balanceMap := map[string]float64{}
	for _, value := range results {
		balanceMap[value.ID] = value.Balance
	}

	return balanceMap, nil
}

// UpdateAccount updates an account in the database.
func UpdateAccount(ctx context.Context, accountID string, updates map[string]interface{}) error {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	// Wrapping the updates with $set operator of mongodb.
	updates = bson.M{"$set": updates}

	result, err := getAccountsCollection().UpdateOne(callCtx, bson.M{"_id": accountID}, updates)
	if err != nil {
		err = fmt.Errorf("mongodb UpdateOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return err
	}

	if result.MatchedCount == 0 {
		return errutils.AccountNotFound()
	}
	return nil
}

// DeleteAccount deletes an account in the database.
func DeleteAccount(ctx context.Context, accountID string) error {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	result, err := getAccountsCollection().DeleteOne(callCtx, bson.M{"_id": accountID})
	if err != nil {
		err = fmt.Errorf("mongodb DeleteOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return err
	}

	if result.DeletedCount == 0 {
		return errutils.AccountNotFound()
	}
	return nil
}
