package okex

import (
	okexconstants "github.com/djpken/go-exc/exchanges/okex/constants"
	commontypes "github.com/djpken/go-exc/types"
)

// ConstantsConverter 負責在 OKEx 特定常數和共用常數之間進行轉換
type ConstantsConverter struct{}

// NewConstantsConverter 創建新的常數轉換器
func NewConstantsConverter() *ConstantsConverter {
	return &ConstantsConverter{}
}

// ===========================================
// OrderSide 轉換
// ===========================================

// ToOKExOrderSide 將共用 OrderSide 轉換為 OKEx OrderSide
func (c *ConstantsConverter) ToOKExOrderSide(side commontypes.OrderSide) okexconstants.OrderSide {
	switch side {
	case commontypes.OrderSideBuy:
		return okexconstants.OrderBuy
	case commontypes.OrderSideSell:
		return okexconstants.OrderSell
	default:
		return okexconstants.OrderSide(side)
	}
}

// FromOKExOrderSide 將 OKEx OrderSide 轉換為共用 OrderSide
func (c *ConstantsConverter) FromOKExOrderSide(side okexconstants.OrderSide) commontypes.OrderSide {
	switch side {
	case okexconstants.OrderBuy:
		return commontypes.OrderSideBuy
	case okexconstants.OrderSell:
		return commontypes.OrderSideSell
	default:
		return commontypes.OrderSide(side)
	}
}

// ===========================================
// OrderType 轉換
// ===========================================

// ToOKExOrderType 將共用 OrderType 轉換為 OKEx OrderType
func (c *ConstantsConverter) ToOKExOrderType(orderType commontypes.OrderType) okexconstants.OrderType {
	switch orderType {
	case commontypes.OrderTypeLimit:
		return okexconstants.OrderLimit
	case commontypes.OrderTypeMarket:
		return okexconstants.OrderMarket
	case commontypes.OrderTypePostOnly:
		return okexconstants.OrderPostOnly
	case commontypes.OrderTypeFOK:
		return okexconstants.OrderFOK
	case commontypes.OrderTypeIOC:
		return okexconstants.OrderIOC
	case commontypes.OrderTypeOptimalLimitIOC:
		return okexconstants.OrderOptimalLimitIoc
	default:
		return okexconstants.OrderType(orderType)
	}
}

// FromOKExOrderType 將 OKEx OrderType 轉換為共用 OrderType
func (c *ConstantsConverter) FromOKExOrderType(orderType okexconstants.OrderType) commontypes.OrderType {
	switch orderType {
	case okexconstants.OrderLimit:
		return commontypes.OrderTypeLimit
	case okexconstants.OrderMarket:
		return commontypes.OrderTypeMarket
	case okexconstants.OrderPostOnly:
		return commontypes.OrderTypePostOnly
	case okexconstants.OrderFOK:
		return commontypes.OrderTypeFOK
	case okexconstants.OrderIOC:
		return commontypes.OrderTypeIOC
	case okexconstants.OrderOptimalLimitIoc:
		return commontypes.OrderTypeOptimalLimitIOC
	default:
		return commontypes.OrderType(orderType)
	}
}

// ===========================================
// OrderStatus 轉換
// ===========================================

// ToOKExOrderStatus 將共用 OrderStatus 轉換為 OKEx OrderState
func (c *ConstantsConverter) ToOKExOrderStatus(status commontypes.OrderStatus) okexconstants.OrderState {
	switch status {
	case commontypes.OrderStatusCanceled:
		return okexconstants.OrderCancel
	case commontypes.OrderStatusLive, commontypes.OrderStatusOpen, commontypes.OrderStatusNew:
		return okexconstants.OrderLive
	case commontypes.OrderStatusPartiallyFilled:
		return okexconstants.OrderPartiallyFilled
	case commontypes.OrderStatusFilled:
		return okexconstants.OrderFilled
	default:
		return okexconstants.OrderState(status)
	}
}

