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

// PositionV2 represents position details from V2 API
type PositionV2 struct {
	Symbol            string `json:"symbol"`              // Contract trading pair (e.g., BTCUSDT)
	Leverage          string `json:"leverage"`            // Leverage multiplier
	Timestamp         int64  `json:"timestamp"`           // Timestamp
	CurrentFee        string `json:"current_fee"`         // Current fee
	OpenTimestamp     int64  `json:"open_timestamp"`      // Open timestamp
	CurrentValue      string `json:"current_value"`       // Current value
	MarkPrice         string `json:"mark_price"`          // Mark price
	PositionValue     string `json:"position_value"`      // Position value
	PositionCross     string `json:"position_cross"`      // Position cross
	MaintenanceMargin string `json:"maintenance_margin"`  // Maintenance margin
	CloseVol          string `json:"close_vol"`           // Close volume
	CloseAvgPrice     string `json:"close_avg_price"`     // Close average price
	OpenAvgPrice      string `json:"open_avg_price"`      // Open average price (entry price)
	EntryPrice        string `json:"entry_price"`         // Entry price
	CurrentAmount     string `json:"current_amount"`      // Current amount (position size)
	PositionAmount    string `json:"position_amount"`     // Position amount (max position)
	RealizedValue     string `json:"realized_value"`      // Realized value
	MarkValue         string `json:"mark_value"`          // Mark value
	Account           string `json:"account"`             // Account type
	OpenType          string `json:"open_type"`           // Open type (cross/isolated)
	PositionSide      string `json:"position_side"`       // Position side (long/short/both)
	UnrealizedPnl     string `json:"unrealized_pnl"`      // Unrealized PnL
	LiquidationPrice  string `json:"liquidation_price"`   // Liquidation price
	MaxNotionalValue  string `json:"max_notional_value"`  // Max notional value
	InitialMargin     string `json:"initial_margin"`      // Initial margin
}

// ContractAsset represents a single contract asset detail
type ContractAsset struct {
	Currency         string `json:"currency"`           // Currency (USDT, BTC, ETH, etc.)
	PositionDeposit  string `json:"position_deposit"`   // Position margin
	FrozenBalance    string `json:"frozen_balance"`     // Trading frozen amount
	AvailableBalance string `json:"available_balance"`  // Available balance
	Equity           string `json:"equity"`             // Total equity
	Unrealized       string `json:"unrealized"`         // Unrealized PnL
}

// GetContractAssetsResponse represents contract assets detail API response
// API: GET /contract/private/assets-detail
type GetContractAssetsResponse struct {
	BaseResponse
	Data []ContractAsset `json:"data"`
}

// ContractTrade represents a single contract trade execution
type ContractTrade struct {
	OrderID        string `json:"order_id"`         // Order ID
	TradeID        string `json:"trade_id"`         // Trade ID
	Symbol         string `json:"symbol"`           // Contract symbol (e.g., BTCUSDT)
	Side           int    `json:"side"`             // Order direction (1=open long, 2=close short, 3=close long, 4=open short)
	Price          string `json:"price"`            // Execution price
	Vol            string `json:"vol"`              // Execution volume
	ExecType       string `json:"exec_type"`        // Liquidity type (Maker/Taker)
	Profit         bool   `json:"profit"`           // Whether profitable
	RealisedProfit string `json:"realised_profit"`  // Realized PnL
	PaidFees       string `json:"paid_fees"`        // Trading fees
	Account        string `json:"account"`          // Account type (futures/copy_trading)
	CreateTime     int64  `json:"create_time"`      // Trade creation time (ms)
}

// GetContractTradesResponse represents contract trades API response
// API: GET /contract/private/trades
type GetContractTradesResponse struct {
	BaseResponse
	Data []ContractTrade `json:"data"`
}

// GetPositionV2Response represents position details V2 API response
type GetPositionV2Response struct {
	BaseResponse
	Data []PositionV2 `json:"data"`
}

// SubmitLeverageResponse represents submit leverage API response
type SubmitLeverageResponse struct {
	BaseResponse
	Data struct {
		Symbol   string `json:"symbol"`    // Contract trading pair
		Leverage string `json:"leverage"`  // Current leverage multiplier
		OpenType string `json:"open_type"` // Margin mode (cross/isolated)
		MaxValue string `json:"max_value"` // Maximum leverage allowed
	} `json:"data"`
}

// ContractKlineData represents a single kline/candlestick data point
type ContractKlineData struct {
	Timestamp  int64  `json:"timestamp"`   // Time window (seconds)
	OpenPrice  string `json:"open_price"`  // Opening price
	ClosePrice string `json:"close_price"` // Closing price
	HighPrice  string `json:"high_price"`  // Highest price
	LowPrice   string `json:"low_price"`   // Lowest price
	Volume     string `json:"volume"`      // Trading volume
}

// GetContractKlineResponse represents contract kline API response
// API: GET /contract/public/kline
type GetContractKlineResponse struct {
	BaseResponse
	Data []ContractKlineData `json:"data"`
}
