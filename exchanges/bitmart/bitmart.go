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
	client, err := NewClient(ctx, apiKey, secretKey, memo)
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

// GetOrderBook gets order book
func (e *BitMartExchange) GetOrderBook(ctx context.Context, symbol string, depth int) (*commontypes.OrderBook, error) {
	return e.restAPI.Market().GetOrderBook(ctx, symbol, depth)
}

// GetBalance gets account balance
func (e *BitMartExchange) GetBalance(ctx context.Context, currencies ...string) (*commontypes.AccountBalance, error) {
	return e.restAPI.Account().GetBalance(ctx)
}

// GetPositions gets account positions
func (e *BitMartExchange) GetPositions(ctx context.Context, symbols ...string) ([]*commontypes.Position, error) {
	return e.restAPI.Account().GetPositions(ctx)
}

// PlaceOrder places a new order
func (e *BitMartExchange) PlaceOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().PlaceOrder(ctx, req.Symbol, req.Side, req.Type, req.Quantity, req.Price, req.Extra)
}

// CancelOrder cancels an existing order
func (e *BitMartExchange) CancelOrder(ctx context.Context, req commontypes.CancelOrderRequest) error {
	return e.restAPI.Trade().CancelOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

// GetOrder gets order details
func (e *BitMartExchange) GetOrder(ctx context.Context, req commontypes.GetOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().GetOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

// GetDepositAddress gets deposit address
func (e *BitMartExchange) GetDepositAddress(ctx context.Context, currency string) (string, error) {
	return e.restAPI.Funding().GetDepositAddress(ctx, currency)
}

// Withdraw initiates a withdrawal
func (e *BitMartExchange) Withdraw(ctx context.Context, req commontypes.WithdrawRequest) (string, error) {
	return e.restAPI.Funding().Withdraw(ctx, req.Currency, req.Amount, req.Address, req.Tag, req.Extra)
}
