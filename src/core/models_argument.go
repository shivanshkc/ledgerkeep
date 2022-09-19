package core

// ParamsAccountCreate are the params required to create an account.
type ParamsAccountCreate struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string

	// AccountName for the account.
	AccountName string
	// AccountNumber for the account.
	AccountNumber string
}

// ParamsAccountUpdate are the params required to update an account.
type ParamsAccountUpdate struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string
	// AccountID is the identifier of the account that is to be updated.
	AccountID string

	// AccountName is the new value of AccountName. It is optional.
	AccountName *string
	// AccountNumber is the new value of AccountNumber. It is optional.
	AccountNumber *string
}

// ParamsAccountDelete are params required to delete an account.
type ParamsAccountDelete struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string
	// AccountID is the identifier of the account that is to be deleted.
	AccountID string
}

// ParamsAccountList are the params required to list accounts.
type ParamsAccountList struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string

	// Search through the AccountName and AccountNumber fields. It follows the "similar" search approach.
	Search *string

	PaginationLimit *int64
	PaginationSkip  *int64

	// SortField is the name of the field that will be used for sorting the results.
	SortField *string
	// SortOrder is the order of sorting. asc or desc.
	SortOrder *string
}

// ParamsTransactionCreate are the params required to create a transaction.
type ParamsTransactionCreate struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string
	// AccountID is the identifier of the account to which this transaction will be mapped.
	AccountID string

	// Amount of the transaction.
	Amount float64
	// AmountPerCategory for the transaction.
	AmountPerCategory *AmountPerCategory
	// Timestamp is the time of occurrence for this transaction.
	Timestamp *int64
	// Tags of the transaction.
	Tags []string
	// Notes are any additional details about the transaction.
	Notes string
}

// ParamsTransactionUpdate are the params required to update a transaction.
type ParamsTransactionUpdate struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string
	// TransactionID is the identifier of the transaction to be updated.
	TransactionID string

	// AccountID will be the new accountID of the transaction. It is optional.
	AccountID *string
	// Amount will be the new amount of the transaction. It is optional.
	Amount *float64
	// AmountPerCategory will be the new amount distribution for the transaction.
	AmountPerCategory *AmountPerCategory
	// Timestamp will be the new timestamp of the transaction. It is optional.
	Timestamp *int64
	// Tags will be the new tags of the transaction. It is optional.
	Tags []string
	// Notes will be the new notes of the transaction. It is optional.
	Notes *string
}

// ParamsTransactionDelete are the params required to delete a transaction.
type ParamsTransactionDelete struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string
	// TransactionID is the identifier of the transaction to be deleted.
	TransactionID string
}

// ParamsTransactionGet are the params required to fetch a single transaction.
type ParamsTransactionGet struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string
	// TransactionID is the identifier of the transaction to be fetched.
	TransactionID string
}

// ParamsTransactionList are the params required to list transactions.
type ParamsTransactionList struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string

	// AccountID can be used to filter based on the account of the transaction.
	AccountID *string

	// StartAmount can be used to filter out transactions with a lower amount than provided.
	StartAmount *float64
	// EndAmount can be used to filter out transactions with a higher amount than provided.
	EndAmount *float64

	// StartTime can be used to filter out transactions with an earlier timestamp than provided.
	StartTime *int64
	// EndTime can be used to filter out transactions with a later timestamp than provided.
	EndTime *int64

	// Category can be used to filter based on the category of the transaction.
	Category *string

	// Tags can be used to filter based on the tags of the transaction.
	// The nature of the query ("all" or "any") can be switched using the TagsAny flag.
	Tags []string
	// TagsAny tells whether the tags provided in the Tags field should be treated as "any" or "all".
	TagsAny *bool

	// NotesHint provides full-text search on the Notes field.
	NotesHint *string

	PaginationLimit *int64
	PaginationSkip  *int64

	// SortField is the name of the field that will be used for sorting the results.
	SortField *string
	// SortOrder is the order of sorting. asc or desc.
	SortOrder *string
}

// ParamsStatsGetAmountDistribution are the params required to get the debit or credit distribution by category or tags.
type ParamsStatsGetAmountDistribution struct {
	// UserID is the identifier of the user who is performing this operation.
	UserID string

	// StartTime can be used to only consider transactions that occurred after the specified time.
	StartTime *int64
	// EndTime can be used to only consider transactions that occurred before the specified time.
	EndTime *int64
}
