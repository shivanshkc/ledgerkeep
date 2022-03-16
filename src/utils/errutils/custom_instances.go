package errutils

import (
	"net/http"
)

// AccountNotFound is for requests that want to access a non-existent account.
func AccountNotFound() *HTTPError {
	return &HTTPError{StatusCode: http.StatusNotFound, CustomCode: "ACCOUNT_NOT_FOUND"}
}

// AccountAlreadyExists is for requests that want to duplicate an existing account.
func AccountAlreadyExists() *HTTPError {
	return &HTTPError{StatusCode: http.StatusConflict, CustomCode: "ACCOUNT_ALREADY_EXISTS"}
}

// AccountIsInUse is for requests that want to delete an account that is being used by transactions.
func AccountIsInUse() *HTTPError {
	return &HTTPError{StatusCode: http.StatusConflict, CustomCode: "ACCOUNT_IS_IN_USE"}
}

// TransactionNotFound is for requests that want to access a non-existent transaction.
func TransactionNotFound() *HTTPError {
	return &HTTPError{StatusCode: http.StatusNotFound, CustomCode: "TRANSACTION_NOT_FOUND"}
}
