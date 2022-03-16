package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/models"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"

	"github.com/gorilla/mux"
)

// UpdateAccountHandler updates an account by its ID.
func UpdateAccountHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	accountID := mux.Vars(request)["account_id"]
	// Validating account ID.
	if !accountIDRegexp.MatchString(accountID) {
		err := errutils.BadRequest().AddErrors(errInvalidAccountID)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Decoding the request.
	var requestBody *models.AccountDTO
	if err := httputils.UnmarshalBody(request, &requestBody); err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Validating account name.
	if !accountNameRegexp.MatchString(requestBody.Name) {
		err := errutils.BadRequest().AddErrors(errInvalidAccountName)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Database call.
	if err := database.UpdateAccount(ctx, accountID, msi{"name": requestBody.Name}); err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "ACCOUNT_UPDATED",
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
