// Package exc provides a unified interface for multiple cryptocurrency exchanges
package exc

import (
	"context"

	"github.com/djpken/go-exc/types"
)

// Import common types from types package
type (
	Order         = types.Order
	Balance       = types.AccountBalance
	Position      = types.Position
	Ticker        = types.Ticker
	OrderBook     = types.OrderBook
	AccountConfig = types.AccountConfig

	// Request types
	PlaceOrderRequest  = types.PlaceOrderRequest
	CancelOrderRequest = types.CancelOrderRequest
	GetOrderRequest    = types.GetOrderRequest
	WithdrawRequest    = types.WithdrawRequest
)

// Exchange represents a cryptocurrency exchange with a unified API interface.
// All exchanges implement the same methods, allowing you to write exchange-agnostic code.
//
// Usage:
//   client, _ := exc.NewExchange(ctx, exc.BitMart, config)
//   ticker, _ := client.GetTicker(ctx, "BTC_USDT")  // Works for all exchanges
type Exchange interface {
	// ========== Basic Methods ==========

	// Name returns the exchange name (e.g., "BitMart", "OKEx")
	Name() string

	// REST returns the exchange-specific REST API adapter
	// Returns *bitmart.RESTAdapter or *okex.RESTAdapter
	// Use this for exchange-specific features not in the unified API
	REST() interface{}

	// WebSocket returns the exchange-specific WebSocket API adapter
	// Returns *bitmart.WebSocketAdapter or *okex.WebSocketAdapter
	// Use this for exchange-specific WebSocket features
	WebSocket() interface{}

	// Close closes all connections and cleans up resources
	// Should be called when done using the exchange client
	Close() error

	// ========== Unified API Methods ==========
	// These methods work the same way across all exchanges.
	// They provide a consistent interface for common operations.

	// --- Market Data ---

	// GetTicker gets the latest ticker/price information for a symbol
	// symbol: Trading pair symbol (e.g., "BTC_USDT" for BitMart, "BTC-USDT" for OKEx)
	// Returns: Ticker with price, volume, and timestamp information
	GetTicker(ctx context.Context, symbol string) (*Ticker, error)

	// GetOrderBook gets the order book (bids and asks) for a symbol
	// symbol: Trading pair symbol
	// depth: Number of price levels to retrieve (e.g., 5, 20, 50)
	// Returns: OrderBook with bids and asks at each price level
	GetOrderBook(ctx context.Context, symbol string, depth int) (*OrderBook, error)

	// --- Account Information ---

	// GetConfig gets account configuration settings
	// Returns: AccountConfig with user ID, account level, position mode, etc.
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	GetConfig(ctx context.Context) (*AccountConfig, error)

	// GetBalance gets account balance information
	// currencies: Optional list of currencies to query (empty = all currencies)
	// Returns: AccountBalance with balances for each currency
	GetBalance(ctx context.Context, currencies ...string) (*Balance, error)

	// GetPositions gets current open positions (for futures/margin trading)
	// symbols: Optional list of symbols to query (empty = all positions)
	// Returns: List of Position objects
	// Note: Returns empty list for spot-only exchanges like BitMart
	GetPositions(ctx context.Context, symbols ...string) ([]*Position, error)

	// --- Trading Operations ---

	// PlaceOrder places a new order on the exchange
	// req: PlaceOrderRequest with symbol, side (buy/sell), type (limit/market), quantity, price
	// Returns: Order object with order ID and current status
	PlaceOrder(ctx context.Context, req PlaceOrderRequest) (*Order, error)

	// CancelOrder cancels an existing order
	// req: CancelOrderRequest with symbol and order ID
	// Returns: Error if cancellation failed
	CancelOrder(ctx context.Context, req CancelOrderRequest) error

	// GetOrder gets details of a specific order
	// req: GetOrderRequest with symbol and order ID
	// Returns: Order object with current status, filled quantity, etc.
	GetOrder(ctx context.Context, req GetOrderRequest) (*Order, error)

	// --- Funding Operations ---

	// GetDepositAddress gets the deposit address for a currency
	// currency: Currency code (e.g., "USDT", "BTC")
	// Returns: Deposit address string (may include tag/memo separated by ":")
	GetDepositAddress(ctx context.Context, currency string) (string, error)

	// Withdraw initiates a withdrawal to an external address
	// req: WithdrawRequest with currency, amount, address, and optional tag
	// Returns: Withdrawal ID string
	Withdraw(ctx context.Context, req WithdrawRequest) (string, error)
}

// RESTClient provides access to REST API endpoints
type RESTClient interface {
	// Trade returns the trading API
	Trade() TradeAPI

	// Account returns the account API
	Account() AccountAPI

	// Market returns the market data API
	Market() MarketAPI

	// Funding returns the funding API
	Funding() FundingAPI
}

// WebSocketClient provides access to WebSocket API
type WebSocketClient interface {
	// Connect establishes the WebSocket connection
	Connect() error

	// Close closes the WebSocket connection
	Close() error

	// Subscribe subscribes to a channel
	Subscribe(channel ChannelType, params SubscribeParams) error

	// SetEventHandler sets the event handler for WebSocket events
	SetEventHandler(handler EventHandler)
}

// TradeAPI provides trading operations
type TradeAPI interface {
	// PlaceOrder places a new order
	PlaceOrder(ctx context.Context, req PlaceOrderRequest) (*Order, error)

	// CancelOrder cancels an existing order
	CancelOrder(ctx context.Context, req CancelOrderRequest) error

	// GetOrder gets order details
	GetOrder(ctx context.Context, req GetOrderRequest) (*Order, error)
}

// AccountAPI provides account operations
type AccountAPI interface {
	// GetBalance gets account balance
	GetBalance(ctx context.Context) (*Balance, error)

	// GetPositions gets account positions
	GetPositions(ctx context.Context) ([]*Position, error)
}

// MarketAPI provides market data operations
type MarketAPI interface {
	// GetTicker gets ticker information
	GetTicker(ctx context.Context, symbol string) (*Ticker, error)

	// GetOrderBook gets order book
	GetOrderBook(ctx context.Context, symbol string, depth int) (*OrderBook, error)
}

// FundingAPI provides funding operations
type FundingAPI interface {
	// GetDepositAddress gets deposit address
	GetDepositAddress(ctx context.Context, currency string) (string, error)

	// Withdraw initiates a withdrawal
	Withdraw(ctx context.Context, req WithdrawRequest) (string, error)
}

// ChannelType represents a WebSocket channel type
type ChannelType string

// Common channel types
const (
	ChannelTicker    ChannelType = "ticker"
	ChannelOrderBook ChannelType = "orderbook"
	ChannelTrades    ChannelType = "trades"
	ChannelOrders    ChannelType = "orders"
	ChannelPositions ChannelType = "positions"
)

// SubscribeParams contains parameters for subscribing to a channel
type SubscribeParams map[string]interface{}

// EventHandler handles WebSocket events
type EventHandler func(event interface{})
