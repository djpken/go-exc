package types

// Balance represents account balance
type Balance struct {
	// Currency is the currency code
	Currency string

	// Available is the available balance
	Available Decimal

	// Frozen is the frozen balance
	Frozen Decimal

	// Total is the total balance (Available + Frozen)
	Total Decimal

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// AccountBalance represents the complete account balance information
type AccountBalance struct {
	// Balances is the list of currency balances
	Balances []*Balance

	// TotalEquity is the total equity in base currency
	TotalEquity Decimal

	// UpdatedAt is the last update time
	UpdatedAt Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}
