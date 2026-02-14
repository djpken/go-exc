package types

type BarSize string

const (
	Bar1m     BarSize = "1m"
	Bar3m     BarSize = "3m"
	Bar5m     BarSize = "5m"
	Bar15m    BarSize = "15m"
	Bar30m    BarSize = "30m"
	Bar1H     BarSize = "1H"
	Bar2H     BarSize = "2H"
	Bar4H     BarSize = "4H"
	Bar6H     BarSize = "6H"
	Bar12H    BarSize = "12H"
	Bar6Hutc  BarSize = "6Hutc"
	Bar12Hutc BarSize = "12Hutc"
	Bar1Dutc  BarSize = "1Dutc"
	Bar2Dutc  BarSize = "2Dutc"
	Bar3Dutc  BarSize = "3Dutc"
	Bar1Wutc  BarSize = "1Wutc"
	Bar1Mutc  BarSize = "1Mutc"
	Bar3Mutc  BarSize = "3Mutc"
)

// Position represents a trading position
type Position struct {
	// Symbol is the trading symbol
	Symbol string

	// PosSide is the position side (long/short/net)
	// Use PositionSideLong, PositionSideShort, or PositionSideNet constants
	PosSide PositionSide

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
	Leverage int

	// MarginMode is the margin mode (cross/isolated)
	// Use MarginModeCross or MarginModeIsolated constants
	MarginMode MarginMode

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
	// Use MarginModeCross or MarginModeIsolated constants
	MarginMode MarginMode

	// PosSide is the position side (for hedge mode)
	// Use PositionSideLong, PositionSideShort constants, or empty for one-way mode
	PosSide PositionSide

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
	// Use MarginModeCross or MarginModeIsolated constants
	MarginMode MarginMode

	// PosSide is the position side for hedge mode (optional)
	// Use PositionSideLong, PositionSideShort constants, or empty for one-way mode
	PosSide PositionSide

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}

// GetLeverageRequest contains parameters for getting leverage
type GetLeverageRequest struct {
	// Symbols is the list of trading symbols
	Symbols []string

	// MarginMode is the margin mode (cross/isolated)
	// Use MarginModeCross or MarginModeIsolated constants
	MarginMode MarginMode

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}
