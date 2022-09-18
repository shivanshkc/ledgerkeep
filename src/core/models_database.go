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
	// Timestamp is the actual time of occurrence of this transaction.
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	// Category of the transaction.
	Category string `json:"category" bson:"category"`
	// Tags of the transaction.
	Tags []string `json:"tags" bson:"tags"`
	// Notes contain any additional information about the transaction.
	Notes string `json:"notes" bson:"notes"`

	// DocCreatedAt is the time of creation of this document in the database.
	DocCreatedAt int64 `json:"doc_created_at" bson:"doc_created_at"`
	// DocUpdatedAt is the time of last update of this document is the database.
	DocUpdatedAt int64 `json:"doc_updated_at" bson:"doc_updated_at"`
}
