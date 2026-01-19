package funding

// GetDepositAddressRequest represents request for getting deposit address
type GetDepositAddressRequest struct {
	Currency string `json:"currency"`
	Chain    string `json:"chain,omitempty"` // Optional: blockchain network
}

// WithdrawRequest represents request for withdrawal
type WithdrawRequest struct {
	Currency   string `json:"currency"`
	Amount     string `json:"amount"`
	Address    string `json:"address"`
	AddressTag string `json:"address_tag,omitempty"` // For currencies like XRP, EOS
	Chain      string `json:"chain,omitempty"`       // Blockchain network
}

// GetDepositHistoryRequest represents request for getting deposit history
type GetDepositHistoryRequest struct {
	Currency  string `json:"currency,omitempty"`
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"` // Default 50, max 200
}

// GetWithdrawHistoryRequest represents request for getting withdrawal history
type GetWithdrawHistoryRequest struct {
	Currency  string `json:"currency,omitempty"`
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"` // Default 50, max 200
}
