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
	Instrument    = types.Instrument
	Candle        = types.Candle

	// Request types
	PlaceOrderRequest     = types.PlaceOrderRequest
	CancelOrderRequest    = types.CancelOrderRequest
	GetOrderRequest       = types.GetOrderRequest
	WithdrawRequest       = types.WithdrawRequest
	SetLeverageRequest    = types.SetLeverageRequest
	GetLeverageRequest    = types.GetLeverageRequest
	GetInstrumentsRequest = types.GetInstrumentsRequest
	GetTickersRequest     = types.GetTickersRequest
	GetCandlesRequest     = types.GetCandlesRequest
	Leverage              = types.Leverage

	// WebSocket types
	BalanceAndPositionUpdate  = types.BalanceAndPositionUpdate
	AccountUpdate             = types.AccountUpdate
	PositionUpdate            = types.PositionUpdate
	OrderUpdate               = types.OrderUpdate
	TickerUpdate              = types.TickerUpdate
	CandleUpdate              = types.CandleUpdate
	WebSocketSubscribeRequest = types.WebSocketSubscribeRequest
	WebSocketError            = types.WebSocketError
	WebSocketSubscribe        = types.WebSocketSubscribe
	WebSocketUnsubscribe      = types.WebSocketUnsubscribe
	WebSocketLogin            = types.WebSocketLogin
	WebSocketSuccess          = types.WebSocketSuccess
	WebSocketSystemMessage    = types.WebSocketSystemMessage
	WebSocketSystemError      = types.WebSocketSystemError
)

// Exchange represents a cryptocurrency exchange with a unified API interface.
// All exchanges implement the same methods, allowing you to write exchange-agnostic code.
//
// Usage:
//
//	client, _ := exc.NewExchange(ctx, exc.BitMart, config)
//	ticker, _ := client.GetTicker(ctx, "BTC_USDT")  // Works for all exchanges
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

	// GetTickers gets ticker information for all trading pairs
	// req: GetTickersRequest with optional instrument type filter
	// Returns: List of Ticker objects for all available symbols
	GetTickers(ctx context.Context, req GetTickersRequest) ([]*Ticker, error)

	// GetInstruments gets information about available trading instruments
	// req: GetInstrumentsRequest with optional instrument type filter
	// Returns: List of Instrument objects with trading pair details
	GetInstruments(ctx context.Context, req GetInstrumentsRequest) ([]*Instrument, error)

	// GetCandles gets historical candlestick/kline data
	// req: GetCandlesRequest with symbol, interval, limit, and optional time range
	// Returns: List of Candle objects with OHLCV data
	GetCandles(ctx context.Context, req GetCandlesRequest) ([]*Candle, error)

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

	// GetLeverage gets leverage configuration for trading pairs
	// req: GetLeverageRequest with symbols and margin mode
	// Returns: List of Leverage objects with current leverage settings
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	GetLeverage(ctx context.Context, req GetLeverageRequest) ([]*Leverage, error)

	// SetLeverage sets leverage for a trading pair
	// req: SetLeverageRequest with symbol/currency, leverage multiplier, margin mode
	// Returns: Updated Leverage object
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	SetLeverage(ctx context.Context, req SetLeverageRequest) (*Leverage, error)

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

	// --- WebSocket Subscriptions ---

	// SubscribeTickers subscribes to ticker updates for specified symbols via WebSocket
	// symbols: List of trading symbols to subscribe to
	// ch: Channel to receive TickerUpdate events
	// Returns: Error if subscription failed
	SubscribeTickers(ch chan *TickerUpdate, symbols ...string) error

	// UnsubscribeTickers unsubscribes from ticker updates for specified symbols
	// symbols: List of trading symbols to unsubscribe from
	// Returns: Error if unsubscription failed
	UnsubscribeTickers(symbols ...string) error

	// SubscribeCandles subscribes to candlestick/kline updates for specified symbols via WebSocket
	// ch: Channel to receive CandleUpdate events
	// interval: Candlestick interval (e.g., "1m", "5m", "1H", "1D")
	// symbols: List of trading symbols to subscribe to
	// Returns: Error if subscription failed
	SubscribeCandles(ch chan *CandleUpdate, interval string, symbols ...string) error

	// UnsubscribeCandles unsubscribes from candlestick updates for specified symbols
	// interval: Candlestick interval
	// symbols: List of trading symbols to unsubscribe from
	// Returns: Error if unsubscription failed
	UnsubscribeCandles(interval string, symbols ...string) error

	// SubscribeBalanceAndPosition subscribes to balance and position updates via WebSocket
	// ch: Channel to receive BalanceAndPositionUpdate events
	// Returns: Error if subscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	SubscribeBalanceAndPosition(ch chan *BalanceAndPositionUpdate) error

	// UnsubscribeBalanceAndPosition unsubscribes from balance and position updates
	// Returns: Error if unsubscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	UnsubscribeBalanceAndPosition() error

	// SubscribeAccount subscribes to account balance updates via WebSocket
	// currencies: Optional list of currencies to subscribe (empty = all currencies)
	// ch: Channel to receive AccountUpdate events
	// Returns: Error if subscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	SubscribeAccount(ch chan *AccountUpdate, currencies ...string) error

	// UnsubscribeAccount unsubscribes from account balance updates
	// currencies: Optional list of currencies to unsubscribe (empty = all currencies)
	// Returns: Error if unsubscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	UnsubscribeAccount(currencies ...string) error

	// SubscribePosition subscribes to position updates via WebSocket
	// ch: Channel to receive PositionUpdate events
	// req: WebSocketSubscribeRequest with subscription parameters
	// Returns: Error if subscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	SubscribePosition(ch chan *PositionUpdate, req WebSocketSubscribeRequest) error

	// UnsubscribePosition unsubscribes from position updates
	// req: WebSocketSubscribeRequest with subscription parameters
	// Returns: Error if unsubscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	UnsubscribePosition(req WebSocketSubscribeRequest) error

	// SubscribeOrders subscribes to order updates via WebSocket
	// ch: Channel to receive OrderUpdate events
	// req: WebSocketSubscribeRequest with subscription parameters
	// Returns: Error if subscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	SubscribeOrders(ch chan *OrderUpdate, req WebSocketSubscribeRequest) error

	// UnsubscribeOrders unsubscribes from order updates
	// req: WebSocketSubscribeRequest with subscription parameters
	// Returns: Error if unsubscription failed
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	UnsubscribeOrders(req WebSocketSubscribeRequest) error

	// SetChannels sets channels for receiving WebSocket events
	// errCh: Channel to receive error events
	// subCh: Channel to receive subscription events
	// unsubCh: Channel to receive unsubscription events
	// loginCh: Channel to receive login events
	// successCh: Channel to receive success events
	// systemMsgCh: Channel to receive system messages (connection, reconnection, etc.)
	// systemErrCh: Channel to receive system errors (connection failures, etc.)
	// Note: Not all exchanges support this (BitMart returns ErrNotSupported)
	SetChannels(
		errCh chan *WebSocketError,
		subCh chan *WebSocketSubscribe,
		unsubCh chan *WebSocketUnsubscribe,
		loginCh chan *WebSocketLogin,
		successCh chan *WebSocketSuccess,
		systemMsgCh chan *WebSocketSystemMessage,
		systemErrCh chan *WebSocketSystemError,
	) error
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
