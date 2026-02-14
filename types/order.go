package types

// Order represents a trading order
type Order struct {
	ClientOrderID string
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

// ===========================================
// 共用常數（Common Constants）
// ===========================================

// OrderSide 訂單方向（買/賣）
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"  // 買入
	OrderSideSell OrderSide = "sell" // 賣出
)

// OrderType 訂單類型
type OrderType string

const (
	OrderTypeLimit           OrderType = "limit"             // 限價單
	OrderTypeMarket          OrderType = "market"            // 市價單
	OrderTypeLimitMaker      OrderType = "limit_maker"       // 只做 Maker 限價單
	OrderTypeIOC             OrderType = "ioc"               // Immediate or Cancel
	OrderTypeFOK             OrderType = "fok"               // Fill or Kill
	OrderTypePostOnly        OrderType = "post_only"         // 只掛單
	OrderTypeOptimalLimitIOC OrderType = "optimal_limit_ioc" // 最優限價IOC
)

// OrderStatus 訂單狀態
type OrderStatus string

const (
	OrderStatusPending         OrderStatus = "pending"          // 待處理
	OrderStatusNew             OrderStatus = "new"              // 新訂單
	OrderStatusOpen            OrderStatus = "open"             // 開啟
	OrderStatusPartiallyFilled OrderStatus = "partially_filled" // 部分成交
	OrderStatusFilled          OrderStatus = "filled"           // 完全成交
	OrderStatusCanceled        OrderStatus = "canceled"         // 已取消
	OrderStatusCanceling       OrderStatus = "canceling"        // 取消中
	OrderStatusRejected        OrderStatus = "rejected"         // 已拒絕
	OrderStatusExpired         OrderStatus = "expired"          // 已過期
	OrderStatusLive            OrderStatus = "live"             // 活躍
)
