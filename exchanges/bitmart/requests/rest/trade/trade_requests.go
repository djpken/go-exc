package trade

// PlaceOrderRequest represents request for placing an order
type PlaceOrderRequest struct {
	Symbol        string `json:"symbol"`
	Side          string `json:"side"`              // buy or sell
	Type          string `json:"type"`              // limit, market, limit_maker, ioc
	Size          string `json:"size,omitempty"`    // For limit orders
	Price         string `json:"price,omitempty"`   // For limit orders
	Notional      string `json:"notional,omitempty"` // For market buy orders
	ClientOrderID string `json:"client_order_id,omitempty"`
}

// CancelOrderRequest represents request for canceling an order
type CancelOrderRequest struct {
	Symbol  string `json:"symbol"`
	OrderID string `json:"order_id,omitempty"`
	ClientOrderID string `json:"client_order_id,omitempty"`
}

// CancelAllOrdersRequest represents request for canceling all orders
type CancelAllOrdersRequest struct {
	Symbol string `json:"symbol,omitempty"` // If empty, cancels all orders
	Side   string `json:"side,omitempty"`   // buy, sell, or empty for both
}

// GetOrderRequest represents request for getting order details
type GetOrderRequest struct {
	OrderID       string `json:"order_id,omitempty"`
	ClientOrderID string `json:"client_order_id,omitempty"`
}

// GetOrdersRequest represents request for getting order list
type GetOrdersRequest struct {
	Symbol    string `json:"symbol"`
	Status    string `json:"status,omitempty"`    // new, partially_filled, filled, canceled
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"`     // Default 50, max 200
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
}

// GetTradesRequest represents request for getting trade history
type GetTradesRequest struct {
	Symbol    string `json:"symbol"`
	OrderID   string `json:"order_id,omitempty"`
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"` // Default 50, max 200
}
