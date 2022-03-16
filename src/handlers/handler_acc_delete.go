package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"

	"github.com/gorilla/mux"
)

// DeleteAccountHandler deletes an account by its ID.
func DeleteAccountHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	accountID := mux.Vars(request)["account_id"]
	// Validating account ID.
	if !accountIDRegexp.MatchString(accountID) {
		err := errutils.BadRequest().AddErrors(errInvalidAccountID)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Checking if this account is in use.
	isUsed, err := database.IsAccountUsed(ctx, accountID)
	if err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// If account is in use, we cannot allow its deletion.
	if isUsed {
		httputils.WriteErrAndLog(ctx, writer, errutils.AccountIsInUse(), log)
		return
	}

	// Database call.
	if err := database.DeleteAccount(ctx, accountID); err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "ACCOUNT_DELETED",
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
