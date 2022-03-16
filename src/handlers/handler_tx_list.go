package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/models"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"

	"golang.org/x/sync/errgroup"
)

// listTransactionsQuery is the schema of the query parameters for the ListTransactions API.
type listTransactionsQuery struct {
	StartAmount *string
	EndAmount   *string
	StartTime   *string
	EndTime     *string
	AccountID   *string
	Category    *string
	NotesHint   *string

	Limit     *string
	Skip      *string
	SortField *string
	SortOrder *string
}

// ListTransactionsHandler lists transactions as per the provided queries.
func ListTransactionsHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	qValues := readListTransactionsQuery(request.URL.Query())
	filter := msi{}

	// Validating the amount values and creating the amount filter.
	amountFilter, err := getStartEndAmountFilter(qValues.StartAmount, qValues.EndAmount)
	if err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}
	// If the amount filter has any entries, we put it inside the main filter.
	if len(amountFilter) > 0 {
		filter["amount"] = amountFilter
	}

	// Validating the timestamp values and creating the timestamp filter.
	timestampFilter, err := getStartEndTimestampFilter(qValues.StartTime, qValues.EndTime)
	if err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}
	// If the timestamp filter has any entries, we put it inside the main filter.
	if len(timestampFilter) > 0 {
		filter["timestamp"] = timestampFilter
	}

	// If account ID filter is provided, we use it.
	if qValues.AccountID != nil && *qValues.AccountID != "" {
		filter["account_id"] = *qValues.AccountID
	}

	// If category filter is provided, we use it.
	if qValues.Category != nil && *qValues.Category != "" {
		filter["category"] = strings.ToLower(*qValues.Category)
	}

	// Full text search on the notes field.
	if qValues.NotesHint != nil && *qValues.NotesHint != "" {
		filter["$text"] = msi{"$search": *qValues.NotesHint}
	}

	// Parsing limit and skip to int.
	limit, skip, err := parseLimitSkip(qValues.Limit, qValues.Skip)
	if err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Parsing sortField and sortOrder.
	sortField, sortOrder, err := parseTransactionSortFieldAndOrder(qValues.SortField, qValues.SortOrder)
	if err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Params required for getting transactions that will calculate closing balance.
	databaseParamsForClosingBal := &database.ListTransactionsParams{
		Filter:          nil,
		RequiredFields:  []string{"_id", "amount", "account_id"},
		PaginationLimit: math.MaxInt64,
		PaginationSkip:  0,
		SortField:       "timestamp",
		SortOrder:       1,
		ExcludeCount:    true,
	}

	// Params required to get the transactions that will actually be shown to the user.
	databaseParamsForList := &database.ListTransactionsParams{
		Filter:          filter,
		RequiredFields:  nil,
		PaginationLimit: limit,
		PaginationSkip:  skip,
		SortField:       sortField,
		SortOrder:       sortOrder,
		ExcludeCount:    false,
	}

	errs, errCtx := errgroup.WithContext(ctx)
	// We will run two goroutines. One for getting the original transaction list
	// and the other to get the transaction list to calculate closing balances.
	// These channels will receive the final call results.
	transactionsChan := make(chan []*models.TransactionDTO, 1)
	transactionsCountChan := make(chan int, 1)
	transactions4ClosingBalChan := make(chan []*models.TransactionDTO, 1)

	// Call for transaction list that will calculate closing balances.
	errs.Go(func() error {
		// Channels will be closed upon function return.
		defer close(transactions4ClosingBalChan)
		// Database call.
		allTransactions, _, err := database.ListTransactions(errCtx, databaseParamsForClosingBal)
		if err != nil {
			return fmt.Errorf("failed to list transactions for closing bal: %w", err)
		}
		transactions4ClosingBalChan <- allTransactions
		return nil
	})

	// Call for original transaction list.
	errs.Go(func() error {
		// Channels will be closed upon function return.
		defer close(transactionsChan)
		defer close(transactionsCountChan)
		// Database call.
		transactions, count, err := database.ListTransactions(errCtx, databaseParamsForList)
		if err != nil {
			return fmt.Errorf("failed to get origin transaction list: %w", err)
		}
		transactionsChan <- transactions
		transactionsCountChan <- count
		return nil
	})

	// Checking for errors.
	if err := errs.Wait(); err != nil {
		log.Error(ctx, &logger.Entry{Payload: err})
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Processing the result of the goroutines.
	transactions := <-transactionsChan
	count := <-transactionsCountChan
	closingBalTransactions := <-transactions4ClosingBalChan

	// Getting the closing balance map.
	closingBalMap := generateTransactionClosingBalanceMap(closingBalTransactions)

	// This loop will put closing balance in all transactions.
	for idx, tx := range transactions {
		closingBal, exists := closingBalMap[tx.ID]
		if !exists {
			// This should never happen.
			message := fmt.Sprintf("did not found closing balance for tx: %s", tx.ID)
			log.Error(ctx, &logger.Entry{Payload: message})
			continue
		}
		transactions[idx].ClosingBal = closingBal
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status:  http.StatusOK,
		Headers: map[string]string{"x-total-count": fmt.Sprintf("%d", count)},
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "TRANSACTIONS_LISTED",
			Data:       transactions,
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
