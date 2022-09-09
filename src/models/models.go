package models

// AccountDTO is the schema of an account object as stored in the database.
type AccountDTO struct {
	// ID is the identifier of the account.
	ID string `bson:"_id" json:"id"`
	// Name is displayable name of the account.
	Name string `bson:"name" json:"name"`
}

// TransactionDTO is the schema of a transaction object as stored in the database.
type TransactionDTO struct {
	// ID is the identifier of the transaction.
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	// Amount of the transaction.
	Amount float64 `bson:"amount" json:"amount"`
	// Timestamp of the transaction.
	Timestamp int64 `bson:"timestamp" json:"timestamp"`
	// AccountID is the ID of the account to which the transaction belongs.
	AccountID string `bson:"account_id" json:"account_id"`
	// Category is one of the waterfall categories.
	Category string `bson:"category" json:"category"`
	// Tags is the list of tags on this transaction.
	Tags []string `bson:"tags" json:"tags"`
	// Notes are any details about the transaction.
	Notes string `bson:"notes" json:"notes"`

	// ClosingBal for this transaction.
	// This is calculated before returning a response, and not stored in the database.
	ClosingBal float64 `bson:"-" json:"closing_bal"`
}

// Budget is the schema of a budget object.
// A budget provides information on the planned expense and the actual expense for a period.
type Budget struct {
	TotalIncome float64 `json:"total_income"`

	EssentialsExpected float64 `json:"essentials_expected"`
	EssentialsActual   float64 `json:"essentials_actual"`

	InvestmentsExpected float64 `json:"investments_expected"`
	InvestmentsActual   float64 `json:"investments_actual"`

	SavingsExpected float64 `json:"savings_expected"`
	SavingsActual   float64 `json:"savings_actual"`

	LuxuryExpected float64 `json:"luxury_expected"`
	LuxuryActual   float64 `json:"luxury_actual"`

	IgnorableExpected float64 `json:"ignorable_expected"`
	IgnorableActual   float64 `json:"ignorable_actual"`
}
