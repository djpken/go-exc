package types

// Position represents a trading position
type Position struct {
	// Symbol is the trading symbol
	Symbol string

	// Side is the position side (long/short)
	Side string

	// Quantity is the position quantity
	Quantity Decimal

	// AvgPrice is the average entry price
	AvgPrice Decimal

	// MarkPrice is the current mark price
	MarkPrice Decimal

	// LiquidationPrice is the liquidation price
	LiquidationPrice Decimal

	// UnrealizedPnL is the unrealized profit and loss
	UnrealizedPnL Decimal

	// RealizedPnL is the realized profit and loss
	RealizedPnL Decimal

	// Leverage is the leverage
	Leverage Decimal

	// MarginMode is the margin mode (cross/isolated)
	MarginMode string

	// CreatedAt is the creation time
	CreatedAt Timestamp

	// UpdatedAt is the last update time
	UpdatedAt Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// PositionSide represents position side
type PositionSide string

const (
	PositionSideLong  PositionSide = "long"
	PositionSideShort PositionSide = "short"
	PositionSideNet   PositionSide = "net"
)

// MarginMode represents margin mode
type MarginMode string

const (
	MarginModeCross    MarginMode = "cross"
	MarginModeIsolated MarginMode = "isolated"
)

// InstrumentType 交易產品類型
type InstrumentType string

const (
	InstrumentAny     InstrumentType = "any"     // 現貨
	InstrumentSpot    InstrumentType = "spot"    // 現貨
	InstrumentMargin  InstrumentType = "margin"  // 槓桿
	InstrumentFutures InstrumentType = "futures" // 期貨
	InstrumentSwap    InstrumentType = "swap"    // 永續合約
	InstrumentOption  InstrumentType = "option"  // 期權
)

// Leverage represents leverage configuration for a trading pair
type Leverage struct {
	// Symbol is the trading symbol
	Symbol string

	// Leverage is the leverage multiplier
	Leverage int

	// MarginMode is the margin mode (cross/isolated)
	MarginMode string

	// PositionSide is the position side (for hedge mode)
	// Empty for one-way mode
	PositionSide string

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// SetLeverageRequest contains parameters for setting leverage
type SetLeverageRequest struct {
	// Symbol is the trading symbol (optional, either Symbol or Currency must be provided)
	Symbol string

	// Currency is the margin currency (optional, either Symbol or Currency must be provided)
	Currency string

	// Leverage is the leverage multiplier
	Leverage int

	// MarginMode is the margin mode (cross/isolated)
	MarginMode string

	// PositionSide is the position side for hedge mode (optional)
	// Values: "long", "short", or empty for one-way mode
	PositionSide string

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}

// GetLeverageRequest contains parameters for getting leverage
type GetLeverageRequest struct {
	// Symbols is the list of trading symbols
	Symbols []string

	// MarginMode is the margin mode (cross/isolated)
	MarginMode string

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}
