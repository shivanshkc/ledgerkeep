package handlers

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	defaultLimit = 100
	defaultSkip  = 0
)

const (
	categoryEssentials  = "essentials"
	categoryInvestments = "investments"
	categorySavings     = "savings"
	categoryLuxury      = "luxury"

	categoryEarnings = "earnings"
	categoryRefunds  = "refunds"
	categoryReturns  = "returns"
	categoryPetty    = "petty"

	categoryIgnorable = "ignorable"
)

const (
	essentialsContrib  = 0.4
	investmentsContrib = 0.2
	savingsContrib     = 0.2
	luxuryContrib      = 0.2
	ignorableContrib   = 0.0
)

// TODO: Apply max and min length validations for all possible fields.
var (
	// Validation regexp(s).
	accountIDRegexp   = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	accountNameRegexp = regexp.MustCompile("^[a-zA-Z0-9-_ ]+$")

	// allowedDebitCategories are the only categories that debit transactions can have.
	allowedDebitCategories = []string{categoryEssentials, categoryInvestments, categorySavings, categoryLuxury, categoryIgnorable}
	// allowedCreditCategories are the only categories that credit transactions can have.
	allowedCreditCategories = []string{categoryEarnings, categoryRefunds, categoryReturns, categoryPetty, categoryIgnorable}

	// allowedTransactionSortFields is the list of transaction field names that can be used for sorting.
	allowedTransactionSortFields = []string{"amount", "timestamp", "category"}
	// defaultTransactionSortField is the default field by which transactions are sorted.
	defaultTransactionSortField = "timestamp"

	// allowedSortOrders are the allowed sort orders for an API.
	allowedSortOrders = []string{"asc", "desc"}
	// defaultSortOrder is the default sorting order for an API.
	defaultSortOrder = -1
)

var (
	errInvalidAccountID   = fmt.Errorf("account id should satisfy regex: %s", accountIDRegexp.String())
	errInvalidAccountName = fmt.Errorf("account name should satisfy regex: %s", accountNameRegexp.String())

	errInvalidTxID            = errors.New("transaction_id is invalid")
	errInvalidTxAmount        = errors.New("amount should be non-zero")
	errAmountCategoryMismatch = errors.New("amount not compatible with current category")
	errInvalidTxTimestamp     = fmt.Errorf("timestamp must be valid epoch seconds")
	errInvalidTxCategory      = fmt.Errorf("allowed categories for debits: %s, and for credits: %s", allowedDebitCategories, allowedCreditCategories)

	errInvalidStartAmount = errors.New("start_amount should be a float")
	errInvalidEndAmount   = errors.New("end_amount should be a float")

	errInvalidLimit = fmt.Errorf("limit should be a positive int and less than %d inclusive", defaultLimit)
	errInvalidSkip  = errors.New("skip should be a non-negative int")

	errInvalidTxSortField = fmt.Errorf("sort_field should be one of: %+v", allowedTransactionSortFields)
	errInvalidSortOrder   = fmt.Errorf("sort_order should be one of: %+v", allowedSortOrders)

	errInvalidBudgetTimestamp = fmt.Errorf("timestamp must be valid epoch seconds")

	errEmptyUpdate = errors.New("no updates provided")
)
