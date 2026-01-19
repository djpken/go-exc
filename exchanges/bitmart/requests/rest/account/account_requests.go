package account

// GetBalanceRequest represents request for getting account balance
type GetBalanceRequest struct {
	Currency string `json:"currency,omitempty"` // Optional: specific currency
}

// GetWalletBalanceRequest represents request for getting wallet balance
type GetWalletBalanceRequest struct {
	WalletType string `json:"wallet_type,omitempty"` // spot, margin, futures
}
