package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"
)

// createTransactionBody is the schema of the body of the CreateTransaction API.
type createTransactionBody struct {
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
	AccountID string  `json:"account_id"`
	Category  string  `json:"category"`
	Notes     string  `json:"notes"`
}

// CreateTransactionHandler creates a new transaction in the system.
func CreateTransactionHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	// Decoding the request.
	var requestBody *createTransactionBody
	if err := httputils.UnmarshalBody(request, &requestBody); err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// This call validates the user input.
	transaction, err := prepareNewTransaction(requestBody)
	if err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Checking account's existence.
	accountExists, err := database.IsAccountExists(ctx, requestBody.AccountID)
	if err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}
	if !accountExists {
		httputils.WriteErrAndLog(ctx, writer, errutils.AccountNotFound(), log)
		return
	}

	// Database call.
	insertedID, err := database.InsertTransaction(ctx, transaction)
	if err != nil {
		httputils.WriteErrAndLog(ctx, writer, errutils.AccountNotFound(), log)
		return
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusCreated,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusCreated,
			CustomCode: "TRANSACTION_CREATED",
			Data:       map[string]interface{}{"id": insertedID},
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
