package types

// Destination BitMart API server destinations
type Destination string

const (
	// ProductionSpotServer is the production API server
	ProductionSpotServer Destination = "https://api-cloud.bitmart.com"
	ProductionSwapServer Destination = "https://api-cloud-v2.bitmart.com"
	DemoSwapServer       Destination = "https://demo-api-cloud-v2.bitmart.com"

	// ProductionAPIWSServer is the production WebSocket server
	ProductionAPIWSServer Destination = "wss://ws-manager-compress.bitmart.com/api?protocol=1.1"
	DemoAPIWSServer       Destination = "wss://openapi-wsdemo-v2.bitmart.com/api?protocol=1.1"
	DemoUserWSServer      Destination = "wss://openapi-wsdemo-v2.bitmart.com/user?protocol=1.1"
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
