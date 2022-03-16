package handlers

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shivanshkc/ledgerkeep/src/models"
)

type msi = map[string]interface{}

// prepareNewTransaction validates all params of createTransactionBody struct.
// It also creates a *models.TransactionDTO struct out of it.
func prepareNewTransaction(body *createTransactionBody) (*models.TransactionDTO, error) {
	// Validating transaction amount.
	if body.Amount == 0 {
		return nil, errInvalidTxAmount
	}

	// Validating account ID.
	if !accountIDRegexp.MatchString(body.AccountID) {
		return nil, errInvalidAccountID
	}

	// Validating category.
	allowedCats := getAllowedCategoriesForTxAmount(body.Amount)
	if !stringPresentCaseInsensitive(body.Category, allowedCats) {
		return nil, errInvalidTxCategory
	}

	return &models.TransactionDTO{
		Amount:    body.Amount,
		Timestamp: body.Timestamp,
		AccountID: body.AccountID,
		Category:  strings.ToLower(body.Category),
		Notes:     body.Notes,
	}, nil
}

// prepareUpdateTransactionQuery validates all params of the updateTransactionBody and creates a map of updates.
// Any nil parameters are ignored in the process.
//
// This method requires the current transaction object for the following reasons:
//
// 1. If the user did not provide an amount but provided a category, we need the current transaction object to determine
// the valid list of categories.
//
// 2. If the user did not provide a category but provided an amount, we need the current category to determine if the
// transaction should be a debit or credit.
//
// This whole mess is obviously because some categories are only valid for debit transactions and some are only valid
// for credit transactions.
func prepareUpdateTransactionQuery(body *updateTransactionBody, currentTx *models.TransactionDTO) (msi, error) {
	updates := msi{}

	// Validating transaction amount.
	if body.Amount != nil {
		if *body.Amount == 0 {
			return nil, errInvalidTxAmount
		}
		// If a new category has not been provided,
		// then this new amount should be in agreement with the current category.
		if body.Category == nil {
			validCase1 := stringPresentCaseInsensitive(currentTx.Category, allowedCreditCategories) && *body.Amount > 0
			validCase2 := stringPresentCaseInsensitive(currentTx.Category, allowedDebitCategories) && *body.Amount < 0
			if !validCase1 && !validCase2 {
				return nil, errAmountCategoryMismatch
			}
		}

		updates["amount"] = *body.Amount
	}

	// Validating transaction timestamp.
	if body.Timestamp != nil {
		updates["timestamp"] = *body.Timestamp
	}

	// Validating account ID.
	if body.AccountID != nil {
		if !accountIDRegexp.MatchString(*body.AccountID) {
			return nil, errInvalidAccountID
		}
		updates["account_id"] = *body.AccountID
	}

	// Validating category.
	if body.Category != nil {
		var amount float64
		if body.Amount != nil {
			amount = *body.Amount
		} else {
			amount = currentTx.Amount
		}

		allowedCats := getAllowedCategoriesForTxAmount(amount)
		if !stringPresentCaseInsensitive(*body.Category, allowedCats) {
			return nil, errInvalidTxCategory
		}
		updates["category"] = *body.Category
	}

	// Validating notes.
	if body.Notes != nil {
		updates["notes"] = *body.Notes
	}

	return updates, nil
}

// readListTransactionsQuery reads the request query values and loads them into *listTransactionsQuery type.
func readListTransactionsQuery(values url.Values) *listTransactionsQuery {
	qValues := &listTransactionsQuery{}

	if values.Has("start_amount") {
		startAmount := values.Get("start_amount")
		qValues.StartAmount = &startAmount
	}

	if values.Has("end_amount") {
		endAmount := values.Get("end_amount")
		qValues.EndAmount = &endAmount
	}

	if values.Has("start_time") {
		startTime := values.Get("start_time")
		qValues.StartTime = &startTime
	}

	if values.Has("end_time") {
		endTime := values.Get("end_time")
		qValues.EndTime = &endTime
	}

	if values.Has("account_id") {
		accountID := values.Get("account_id")
		qValues.AccountID = &accountID
	}

	if values.Has("category") {
		category := values.Get("category")
		qValues.Category = &category
	}

	if values.Has("notes_hint") {
		notes := values.Get("notes_hint")
		qValues.NotesHint = &notes
	}

	if values.Has("limit") {
		limit := values.Get("limit")
		qValues.Limit = &limit
	}

	if values.Has("skip") {
		skip := values.Get("skip")
		qValues.Skip = &skip
	}

	if values.Has("sort_field") {
		sortField := values.Get("sort_field")
		qValues.SortField = &sortField
	}

	if values.Has("sort_order") {
		sortOrder := values.Get("sort_order")
		qValues.SortOrder = &sortOrder
	}

	return qValues
}

