package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"
)

// GetStatsBalancesHandler serves the info about how total balance has varied over time.
func GetStatsBalancesHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	// Database call.
	transactions, _, err := database.ListTransactions(ctx, &database.ListTransactionsParams{
		Filter:          nil,
		RequiredFields:  []string{"amount", "timestamp"},
		PaginationLimit: math.MaxInt64,
		PaginationSkip:  0,
		SortField:       "timestamp",
		SortOrder:       1,
		ExcludeCount:    true,
	})
	if err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// We will return an epoch -> balance map to the caller.
	responseBalanceMap := map[int64]float64{}
	var balanceTimestampLast *int64
	for _, tx := range transactions {
		// Getting the final timestamp of the month of the transaction.
		// This means that the balances are being grouped by months.
		balanceTimestamp := toLastDayOfMonth(time.Unix(tx.Timestamp, 0)).Unix()

		// Getting the balance of the last month if available, and adding it as the starting balance of this month.
		if balanceTimestampLast != nil && balanceTimestamp != *balanceTimestampLast {
			responseBalanceMap[balanceTimestamp] = responseBalanceMap[*balanceTimestampLast]
		}
		// Using this transaction's balance.
		responseBalanceMap[balanceTimestamp] += tx.Amount
		// Updating the last balance's timestamp for next iteration.
		balanceTimestampLast = &balanceTimestamp
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "BALANCES_FETCHED",
			Data:       responseBalanceMap,
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
