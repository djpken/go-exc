// Package exc provides a unified interface for multiple cryptocurrency exchanges
package exc

import (
	"context"

	"github.com/djpken/go-exc/types"
)

// Import common types from types package
type (
	Order     = types.Order
	Balance   = types.AccountBalance
	Position  = types.Position
	Ticker    = types.Ticker
	OrderBook = types.OrderBook
)

// Exchange represents a cryptocurrency exchange
type Exchange interface {
	// Name returns the name of the exchange
	Name() string

	// REST returns the REST API client
	REST() RESTClient

	// WebSocket returns the WebSocket API client
	WebSocket() WebSocketClient

	// Close closes all connections and cleans up resources
	Close() error
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

// PlaceOrderRequest contains parameters for placing an order
type PlaceOrderRequest struct {
	Symbol   string
	Side     string
	Type     string
	Quantity float64
	Price    float64
	Extra    map[string]interface{}
}

// CancelOrderRequest contains parameters for canceling an order
type CancelOrderRequest struct {
	Symbol  string
	OrderID string
	Extra   map[string]interface{}
}

// GetOrderRequest contains parameters for getting order details
type GetOrderRequest struct {
	Symbol  string
	OrderID string
	Extra   map[string]interface{}
}

// WithdrawRequest contains parameters for withdrawal
type WithdrawRequest struct {
	Currency string
	Amount   float64
	Address  string
	Tag      string
	Extra    map[string]interface{}
}
