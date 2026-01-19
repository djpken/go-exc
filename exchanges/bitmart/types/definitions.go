package types

// BitMart API server destinations
type Destination string

const (
	// ProductionServer is the production API server
	ProductionServer Destination = "https://api-cloud.bitmart.com"

	// ProductionWSServer is the production WebSocket server
	ProductionWSServer Destination = "wss://ws-manager-compress.bitmart.com/api?protocol=1.1"
)

// Common BitMart constants
const (
	// APIVersion is the current API version
	APIVersion = "v1"

	// MaxRequestsPerSecond defines rate limit
	MaxRequestsPerSecond = 10
)

// OrderSide represents the order side
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderType represents the order type
type OrderType string

const (
	OrderTypeLimit      OrderType = "limit"
	OrderTypeMarket     OrderType = "market"
	OrderTypeLimitMaker OrderType = "limit_maker"
	OrderTypeIOC        OrderType = "ioc"
)

// OrderStatus represents the order status
type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "new"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusCanceled        OrderStatus = "canceled"
	OrderStatusPendingCancel   OrderStatus = "pending_cancel"
)
