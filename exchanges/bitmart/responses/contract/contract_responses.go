package contract

import "github.com/djpken/go-exc/exchanges/bitmart/models/contract"

// BaseResponse represents the base API response structure
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

// ContractDetailsResponse represents contract details API response
type ContractDetailsResponse struct {
	BaseResponse
	Data struct {
		Symbols []contract.ContractDetail `json:"symbols"`
	} `json:"data"`
}

// SubmitOrderResponse represents submit contract order API response
type SubmitOrderResponse struct {
	BaseResponse
	Data struct {
		OrderID int64  `json:"order_id"` // Order ID
		Price   string `json:"price"`    // Order price, returns "market price" for market orders
	} `json:"data"`
}
