package account

// GetBalanceRequest represents request for getting account balance
type GetBalanceRequest struct {
	Currency string `json:"currency,omitempty"` // Optional: specific currency
}

// GetWalletBalanceRequest represents request for getting wallet balance
type GetWalletBalanceRequest struct {
	Currency           string `json:"currency,omitempty"`             // Optional: specific currency, e.g., BTC
	NeedUsdValuation   bool   `json:"need_usd_valuation,omitempty"`   // Optional: whether to return USD valuation, default false
}
