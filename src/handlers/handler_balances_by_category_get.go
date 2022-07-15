package handlers

import (
	"math"
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/models"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"
)

// getBalByCatQuery are the query parameters for the Get Balances by Category API.
type getBalByCatQuery struct {
	StartTime *string
	EndTime   *string
}

// GetBalancesByCategoryHandler provides the balances grouped by categories.
//
// It helps show how much expenditure is happening in each category.
func GetBalancesByCategoryHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	qValues := readGetBalByCatQuery(request.URL.Query())
	filter := msi{}

	// Validating the timestamp values and creating the timestamp filter.
	timestampFilter, err := getBudgetTimestampFilter(qValues.StartTime, qValues.EndTime)
	if err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// If the timestamp filter has any entries, we put it inside the main filter.
	if len(timestampFilter) > 0 {
		filter["timestamp"] = timestampFilter
	}

	// Database params for getting transactions that will calculate budget.
	databaseCallParams := &database.ListTransactionsParams{
		Filter:          filter,
		RequiredFields:  nil,
		PaginationLimit: math.MaxInt64,
		PaginationSkip:  0,
		SortField:       "timestamp",
		SortOrder:       1,
		ExcludeCount:    true,
	}

	// Database call.
	transactions, _, err := database.ListTransactions(ctx, databaseCallParams)
	if err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// This is the budget that will be finally returned.
	budget := &models.Budget{}

	for _, tx := range transactions {
		if tx.Amount > 0 && tx.Category != categoryIgnorable {
			budget.TotalIncome += tx.Amount
		}

		switch tx.Category {
		case categoryEssentials:
			budget.EssentialsActual += -tx.Amount
		case categoryInvestments:
			budget.InvestmentsActual += -tx.Amount
		// We do nothing for the "savings" category because it is actually the unspent amount.
		// To calculate it, we will later subtract all expenses from the total income.
		case categorySavings:
		case categoryLuxury:
			budget.LuxuryActual += -tx.Amount
		// This is a special case. That's because "ignorable" is the only common category between Credit and Debit
		// transactions.
		// We record this amount in the budget because ideally it should resolve to zero and the budget information
		// helps the user to easily keep track of it.
		case categoryIgnorable:
			budget.IgnorableActual += -tx.Amount
		default:
		}
	}

	// Calculating the SavingsActual amount by subtracting all expenses from the total income.
	budget.SavingsActual = budget.TotalIncome -
		(budget.EssentialsActual + budget.InvestmentsActual + budget.LuxuryActual + budget.IgnorableActual)

	// Calculating the expected amounts using their corresponding expected contributions.
	budget.EssentialsExpected = budget.TotalIncome * essentialsContrib
	budget.InvestmentsExpected = budget.TotalIncome * investmentsContrib
	budget.SavingsExpected = budget.TotalIncome * savingsContrib
	budget.LuxuryExpected = budget.TotalIncome * luxuryContrib
	budget.IgnorableExpected = budget.TotalIncome * ignorableContrib

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "BALANCES_FETCHED",
			Data:       budget,
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
