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

// GetTransactionHandler gets a transaction by its ID.
func GetTransactionHandler(writer http.ResponseWriter, request *http.Request) {
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

	// Database call.
	transaction, err := database.GetTransaction(ctx, transactionID)
	if err != nil {
		httputils.WriteErrAndLog(ctx, writer, err, log)
		return
	}

	// Final HTTP response.
	response := &httputils.ResponseDTO{
		Status: http.StatusOK,
		Body: &httputils.ResponseBodyDTO{
			StatusCode: http.StatusOK,
			CustomCode: "TRANSACTION_FETCHED",
			Data:       transaction,
		},
	}

	httputils.WriteAndLog(ctx, writer, response, log)
}
