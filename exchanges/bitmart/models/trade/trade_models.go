package trade

// Order represents order information
type Order struct {
	OrderID       string `json:"order_id"`
	Symbol        string `json:"symbol"`
	Side          string `json:"side"`          // buy or sell
	Type          string `json:"type"`          // limit, market, limit_maker, ioc
	Price         string `json:"price"`
	Size          string `json:"size"`
	Notional      string `json:"notional"`
	FilledSize    string `json:"filled_size"`
	FilledNotional string `json:"filled_notional"`
	Status        string `json:"status"`        // new, partially_filled, filled, canceled, pending_cancel
	CreateTime    int64  `json:"create_time"`
	UpdateTime    int64  `json:"update_time"`
}

// OrderDetail represents detailed order information
type OrderDetail struct {
	Order
	AvgPrice     string `json:"avg_price"`
	Fee          string `json:"fee"`
	FeeCurrency  string `json:"fee_currency"`
	ClientOrderID string `json:"client_order_id,omitempty"`
}

// Trade represents trade execution
type Trade struct {
	TradeID    string `json:"trade_id"`
	OrderID    string `json:"order_id"`
	Symbol     string `json:"symbol"`
	Side       string `json:"side"`
	Price      string `json:"price"`
	Size       string `json:"size"`
	Fee        string `json:"fee"`
	FeeCurrency string `json:"fee_currency"`
	ExecTime   int64  `json:"exec_time"`
}
