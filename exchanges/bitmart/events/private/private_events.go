package private

// OrderEvent represents order update WebSocket event
type OrderEvent struct {
	OrderID        string `json:"order_id"`
	ClientOrderID  string `json:"client_order_id"`
	Symbol         string `json:"symbol"`
	Side           string `json:"side"` // buy or sell
	Type           string `json:"type"` // limit, market, etc.
	Price          string `json:"price"`
	Size           string `json:"size"`
	Notional       string `json:"notional"`
	FilledSize     string `json:"filled_size"`
	FilledNotional string `json:"filled_notional"`
	Status         string `json:"status"` // new, partially_filled, filled, canceled
	CreateTime     int64  `json:"create_time"`
	UpdateTime     int64  `json:"update_time"`
}

// BalanceEvent represents balance update WebSocket event
type BalanceEvent struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Total     string `json:"total"`
	Timestamp int64  `json:"timestamp"`
}

// TradeEvent represents trade execution WebSocket event
type TradeEvent struct {
	TradeID       string `json:"trade_id"`
	OrderID       string `json:"order_id"`
	ClientOrderID string `json:"client_order_id"`
	Symbol        string `json:"symbol"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Size          string `json:"size"`
	Fee           string `json:"fee"`
	FeeCurrency   string `json:"fee_currency"`
	ExecTime      int64  `json:"exec_time"`
}
