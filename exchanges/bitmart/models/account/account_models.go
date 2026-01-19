package account

// Balance represents account balance information
type Balance struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Total     string `json:"total"`
}

// WalletBalance represents wallet balance
type WalletBalance struct {
	WalletType string    `json:"wallet_type"` // spot, margin, futures
	Balances   []Balance `json:"balances"`
}
