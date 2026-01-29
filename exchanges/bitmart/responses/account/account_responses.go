package account

import "github.com/djpken/go-exc/exchanges/bitmart/models/account"

// BaseResponse represents the base API response structure
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

// WalletBalanceResponse represents wallet balance API response
// API: GET /account/v1/wallet
type WalletBalanceResponse struct {
	BaseResponse
	Data struct {
		Wallet []account.Balance `json:"wallet"` // Array of wallet balances
	} `json:"data"`
}
