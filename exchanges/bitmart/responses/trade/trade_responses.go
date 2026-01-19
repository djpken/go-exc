package trade

import "github.com/djpken/go-exc/exchanges/bitmart/models/trade"

// BaseResponse represents the base API response structure
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

// PlaceOrderResponse represents place order API response
type PlaceOrderResponse struct {
	BaseResponse
	Data struct {
		OrderID       string `json:"order_id"`
		ClientOrderID string `json:"client_order_id,omitempty"`
	} `json:"data"`
}

// CancelOrderResponse represents cancel order API response
type CancelOrderResponse struct {
	BaseResponse
	Data struct {
		Result bool `json:"result"`
	} `json:"data"`
}

// CancelAllOrdersResponse represents cancel all orders API response
type CancelAllOrdersResponse struct {
	BaseResponse
	Data struct {
		Success []string `json:"success"`
		Failed  []string `json:"failed"`
	} `json:"data"`
}

// OrderResponse represents order detail API response
type OrderResponse struct {
	BaseResponse
	Data trade.OrderDetail `json:"data"`
}

// OrdersResponse represents orders list API response
type OrdersResponse struct {
	BaseResponse
	Data []trade.Order `json:"data"`
}

// TradesResponse represents trades history API response
type TradesResponse struct {
	BaseResponse
	Data []trade.Trade `json:"data"`
}
