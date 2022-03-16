package handlers

import (
	"net/http"

	"github.com/shivanshkc/ledgerkeep/src/database"
	"github.com/shivanshkc/ledgerkeep/src/logger"
	"github.com/shivanshkc/ledgerkeep/src/utils/errutils"
	"github.com/shivanshkc/ledgerkeep/src/utils/httputils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// updateTransactionBody is the schema of the body of the UpdateTransaction API.
type updateTransactionBody struct {
	Amount    *float64 `json:"amount,omitempty"`
	Timestamp *int64   `json:"timestamp,omitempty"`
	AccountID *string  `json:"account_id,omitempty"`
	Category  *string  `json:"category,omitempty"`
	Notes     *string  `json:"notes,omitempty"`
}

// UpdateTransactionHandler updates a transaction by its ID.
func UpdateTransactionHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	log := logger.Get()

	transactionIDStr := mux.Vars(request)["transaction_id"]
	// Converting the transactionID to an ObjectID. This conversion also validates the transaction ID.
	transactionID, err := primitive.ObjectIDFromHex(transactionIDStr)
	if err != nil {
		err = errutils.BadRequest().AddErrors(errInvalidTxID)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Decoding the request.
	var requestBody *updateTransactionBody
	if err := httputils.UnmarshalBody(request, &requestBody); err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Getting current transaction.
	// This has to be done before validating user input because we may require some fields
	// from the existing transaction to actually run the validation.
	currentTransaction, err := database.GetTransaction(ctx, transactionID)
	if err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Validating user input and getting the updates map.
	updates, err := prepareUpdateTransactionQuery(requestBody, currentTransaction)
	if err != nil {
		err = errutils.BadRequest().AddErrors(err)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// If the user wants to update account_id, and it is different from the current account ID...
	if newAccountID, exists := updates["account_id"]; exists && newAccountID != currentTransaction.AccountID {
		// Checking account's existence.
		accountExists, err := database.IsAccountExists(ctx, newAccountID.(string))
		if err != nil {
			httputils.WriteErrAndLog(ctx, writer, err, log)
			return
		}
		if !accountExists {
			httputils.WriteErrAndLog(ctx, writer, errutils.AccountNotFound(), log)
			return
		}
	}

	// If no updates were given, we stop execution.
	if len(updates) == 0 {
		err := errutils.BadRequest().AddErrors(errEmptyUpdate)
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Database call.
	if err := database.UpdateTransaction(ctx, transactionID, updates); err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "TRANSACTION_UPDATED",
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
