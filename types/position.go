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
	InstrumentSpot    InstrumentType = "spot"    // 現貨
	InstrumentMargin  InstrumentType = "margin"  // 槓桿
	InstrumentFutures InstrumentType = "futures" // 期貨
	InstrumentSwap    InstrumentType = "swap"    // 永續合約
	InstrumentOption  InstrumentType = "option"  // 期權
)
