package bitmart

import (
	"context"

	"github.com/djpken/go-exc/exchanges/bitmart/rest"
	"github.com/djpken/go-exc/exchanges/bitmart/ws"
	commontypes "github.com/djpken/go-exc/types"
)

// BitMartExchange implements the Exchange interface for BitMart
type BitMartExchange struct {
	client  *Client
	restAPI *RESTAdapter
	wsAPI   *WebSocketAdapter
	ctx     context.Context
}

// NewBitMartExchange creates a new BitMart exchange instance
// Note: testMode parameter is ignored as BitMart doesn't have a separate test server
func NewBitMartExchange(ctx context.Context, apiKey, secretKey, memo string, testMode bool) (*BitMartExchange, error) {

	// Create native BitMart client
	client, err := NewClient(ctx, apiKey, secretKey, memo, testMode)
	if err != nil {
		return nil, err
	}

	// Create adapters
	restAdapter := NewRESTAdapter(client.Rest)
	wsAdapter := NewWebSocketAdapter(client.Ws)

	return &BitMartExchange{
		client:  client,
		restAPI: restAdapter,
		wsAPI:   wsAdapter,
		ctx:     ctx,
	}, nil
}

// Name returns the exchange name
func (e *BitMartExchange) Name() string {
	return "BitMart"
}

// REST returns the REST API client
func (e *BitMartExchange) REST() interface{} {
	return e.restAPI
}

// WebSocket returns the WebSocket API client
func (e *BitMartExchange) WebSocket() interface{} {
	return e.wsAPI
}

// Close closes all connections
func (e *BitMartExchange) Close() error {
	// Close WebSocket connection if it exists
	if e.client != nil && e.client.Ws != nil {
		return e.client.Ws.Close()
	}
	return nil
}

// GetNativeClient returns the native BitMart client for advanced usage
func (e *BitMartExchange) GetNativeClient() *Client {
	return e.client
}

// GetNativeRest returns the native REST client for advanced usage
func (e *BitMartExchange) GetNativeRest() *rest.ClientRest {
	return e.client.Rest
}

// GetNativeWs returns the native WebSocket client for advanced usage
func (e *BitMartExchange) GetNativeWs() *ws.ClientWs {
	return e.client.Ws
}

// ========== Unified API Implementation ==========
// 以下方法实现统一的跨交易所接口

// GetConfig gets account configuration
// BitMart does not support this feature
func (e *BitMartExchange) GetConfig(ctx context.Context) (*commontypes.AccountConfig, error) {
	return nil, commontypes.ErrNotSupported
}

// GetTicker gets ticker information
func (e *BitMartExchange) GetTicker(ctx context.Context, symbol string) (*commontypes.Ticker, error) {
	return e.restAPI.Market().GetTicker(ctx, symbol)
}

// GetTickers gets ticker information for all symbols
// Note: BitMart doesn't support filtering by instrument type
func (e *BitMartExchange) GetTickers(ctx context.Context, req commontypes.GetTickersRequest) ([]*commontypes.Ticker, error) {
	return e.restAPI.Market().GetTickers(ctx)
}

// GetInstruments gets trading instrument information
// Note: BitMart doesn't support filtering by instrument type
func (e *BitMartExchange) GetInstruments(ctx context.Context, req commontypes.GetInstrumentsRequest) ([]*commontypes.Instrument, error) {
	return e.restAPI.Market().GetInstruments(ctx)
}

// GetCandles gets historical candlestick/kline data
func (e *BitMartExchange) GetCandles(ctx context.Context, req commontypes.GetCandlesRequest) ([]*commontypes.Candle, error) {
	return e.restAPI.Market().GetCandles(ctx, req)
}

// GetBalance gets account balance
func (e *BitMartExchange) GetBalance(ctx context.Context, currencies ...string) (*commontypes.AccountBalance, error) {
	return e.restAPI.Account().GetBalance(ctx)
}

// GetPositions gets account positions
func (e *BitMartExchange) GetPositions(ctx context.Context, symbols ...string) ([]*commontypes.Position, error) {
	return e.restAPI.Account().GetPositions(ctx)
}

// GetLeverage gets leverage configuration
// BitMart does not support this feature
func (e *BitMartExchange) GetLeverage(ctx context.Context, req commontypes.GetLeverageRequest) ([]*commontypes.Leverage, error) {
	return nil, commontypes.ErrNotSupported
}

// SetLeverage sets leverage for a trading pair
// BitMart does not support this feature
func (e *BitMartExchange) SetLeverage(ctx context.Context, req commontypes.SetLeverageRequest) (*commontypes.Leverage, error) {
	return nil, commontypes.ErrNotSupported
}