// FromOKExOrderStatus 將 OKEx OrderState 轉換為共用 OrderStatus
func (c *ConstantsConverter) FromOKExOrderStatus(state okexconstants.OrderState) commontypes.OrderStatus {
	switch state {
	case okexconstants.OrderCancel:
		return commontypes.OrderStatusCanceled
	case okexconstants.OrderLive:
		return commontypes.OrderStatusLive
	case okexconstants.OrderPartiallyFilled:
		return commontypes.OrderStatusPartiallyFilled
	case okexconstants.OrderFilled:
		return commontypes.OrderStatusFilled
	case okexconstants.OrderPause:
		return commontypes.OrderStatusCanceling
	default:
		return commontypes.OrderStatus(state)
	}
}

// ===========================================
// PositionSide 轉換
// ===========================================

// ToOKExPositionSide 將共用 PositionSide 轉換為 OKEx PositionSide
func (c *ConstantsConverter) ToOKExPositionSide(side commontypes.PositionSide) okexconstants.PositionSide {
	switch side {
	case commontypes.PositionSideLong:
		return okexconstants.PositionLongSide
	case commontypes.PositionSideShort:
		return okexconstants.PositionShortSide
	case commontypes.PositionSideNet:
		return okexconstants.PositionNetSide
	default:
		return okexconstants.PositionSide(side)
	}
}

// FromOKExPositionSide 將 OKEx PositionSide 轉換為共用 PositionSide
func (c *ConstantsConverter) FromOKExPositionSide(side okexconstants.PositionSide) commontypes.PositionSide {
	switch side {
	case okexconstants.PositionLongSide:
		return commontypes.PositionSideLong
	case okexconstants.PositionShortSide:
		return commontypes.PositionSideShort
	case okexconstants.PositionNetSide:
		return commontypes.PositionSideNet
	default:
		return commontypes.PositionSide(side)
	}
}

// ===========================================
// MarginMode 轉換
// ===========================================

// ToOKExMarginMode 將共用 MarginMode 轉換為 OKEx MarginMode
func (c *ConstantsConverter) ToOKExMarginMode(mode commontypes.MarginMode) okexconstants.MarginMode {
	switch mode {
	case commontypes.MarginModeCross:
		return okexconstants.MarginCrossMode
	case commontypes.MarginModeIsolated:
		return okexconstants.MarginIsolatedMode
	default:
		return okexconstants.MarginMode(mode)
	}
}

// FromOKExMarginMode 將 OKEx MarginMode 轉換為共用 MarginMode
func (c *ConstantsConverter) FromOKExMarginMode(mode okexconstants.MarginMode) commontypes.MarginMode {
	switch mode {
	case okexconstants.MarginCrossMode:
		return commontypes.MarginModeCross
	case okexconstants.MarginIsolatedMode:
		return commontypes.MarginModeIsolated
	default:
		return commontypes.MarginMode(mode)
	}
}

// ===========================================
// InstrumentType 轉換
// ===========================================

// ToOKExInstrumentType 將共用 InstrumentType 轉換為 OKEx InstrumentType
func (c *ConstantsConverter) ToOKExInstrumentType(instType commontypes.InstrumentType) okexconstants.InstrumentType {
	switch instType {
	case commontypes.InstrumentSpot:
		return okexconstants.SpotInstrument
	case commontypes.InstrumentMargin:
		return okexconstants.MarginInstrument
	case commontypes.InstrumentSwap:
		return okexconstants.SwapInstrument
	case commontypes.InstrumentFutures:
		return okexconstants.FuturesInstrument
	case commontypes.InstrumentOption:
		return okexconstants.OptionsInstrument
	default:
		return okexconstants.InstrumentType(instType)
	}
}

// FromOKExInstrumentType 將 OKEx InstrumentType 轉換為共用 InstrumentType
func (c *ConstantsConverter) FromOKExInstrumentType(instType okexconstants.InstrumentType) commontypes.InstrumentType {
	switch instType {
	case okexconstants.SpotInstrument:
		return commontypes.InstrumentSpot
	case okexconstants.MarginInstrument:
		return commontypes.InstrumentMargin
	case okexconstants.SwapInstrument:
		return commontypes.InstrumentSwap
	case okexconstants.FuturesInstrument:
		return commontypes.InstrumentFutures
	case okexconstants.OptionsInstrument:
		return commontypes.InstrumentOption
	default:
		return commontypes.InstrumentType(instType)
	}
}
