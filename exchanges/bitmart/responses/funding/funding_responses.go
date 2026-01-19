package funding

import "github.com/djpken/go-exc/exchanges/bitmart/models/funding"

// BaseResponse represents the base API response structure
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

// DepositAddressResponse represents deposit address API response
type DepositAddressResponse struct {
	BaseResponse
	Data funding.DepositAddress `json:"data"`
}

// WithdrawResponse represents withdrawal API response
type WithdrawResponse struct {
	BaseResponse
	Data struct {
		WithdrawID string `json:"withdraw_id"`
	} `json:"data"`
}

// DepositHistoryResponse represents deposit history API response
type DepositHistoryResponse struct {
	BaseResponse
	Data struct {
		Records []funding.DepositRecord `json:"records"`
	} `json:"data"`
}

// WithdrawHistoryResponse represents withdrawal history API response
type WithdrawHistoryResponse struct {
	BaseResponse
	Data struct {
		Records []funding.WithdrawRecord `json:"records"`
	} `json:"data"`
}
