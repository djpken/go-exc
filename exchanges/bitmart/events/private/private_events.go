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

// FuturesAssetEvent represents futures asset balance update WebSocket event
// Channel: futures/asset:CURRENCY (e.g., futures/asset:USDT)
type FuturesAssetEvent struct {
	Group string            `json:"group"` // e.g., "futures/asset:USDT"
	Data  FuturesAssetData  `json:"data"`
}

// FuturesAssetData represents the data field in futures asset event
type FuturesAssetData struct {
	Currency         string `json:"currency"`          // Currency (e.g., "USDT", "BTC")
	AvailableBalance string `json:"available_balance"` // Available balance
	PositionDeposit  string `json:"position_deposit"`  // Position margin
	FrozenBalance    string `json:"frozen_balance"`    // Trading frozen balance
}

// FuturesPositionEvent represents futures position update WebSocket event
// Channel: futures/position
type FuturesPositionEvent struct {
	Group string                `json:"group"` // "futures/position"
	Data  []FuturesPositionData `json:"data"`
}

// FuturesPositionData represents a single position in futures position event
type FuturesPositionData struct {
	Symbol         string `json:"symbol"`           // Contract trading pair (e.g., BTCUSDT)
	HoldVolume     string `json:"hold_volume"`      // Holding quantity
	PositionType   int    `json:"position_type"`    // Position type: 1=long, 2=short
	OpenType       int    `json:"open_type"`        // Open type: 1=isolated, 2=cross
	FrozenVolume   string `json:"frozen_volume"`    // Frozen quantity
	CloseVolume    string `json:"close_volume"`     // Close quantity
	HoldAvgPrice   string `json:"hold_avg_price"`   // Holding average price
	CloseAvgPrice  string `json:"close_avg_price"`  // Close average price
	OpenAvgPrice   string `json:"open_avg_price"`   // Open average price
	LiquidatePrice string `json:"liquidate_price"`  // Liquidation price
	CreateTime     int64  `json:"create_time"`      // Position creation time (milliseconds)
	UpdateTime     int64  `json:"update_time"`      // Position update time (milliseconds)
	PositionMode   string `json:"position_mode"`    // Position mode: "hedge_mode" or "one_way_mode"
}
