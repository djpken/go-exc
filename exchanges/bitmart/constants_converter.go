package bitmart

import (
	bitmarttypes "github.com/djpken/go-exc/exchanges/bitmart/types"
	commontypes "github.com/djpken/go-exc/types"
)

// ConstantsConverter 負責在 BitMart 特定常數和共用常數之間進行轉換
type ConstantsConverter struct{}

// NewConstantsConverter 創建新的常數轉換器
func NewConstantsConverter() *ConstantsConverter {
	return &ConstantsConverter{}
}

// ===========================================
// OrderSide 轉換
// ===========================================

// ToBitMartOrderSide 將共用 OrderSide 轉換為 BitMart OrderSide
func (c *ConstantsConverter) ToBitMartOrderSide(side commontypes.OrderSide) bitmarttypes.OrderSide {
	switch side {
	case commontypes.OrderSideBuy:
		return bitmarttypes.OrderSideBuy
	case commontypes.OrderSideSell:
		return bitmarttypes.OrderSideSell
	default:
		return bitmarttypes.OrderSide(side)
	}
}

// FromBitMartOrderSide 將 BitMart OrderSide 轉換為共用 OrderSide
func (c *ConstantsConverter) FromBitMartOrderSide(side bitmarttypes.OrderSide) commontypes.OrderSide {
	switch side {
	case bitmarttypes.OrderSideBuy:
		return commontypes.OrderSideBuy
	case bitmarttypes.OrderSideSell:
		return commontypes.OrderSideSell
	default:
		return commontypes.OrderSide(side)
	}
}

// ===========================================
// OrderType 轉換
// ===========================================

// ToBitMartOrderType 將共用 OrderType 轉換為 BitMart OrderType
func (c *ConstantsConverter) ToBitMartOrderType(orderType commontypes.OrderType) bitmarttypes.OrderType {
	switch orderType {
	case commontypes.OrderTypeLimit:
		return bitmarttypes.OrderTypeLimit
	case commontypes.OrderTypeMarket:
		return bitmarttypes.OrderTypeMarket
	case commontypes.OrderTypeLimitMaker:
		return bitmarttypes.OrderTypeLimitMaker
	case commontypes.OrderTypeIOC:
		return bitmarttypes.OrderTypeIOC
	default:
		// BitMart 不支持其他類型，默認使用 limit
		return bitmarttypes.OrderTypeLimit
	}
}

// FromBitMartOrderType 將 BitMart OrderType 轉換為共用 OrderType
func (c *ConstantsConverter) FromBitMartOrderType(orderType bitmarttypes.OrderType) commontypes.OrderType {
	switch orderType {
	case bitmarttypes.OrderTypeLimit:
		return commontypes.OrderTypeLimit
	case bitmarttypes.OrderTypeMarket:
		return commontypes.OrderTypeMarket
	case bitmarttypes.OrderTypeLimitMaker:
		return commontypes.OrderTypeLimitMaker
	case bitmarttypes.OrderTypeIOC:
		return commontypes.OrderTypeIOC
	default:
		return commontypes.OrderType(orderType)
	}
}

// ===========================================
// OrderStatus 轉換
// ===========================================

// ToBitMartOrderStatus 將共用 OrderStatus 轉換為 BitMart OrderStatus
func (c *ConstantsConverter) ToBitMartOrderStatus(status commontypes.OrderStatus) bitmarttypes.OrderStatus {
	switch status {
	case commontypes.OrderStatusNew, commontypes.OrderStatusOpen:
		return bitmarttypes.OrderStatusNew
	case commontypes.OrderStatusPartiallyFilled:
		return bitmarttypes.OrderStatusPartiallyFilled
	case commontypes.OrderStatusFilled:
		return bitmarttypes.OrderStatusFilled
	case commontypes.OrderStatusCanceled:
		return bitmarttypes.OrderStatusCanceled
	case commontypes.OrderStatusCanceling:
		return bitmarttypes.OrderStatusPendingCancel
	default:
		return bitmarttypes.OrderStatus(status)
	}
}

// FromBitMartOrderStatus 將 BitMart OrderStatus 轉換為共用 OrderStatus
func (c *ConstantsConverter) FromBitMartOrderStatus(status bitmarttypes.OrderStatus) commontypes.OrderStatus {
	switch status {
	case bitmarttypes.OrderStatusNew:
		return commontypes.OrderStatusNew
	case bitmarttypes.OrderStatusPartiallyFilled:
		return commontypes.OrderStatusPartiallyFilled
	case bitmarttypes.OrderStatusFilled:
		return commontypes.OrderStatusFilled
	case bitmarttypes.OrderStatusCanceled:
		return commontypes.OrderStatusCanceled
	case bitmarttypes.OrderStatusPendingCancel:
		return commontypes.OrderStatusCanceling
	default:
		return commontypes.OrderStatus(status)
	}
}

// ===========================================
// String 轉換（用於簡單的字符串到常數轉換）
// ===========================================

// OrderSideToString 將共用 OrderSide 轉換為字符串
func (c *ConstantsConverter) OrderSideToString(side commontypes.OrderSide) string {
	return string(side)
}

// StringToOrderSide 將字符串轉換為共用 OrderSide
func (c *ConstantsConverter) StringToOrderSide(side string) commontypes.OrderSide {
	return commontypes.OrderSide(side)
}

// OrderTypeToString 將共用 OrderType 轉換為字符串
func (c *ConstantsConverter) OrderTypeToString(orderType commontypes.OrderType) string {
	return string(orderType)
}

// StringToOrderType 將字符串轉換為共用 OrderType
func (c *ConstantsConverter) StringToOrderType(orderType string) commontypes.OrderType {
	return commontypes.OrderType(orderType)
}

// OrderStatusToString 將共用 OrderStatus 轉換為字符串
func (c *ConstantsConverter) OrderStatusToString(status commontypes.OrderStatus) string {
	return string(status)
}

// StringToOrderStatus 將字符串轉換為共用 OrderStatus
func (c *ConstantsConverter) StringToOrderStatus(status string) commontypes.OrderStatus {
	return commontypes.OrderStatus(status)
}
