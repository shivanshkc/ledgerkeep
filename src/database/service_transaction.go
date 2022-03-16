package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/models"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

// CreateIndexOnTransactionField creates a B-Tree index on the specified field in the transaction collection.
func CreateIndexOnTransactionField(ctx context.Context, indexData []mongo.IndexModel) error {
	log := logger.Get()

	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	// Creating the index.
	if _, err := getTransactionsCollection().Indexes().CreateMany(callCtx, indexData); err != nil {
		err = fmt.Errorf("mongodb Indexes.CreateMany error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return err
	}

	return nil
}

// InsertTransaction creates a new transaction in the database.
// It returns the ID of the inserted document as well as the error if any.
func InsertTransaction(ctx context.Context, transaction *models.TransactionDTO) (interface{}, error) {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	result, err := getTransactionsCollection().InsertOne(callCtx, transaction)
	if err != nil {
		err = fmt.Errorf("mongodb InsertOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return "", err
	}

	return result.InsertedID, nil
}

// GetTransaction returns the transaction record matching the provided ID.
func GetTransaction(ctx context.Context, transactionID primitive.ObjectID) (*models.TransactionDTO, error) {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	result := getTransactionsCollection().FindOne(callCtx, bson.M{"_id": transactionID})
	if err := result.Err(); err != nil {
		// Handling the not-exists case.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errutils.TransactionNotFound()
		}

		err = fmt.Errorf("mongodb FindOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return nil, err
	}

	var transaction *models.TransactionDTO
	if err := result.Decode(&transaction); err != nil {
		err = fmt.Errorf("mongodb Decode error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return nil, err
	}

	return transaction, nil
}

// ListTransactions lists all the transactions that match the provided filter, pagination and sort params.
func ListTransactions(ctx context.Context, params *ListTransactionsParams) ([]*models.TransactionDTO, int, error) {
	log := logger.Get()

	errs, errCtx := errgroup.WithContext(ctx)
	// We need to fetch the list of transactions as well as the total count for pagination purposes.
	// Both these calls will be in parallel.
	// These channels will be used to process the calls.
	transactionsChan := make(chan []*models.TransactionDTO, 1)
	countChan := make(chan int64, 1)

	// Sending a parallel call to database for counting transactions.
	errs.Go(func() error {
		// Closing the channels upon function return.
		defer close(countChan)
		// If the count is not required, we don't query the DB.
		if params.ExcludeCount {
			countChan <- 0
			return nil
		}

		// Creating timeout context for the database call.
		callCtx, cancelFunc := getTimeoutContext(errCtx)
		defer cancelFunc()

		count, err := getTransactionsCollection().CountDocuments(callCtx, params.Filter)
		if err != nil {
			return fmt.Errorf("mongodb CountDocuments error: %w", err)
		}
		countChan <- count
		return nil
	})

	// Sending another parallel call for listing transactions.
	errs.Go(func() error {
		// Closing the channels upon function return.
		defer close(transactionsChan)

		opts := options.Find()
		// Setting the projection value to only include the specified fields.
		projectionBson := bson.D{}
		for _, field := range params.RequiredFields {
			projectionBson = append(projectionBson, bson.E{Key: field, Value: 1})
		}
		opts.SetProjection(projectionBson)                                                 // Specifying the fields to include.
		opts.SetLimit(int64(params.PaginationLimit)).SetSkip(int64(params.PaginationSkip)) // Specifying limit and skip.

		// Specifying sortField and sortOrder.
		opts.SetSort(bson.D{
			{Key: params.SortField, Value: params.SortOrder},
			// The second sortField is always ID, which takes the same order as the other field.
			// This brings better consistency in transaction ordering and is a must for closing balance's calculation.
			{Key: "_id", Value: params.SortOrder},
		})

		// Creating timeout context for the database call.
		callCtx, cancelFunc := getTimeoutContext(errCtx)
		defer cancelFunc()

		cursor, err := getTransactionsCollection().Find(callCtx, params.Filter, opts)
		if err != nil {
			return fmt.Errorf("mongodb Find error: %w", err)
		}

		// Decoding the transaction list into known slice type.
		var transactions []*models.TransactionDTO
		if err := cursor.All(ctx, &transactions); err != nil {
			return fmt.Errorf("mongodb cursor.All error: %w", err)
		}
		transactionsChan <- transactions
		return nil
	})

	// Checking for errors.
	if err := errs.Wait(); err != nil {
		log.Error(ctx, &logger.Entry{Payload: err})
		return nil, 0, fmt.Errorf("error in one of the goroutines: %w", err)
	}

	// Retrieving the call results.
	transactions := <-transactionsChan
	count := <-countChan

	return transactions, int(count), nil
}

// UpdateTransaction updates a transaction in the database.
func UpdateTransaction(ctx context.Context, transactionID primitive.ObjectID, updates map[string]interface{}) error {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	// Wrapping the updates with $set operator of mongodb.
	updates = bson.M{"$set": updates}

	result, err := getTransactionsCollection().UpdateOne(callCtx, bson.M{"_id": transactionID}, updates)
	if err != nil {
		err = fmt.Errorf("mongodb UpdateOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return err
	}

	if result.MatchedCount == 0 {
		return errutils.TransactionNotFound()
	}
	return nil
}

// DeleteTransaction deletes a transaction from the database.
func DeleteTransaction(ctx context.Context, transactionID primitive.ObjectID) error {
	log := logger.Get()

	// Creating timeout context for the database call.
	callCtx, cancelFunc := getTimeoutContext(ctx)
	defer cancelFunc()

	result, err := getTransactionsCollection().DeleteOne(callCtx, bson.M{"_id": transactionID})
	if err != nil {
		err = fmt.Errorf("mongodb DeleteOne error: %w", err)
		log.Error(ctx, &logger.Entry{Payload: err})
		return err
	}

	if result.DeletedCount == 0 {
		return errutils.TransactionNotFound()
	}
	return nil
}
