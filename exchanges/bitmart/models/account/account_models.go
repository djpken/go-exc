package account

// Balance represents account balance information
type Balance struct {
	Currency              string `json:"currency"`                // Currency code
	Name                  string `json:"name"`                    // Full name of cryptocurrency
	Available             string `json:"available"`               // Available amount
	AvailableUsdValuation string `json:"available_usd_valuation"` // Available amount USD valuation
	Frozen                string `json:"frozen"`                  // Frozen amount in trading
	UnAvailable           string `json:"unAvailable"`             // Trading frozen + other frozen amounts
}

// WalletBalance represents wallet balance
type WalletBalance struct {
	WalletType string    `json:"wallet_type"` // spot, margin, futures
	Balances   []Balance `json:"balances"`
}
