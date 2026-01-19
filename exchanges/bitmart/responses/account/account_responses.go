package account

import "github.com/djpken/go-exc/exchanges/bitmart/models/account"

// BaseResponse represents the base API response structure
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

// BalanceResponse represents account balance API response
type BalanceResponse struct {
	BaseResponse
	Data struct {
		Wallet []account.Balance `json:"wallet"`
	} `json:"data"`
}

// WalletBalanceResponse represents wallet balance API response
type WalletBalanceResponse struct {
	BaseResponse
	Data account.WalletBalance `json:"data"`
}
