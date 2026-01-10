package types

// Order represents a trading order
type Order struct {
	// ID is the order ID
	ID string

	// Symbol is the trading symbol
	Symbol string

	// Side is the order side (buy/sell)
	Side string

	// Type is the order type (limit/market/etc)
	Type string

	// Status is the order status
	Status string

	// Price is the order price
	Price Decimal

	// Quantity is the order quantity
	Quantity Decimal

	// FilledQuantity is the filled quantity
	FilledQuantity Decimal

	// RemainingQuantity is the remaining quantity
	RemainingQuantity Decimal

	// Fee is the trading fee
	Fee Decimal

	// FeeCurrency is the fee currency
	FeeCurrency string

	// CreatedAt is the creation time
	CreatedAt Timestamp

	// UpdatedAt is the last update time
	UpdatedAt Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// OrderSide represents order side
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType represents order type
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusPending         OrderStatus = "pending"
	OrderStatusOpen            OrderStatus = "open"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusCanceled        OrderStatus = "canceled"
	OrderStatusRejected        OrderStatus = "rejected"
)
