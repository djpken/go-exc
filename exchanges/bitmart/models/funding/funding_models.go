package funding

// DepositAddress represents deposit address information
type DepositAddress struct {
	Currency       string `json:"currency"`
	Chain          string `json:"chain"`
	Address        string `json:"address"`
	Tag            string `json:"tag,omitempty"`
	AddressTag     string `json:"address_tag,omitempty"`
}

// WithdrawRecord represents withdrawal record
type WithdrawRecord struct {
	WithdrawID string `json:"withdraw_id"`
	Currency   string `json:"currency"`
	Chain      string `json:"chain"`
	Amount     string `json:"amount"`
	Fee        string `json:"fee"`
	Status     string `json:"status"` // processing, completed, failed, canceled
	Address    string `json:"address"`
	Tag        string `json:"tag,omitempty"`
	TxID       string `json:"tx_id,omitempty"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

// DepositRecord represents deposit record
type DepositRecord struct {
	DepositID  string `json:"deposit_id"`
	Currency   string `json:"currency"`
	Chain      string `json:"chain"`
	Amount     string `json:"amount"`
	Status     string `json:"status"` // pending, completed, failed
	Address    string `json:"address"`
	Tag        string `json:"tag,omitempty"`
	TxID       string `json:"tx_id"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