// getStartEndAmountFilter creates a MongoDB style filter for start and end amount of a transaction.
func getStartEndAmountFilter(startAmount *string, endAmount *string) (msi, error) {
	filter := msi{}

	if startAmount != nil && *startAmount != "" {
		amount, err := strconv.ParseFloat(*startAmount, 64)
		if err != nil {
			return nil, errInvalidStartAmount
		}
		filter["$gte"] = amount
	}
	if endAmount != nil && *endAmount != "" {
		amount, err := strconv.ParseFloat(*endAmount, 64)
		if err != nil {
			return nil, errInvalidEndAmount
		}
		filter["$lte"] = amount
	}

	return filter, nil
}

// getStartEndTimestampFilter creates a MongoDB style filter for start and end timestamps of a transaction.
func getStartEndTimestampFilter(startTime *string, endTime *string) (msi, error) {
	filter := msi{}

	if startTime != nil && *startTime != "" {
		txTimestamp, err := parseTransactionTimestampString(*startTime)
		if err != nil {
			return nil, errInvalidTxTimestamp
		}
		filter["$gte"] = txTimestamp
	}
	if endTime != nil && *endTime != "" {
		txTimestamp, err := parseTransactionTimestampString(*endTime)
		if err != nil {
			return nil, errInvalidTxTimestamp
		}
		filter["$lte"] = txTimestamp
	}

	return filter, nil
}

// parseLimitSkip parses the limit, skip from *string type (as received from the query params) to int type.
// If any of limit and skip or both are nil, default values are used.
func parseLimitSkip(limit *string, skip *string) (int, int, error) {
	var parsedLimit, parsedSkip int

	if limit == nil || *limit == "" {
		parsedLimit = defaultLimit
	} else {
		parsed, err := strconv.ParseInt(*limit, 10, 64)
		if err != nil || parsed < 1 || parsed > defaultLimit { // Limit should be positive int.
			return 0, 0, errInvalidLimit
		}
		parsedLimit = int(parsed)
	}

	if skip == nil || *skip == "" {
		parsedSkip = defaultSkip
	} else {
		parsed, err := strconv.ParseInt(*skip, 10, 64)
		if err != nil || parsed < 0 { // Skip should be a non-negative int.
			return 0, 0, errInvalidSkip
		}
		parsedSkip = int(parsed)
	}

	return parsedLimit, parsedSkip, nil
}

// parseTransactionSortFieldAndOrder parses the sortField and sortOrder values for a transaction.
// sortField has to be one of allowedTransactionSortFields, and sortOrder has to be one of allowedSortOrders.
// If any of these are nil, default values are used.
func parseTransactionSortFieldAndOrder(sortField *string, sortOrder *string) (string, int, error) {
	var parsedSortField string
	var parsedSortOrder int

	if sortField == nil || *sortField == "" {
		parsedSortField = defaultTransactionSortField
	} else {
		if !stringPresentCaseInsensitive(*sortField, allowedTransactionSortFields) {
			return "", 0, errInvalidTxSortField
		}
		parsedSortField = strings.ToLower(*sortField)
	}

	if sortOrder == nil || *sortOrder == "" {
		parsedSortOrder = defaultSortOrder
	} else {
		if !stringPresentCaseInsensitive(*sortOrder, allowedSortOrders) {
			return "", 0, errInvalidSortOrder
		}

		if strings.ToLower(*sortOrder) == "asc" {
			parsedSortOrder = 1
		} else {
			parsedSortOrder = -1
		}
	}

	return parsedSortField, parsedSortOrder, nil
}

