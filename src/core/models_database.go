package core

// AccountDoc is the schema of the account document as stored in the database.
type AccountDoc struct {
	// AccountID is the unique identifier of this account.
	AccountID string `json:"account_id" bson:"_id"`
	// UserID is the identifier of the user to which this account belongs.
	UserID string `json:"user_id" bson:"user_id"`

	// AccountNumber can be used to store the bank account number.
	AccountNumber string `json:"account_number" bson:"account_number"`
	// AccountName is the easily rememberable name of this account.
	AccountName string `json:"account_name" bson:"account_name"`

	// DocCreatedAt is the time of creation of this document in the database.
	DocCreatedAt int64 `json:"doc_created_at" bson:"doc_created_at"`
	// DocUpdatedAt is the time of last update of this document is the database.
	DocUpdatedAt int64 `json:"doc_updated_at" bson:"doc_updated_at"`
}

// TransactionDoc is the schema of the transaction document as stored in the database.
type TransactionDoc struct {
	// TransactionID is the unique identifier of this transaction.
	TransactionID string `json:"transaction_id" bson:"_id"`
	// UserID is the identifier of the user to which this transaction belongs.
	UserID string `json:"user_id" bson:"user_id"`
	// AccountID is the identifier of the account to which this transaction belongs.
	AccountID string `json:"account_id" bson:"account_id"`

	// Amount of the transaction.
	Amount float64 `json:"amount" bson:"amount"`
	// AmountPerCategory holds the distribution of a transaction's amount over all categories.
	AmountPerCategory *AmountPerCategory `json:"amount_per_category" bson:"amount_per_category"`
	// Timestamp is the actual time of occurrence of this transaction.
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	// Tags of the transaction.
	Tags []string `json:"tags" bson:"tags"`
	// Notes contain any additional information about the transaction.
	Notes string `json:"notes" bson:"notes"`

	// DocCreatedAt is the time of creation of this document in the database.
	DocCreatedAt int64 `json:"doc_created_at" bson:"doc_created_at"`
	// DocUpdatedAt is the time of last update of this document is the database.
	DocUpdatedAt int64 `json:"doc_updated_at" bson:"doc_updated_at"`
}

// AmountPerCategory holds the distribution of a transaction's amount over all categories.
type AmountPerCategory struct {
	// DEBIT CATEGORIES ##############################

	// Essentials are those debits that a person cannot avoid. Example: House EMI, electricity bills, anniversaries.
	Essentials float64 `json:"essentials" bson:"essentials"`
	// Investments can be stocks, equity, real-estate, crypto etc.
	Investments float64 `json:"investments" bson:"investments"`
	// Luxury is money that is deliberately spent on comforts.
	Luxury float64 `json:"luxury" bson:"luxury"`
	// Savings are required in case of an immediate emergency.
	Savings float64 `json:"savings" bson:"savings"`
	// ###############################################

	// CREDIT CATEGORIES #############################

	// Salary is primary source of income.
	Salary float64 `json:"salary" bson:"salary"`
	// Returns can be any investment return, including bank account interest.
	Returns float64 `json:"returns" bson:"returns"`
	// Misc are all other kinds of income. Including petty credits.
	Misc float64 `json:"misc" bson:"misc"`
	// ###############################################

	// COMMON CATEGORIES #############################

	// Ignorable contains those transactions that add up to zero, and hence should not contribute to any stat.
	Ignorable float64 `json:"ignorable" bson:"ignorable"`
	// ###############################################
}
