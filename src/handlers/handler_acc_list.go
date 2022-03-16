package handlers

import (
	"fmt"
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/models"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"

	"golang.org/x/sync/errgroup"
)

// accountsListItem is an item of the account list as written by the ListAccountsHandler.
type accountsListItem struct {
	models.AccountDTO
	Balance float64 `json:"balance"`
}

// ListAccountsHandler lists all accounts along with their balances.
func ListAccountsHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	errs, errCtx := errgroup.WithContext(ctx)
	// Creating channels because we intend to make 2 database calls in parallel.
	accountsChan := make(chan []*models.AccountDTO, 1) // One call to fetch the accounts list.
	balancesChan := make(chan map[string]float64, 1)   // Second call to fetch the account balances.

	// Call 1: Fetching account list.
	errs.Go(func() error {
		defer close(accountsChan)
		accounts, err := database.ListAccounts(errCtx)
		if err != nil {
			return fmt.Errorf("failure in databae.ListAccounts: %w", err)
		}
		accountsChan <- accounts
		return nil
	})

	// Call 2: Fetching account balances.
	errs.Go(func() error {
		defer close(balancesChan)
		balances, err := database.GetAccountBalances(errCtx)
		if err != nil {
			return fmt.Errorf("failure in database.GetAccountBalances: %w", err)
		}
		balancesChan <- balances
		return nil
	})

	// Checking for errors.
	if err := errs.Wait(); err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Getting the results of calls made above.
	accounts := <-accountsChan
	balances := <-balancesChan

	// Creating the accounts list.
	accountsList := make([]*accountsListItem, len(accounts))
	for idx, acc := range accounts {
		accountsList[idx] = &accountsListItem{
			AccountDTO: *acc,
			Balance:    balances[acc.ID],
		}
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "ACCOUNTS_LISTED",
			Data:       accountsList,
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