// PlaceOrder places a new order
func (e *BitMartExchange) PlaceOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	// Add PosSide to extra parameters if provided
	// Note: BitMart spot trading doesn't use PosSide, but we keep it for consistency
	extra := req.Extra
	if extra == nil {
		extra = make(map[string]interface{})
	}
	if req.PosSide != "" {
		extra["posSide"] = req.PosSide
	}
	return e.restAPI.Trade().PlaceOrder(ctx, req.Symbol, req.Side, req.Type, req.Quantity, req.Price, extra)
}

// CancelOrder cancels an existing order
func (e *BitMartExchange) CancelOrder(ctx context.Context, req commontypes.CancelOrderRequest) error {
	return e.restAPI.Trade().CancelOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

// GetOrder gets order details
func (e *BitMartExchange) GetOrder(ctx context.Context, req commontypes.GetOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().GetOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

// ========== WebSocket Subscription Methods ==========
// BitMart WebSocket subscriptions are not supported through the unified interface
// Use the native WebSocket client directly for BitMart-specific WebSocket features

// SubscribeTickers subscribes to ticker updates for specified symbols via WebSocket
func (e *BitMartExchange) SubscribeTickers(ch chan *commontypes.TickerUpdate, symbols ...string) error {
	return e.wsAPI.SubscribeTickers(ch, symbols...)
}

// UnsubscribeTickers unsubscribes from ticker updates for specified symbols
func (e *BitMartExchange) UnsubscribeTickers(symbols ...string) error {
	return e.wsAPI.UnsubscribeTickers(symbols...)
}

// SubscribeCandles subscribes to candlestick/kline updates for specified symbols via WebSocket
func (e *BitMartExchange) SubscribeCandles(ch chan *commontypes.CandleUpdate, interval string, symbols ...string) error {
	return e.wsAPI.SubscribeCandles(ch, interval, symbols...)
}

// UnsubscribeCandles unsubscribes from candlestick updates for specified symbols
func (e *BitMartExchange) UnsubscribeCandles(interval string, symbols ...string) error {
	return e.wsAPI.UnsubscribeCandles(interval, symbols...)
}

// SubscribeBalanceAndPosition subscribes to balance and position updates via WebSocket
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) SubscribeBalanceAndPosition(ch chan *commontypes.BalanceAndPositionUpdate) error {
	return commontypes.ErrNotSupported
}

// UnsubscribeBalanceAndPosition unsubscribes from balance and position updates
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) UnsubscribeBalanceAndPosition() error {
	return commontypes.ErrNotSupported
}

// SubscribeAccount subscribes to account balance updates via WebSocket
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) SubscribeAccount(ch chan *commontypes.AccountUpdate, currencies ...string) error {
	return commontypes.ErrNotSupported
}

// UnsubscribeAccount unsubscribes from account balance updates
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) UnsubscribeAccount(currencies ...string) error {
	return commontypes.ErrNotSupported
}

// SubscribePosition subscribes to position updates via WebSocket
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) SubscribePosition(ch chan *commontypes.PositionUpdate, req commontypes.WebSocketSubscribeRequest) error {
	return commontypes.ErrNotSupported
}

// UnsubscribePosition unsubscribes from position updates
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) UnsubscribePosition(req commontypes.WebSocketSubscribeRequest) error {
	return commontypes.ErrNotSupported
}

// SubscribeOrders subscribes to order updates via WebSocket
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) SubscribeOrders(ch chan *commontypes.OrderUpdate, req commontypes.WebSocketSubscribeRequest) error {
	return commontypes.ErrNotSupported
}

// UnsubscribeOrders unsubscribes from order updates
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) UnsubscribeOrders(req commontypes.WebSocketSubscribeRequest) error {
	return commontypes.ErrNotSupported
}

// SetChannels sets channels for receiving WebSocket events
// BitMart does not support this feature through the unified interface
func (e *BitMartExchange) SetChannels(
	errCh chan *commontypes.WebSocketError,
	subCh chan *commontypes.WebSocketSubscribe,
	unsubCh chan *commontypes.WebSocketUnsubscribe,
	loginCh chan *commontypes.WebSocketLogin,
	successCh chan *commontypes.WebSocketSuccess,
	systemMsgCh chan *commontypes.WebSocketSystemMessage,
	systemErrCh chan *commontypes.WebSocketSystemError,
) error {
	return commontypes.ErrNotSupported
}
