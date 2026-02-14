package contract

// GetContractDetailsRequest represents request for getting contract details
type GetContractDetailsRequest struct {
	// Symbol is the trading pair name (optional, returns all if empty)
	Symbol string `json:"symbol,omitempty"`
}

// SubmitOrderRequest represents request for placing a contract order
type SubmitOrderRequest struct {
	// Symbol is the contract trading pair (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// ClientOrderID is user-defined ID (optional)
	// Can be letters (case-sensitive) and numbers, length 1-32
	ClientOrderID string `json:"client_order_id,omitempty"`

	// Type is the order type
	// - "limit" = Limit order (default)
	// - "market" = Market order
	Type string `json:"type,omitempty"`

	// Side is the order direction
	// For dual position mode:
	// - 1 = Open long
	// - 2 = Close short
	// - 3 = Close long
	// - 4 = Open short
	// For single position mode:
	// - 1 = Buy
	// - 2 = Buy (reduce only)
	// - 3 = Sell (reduce only)
	// - 4 = Sell
	Side int `json:"side"`

	// Leverage is the leverage multiplier (optional)
	Leverage string `json:"leverage,omitempty"`

	// OpenType is the position type (optional)
	// - "cross" = Cross margin
	// - "isolated" = Isolated margin
	OpenType string `json:"open_type,omitempty"`

	// Mode is the order mode (optional)
	// - 1 = GTC (Good Till Cancel, default)
	// - 2 = FOK (Fill or Kill)
	// - 3 = IOC (Immediate or Cancel)
	// - 4 = Maker Only
	Mode int `json:"mode,omitempty"`

	// Price is the order price (required for limit orders)
	Price string `json:"price,omitempty"`

	// Size is the order quantity in contracts (required)
	Size int `json:"size"`

	// PresetTakeProfitPriceType is the preset take profit trigger price type (optional)
	// - 1 = Last price (default)
	// - 2 = Mark price
	PresetTakeProfitPriceType int `json:"preset_take_profit_price_type,omitempty"`

	// PresetStopLossPriceType is the preset stop loss trigger price type (optional)
	// - 1 = Last price (default)
	// - 2 = Mark price
	PresetStopLossPriceType int `json:"preset_stop_loss_price_type,omitempty"`

	// PresetTakeProfitPrice is the preset take profit price (optional)
	PresetTakeProfitPrice string `json:"preset_take_profit_price,omitempty"`

	// PresetStopLossPrice is the preset stop loss price (optional)
	PresetStopLossPrice string `json:"preset_stop_loss_price,omitempty"`

	// StpMode is the self-trade prevention mode (optional)
	// - 1 = cancel_maker (default)
	// - 2 = cancel_taker
	// - 3 = cancel_both
	StpMode int `json:"stp_mode,omitempty"`
}

// GetPositionV2Request represents request for getting position details V2
type GetPositionV2Request struct {
	// Symbol is the contract trading pair (optional, e.g., BTCUSDT)
	// If not provided, returns all positions
	// If provided, returns data even if no position exists
	Symbol string `json:"symbol,omitempty"`

	// Account is the trading account (optional)
	// - "futures" = Futures main account (default)
	// - "copy_trading" = Copy trading sub-account
	Account string `json:"account,omitempty"`
}

// SubmitLeverageRequest represents request for setting leverage
type SubmitLeverageRequest struct {
	// Symbol is the contract trading pair (e.g., BTCUSDT, ETHUSDT)
	Symbol string `json:"symbol"`

	// Leverage is the leverage multiplier (e.g., "5", "10", "20")
	// Optional: If not provided, uses current leverage
	Leverage string `json:"leverage,omitempty"`

	// OpenType is the margin mode (required)
	// - "cross" = Cross margin (全仓)
	// - "isolated" = Isolated margin (逐仓)
	OpenType string `json:"open_type"`
}

// GetContractTradesRequest represents request for querying contract trade details
type GetContractTradesRequest struct {
	// Symbol is the contract trading pair (optional, e.g., BTCUSDT)
	Symbol string `url:"symbol,omitempty"`

	// Account is the trading account (optional)
	// - "futures" = Main futures account
	// - "copy_trading" = Copy trading sub-account
	Account string `url:"account,omitempty"`

	// StartTime is the start timestamp in seconds (optional)
	StartTime int64 `url:"start_time,omitempty"`

	// EndTime is the end timestamp in seconds (optional)
	EndTime int64 `url:"end_time,omitempty"`

	// OrderID is the order ID to query (optional)
	OrderID string `url:"order_id,omitempty"`

	// ClientOrderID is the client order ID to query (optional)
	ClientOrderID string `url:"client_order_id,omitempty"`
}

// GetContractKlineRequest represents request for getting contract kline/candlestick data
type GetContractKlineRequest struct {
	// Symbol is the contract trading pair (required, e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Step is the kline interval in minutes (optional, default 1)
	// Valid values: 1, 3, 5, 15, 30, 60, 120, 240, 360, 720, 1440, 4320, 10080
	Step int64 `json:"step,omitempty"`

	// StartTime is the start timestamp in seconds (required)
	StartTime int64 `json:"start_time"`

	// EndTime is the end timestamp in seconds (required)
	EndTime int64 `json:"end_time"`
}
