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