// parseTransactionTimestampString parses the user provided timestamp string for a transaction.
// The timestamp is allowed to be empty, in which case current time will be used.
// It returns the error as the second value if any.
func parseTransactionTimestampString(timestamp string) (int64, error) {
	if timestamp == "" {
		return time.Now().Unix(), nil
	}

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return 0, errInvalidTxTimestamp
	}

	return timestampInt, nil
}

// getAllowedCategoriesForTxAmount provides the list of allowed categories of the specified transaction amount.
func getAllowedCategoriesForTxAmount(amount float64) []string {
	if amount > 0 {
		return allowedCreditCategories
	}
	return allowedDebitCategories
}

// stringPresentCaseInsensitive checks if the "value" is present in the "others" slice case-insensitively.
func stringPresentCaseInsensitive(value string, others []string) bool {
	valueLower := strings.ToLower(value)
	for _, oth := range others {
		if valueLower == strings.ToLower(oth) {
			return true
		}
	}
	return false
}

// readGetBudgetQuery reads the request query values and loads them into *getBudgetQuery type.
func readGetBudgetQuery(values url.Values) *getBudgetQuery {
	qValues := &getBudgetQuery{}

	if values.Has("start_time") {
		startTime := values.Get("start_time")
		qValues.StartTime = &startTime
	}

	if values.Has("end_time") {
		endTime := values.Get("end_time")
		qValues.EndTime = &endTime
	}

	return qValues
}

// getStartEndTimestampBudgetFilter creates a MongoDB style filter for start and end timestamps of a Budget.
func getStartEndTimestampBudgetFilter(startTime *string, endTime *string) (msi, error) {
	filter := msi{}

	if startTime != nil && *startTime != "" {
		timestamp, err := parseBudgetStartTime(*startTime)
		if err != nil {
			return nil, errInvalidBudgetTimestamp
		}
		filter["$gte"] = timestamp
	}
	if endTime != nil && *endTime != "" {
		timestamp, err := parseBudgetEndTime(*endTime)
		if err != nil {
			return nil, errInvalidBudgetTimestamp
		}
		filter["$lte"] = timestamp
	}

	return filter, nil
}

// parseBudgetStartTime parses the start time of the budget period.
// The time is allowed to be empty, in which case zero time will be used.
// It returns the error as the second value if any.
func parseBudgetStartTime(startTime string) (int64, error) {
	// If the startTime is not provided, we take it to be the beginning of time.
	if startTime == "" {
		return 0, nil
	}

	timestampInt, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		return 0, errInvalidBudgetTimestamp
	}

	return timestampInt, nil
}

// parseBudgetEndTime parses the end time of the budget period.
// The time is allowed to be empty, in which case the end of the current month will be used.
// It returns the error as the second value if any.
func parseBudgetEndTime(endTime string) (int64, error) {
	// If the endTime is not provided, we take it to be the current time.
	if endTime == "" {
		return time.Now().Unix(), nil
	}

	timestampInt, err := strconv.ParseInt(endTime, 10, 64)
	if err != nil {
		return 0, errInvalidBudgetTimestamp
	}

	return timestampInt, nil
}

// generateTransactionClosingBalanceMap generates a map of transactionID -> Closing Balance
// for a list of transactions sorted in ascending order of timestamp.
func generateTransactionClosingBalanceMap(transactions []*models.TransactionDTO) map[string]float64 {
	closingBalMap := map[string]float64{}
	accountClosingBalMap := map[string]float64{}

	for _, tx := range transactions {
		currentClosingBal4Account := accountClosingBalMap[tx.AccountID]
		newClosingBal4Account := currentClosingBal4Account + tx.Amount

		accountClosingBalMap[tx.AccountID] = newClosingBal4Account
		closingBalMap[tx.ID] = newClosingBal4Account
	}

	return closingBalMap
}

// toLastDayOfMonth returns a new date that belongs to the first moment of the last day of the month that the given
// date falls in.
func toLastDayOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
}
