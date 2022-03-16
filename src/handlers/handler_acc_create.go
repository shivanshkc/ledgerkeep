package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/models"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"
)

// CreateAccountHandler creates a new account.
func CreateAccountHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	// Decoding the request.
	var requestBody *models.AccountDTO
	if err := httputils.UnmarshalBody(request, &requestBody); err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Validating account ID.
	if !accountIDRegexp.MatchString(requestBody.ID) {
		err := errutils.BadRequest().AddErrors(errInvalidAccountID)
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
	if err := database.InsertAccount(ctx, requestBody); err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusCreated,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusCreated,
			CustomCode: "ACCOUNT_CREATED",
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
