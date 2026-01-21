package okex

import (
	"context"

	"github.com/djpken/go-exc/exchanges/okex/rest"
	"github.com/djpken/go-exc/exchanges/okex/ws"
	commontypes "github.com/djpken/go-exc/types"
)

// OKExExchange implements the Exchange interface for OKEx
type OKExExchange struct {
	client  *Client
	restAPI *RESTAdapter
	wsAPI   *WebSocketAdapter
	ctx     context.Context
}

// NewOKExExchange creates a new OKEx exchange instance
func NewOKExExchange(ctx context.Context, apiKey, secretKey, passphrase string, testMode bool) (*OKExExchange, error) {
	destination := NormalServer
	if testMode {
		destination = DemoServer
	}

	// Create native OKEx client
	client, err := NewClient(ctx, apiKey, secretKey, passphrase, destination)
	if err != nil {
		return nil, err
	}

	// Create adapters
	restAdapter := NewRESTAdapter(client.Rest)
	wsAdapter := NewWebSocketAdapter(client.Ws)

	return &OKExExchange{
		client:  client,
		restAPI: restAdapter,
		wsAPI:   wsAdapter,
		ctx:     ctx,
	}, nil
}

// Name returns the exchange name
func (e *OKExExchange) Name() string {
	return "OKEx"
}

// REST returns the REST API client
func (e *OKExExchange) REST() interface{} {
	return e.restAPI
}

// WebSocket returns the WebSocket API client
func (e *OKExExchange) WebSocket() interface{} {
	return e.wsAPI
}

// Close closes all connections
func (e *OKExExchange) Close() error {
	// OKEx WebSocket client doesn't have a public Close method
	// The connection management is handled internally
	return nil
}

// NativeRest returns the native REST client for advanced usage
func (e *OKExExchange) NativeRest() *rest.ClientRest {
	return e.client.Rest
}

// NativeWs returns the native WebSocket client for advanced usage
func (e *OKExExchange) NativeWs() *ws.ClientWs {
	return e.client.Ws
}

// ========== Unified API Implementation ==========
// 以下方法实现统一的跨交易所接口

// GetConfig gets account configuration
func (e *OKExExchange) GetConfig(ctx context.Context) (*commontypes.AccountConfig, error) {
	// Call native OKEx API
	resp, err := e.client.Rest.Account.GetConfig()
	if err != nil {
		return nil, err
	}

	// Check if we have configs
	if len(resp.Configs) == 0 {
		return nil, nil
	}

	// Convert to common type
	converter := NewConverter()
	return converter.ConvertAccountConfig(resp.Configs[0]), nil
}

// GetTicker gets ticker information
func (e *OKExExchange) GetTicker(ctx context.Context, symbol string) (*commontypes.Ticker, error) {
	return e.restAPI.Market().GetTicker(ctx, symbol)
}

// GetOrderBook gets order book
func (e *OKExExchange) GetOrderBook(ctx context.Context, symbol string, depth int) (*commontypes.OrderBook, error) {
	return e.restAPI.Market().GetOrderBook(ctx, symbol, depth)
}

// GetBalance gets account balance
func (e *OKExExchange) GetBalance(ctx context.Context, currencies ...string) (*commontypes.AccountBalance, error) {
	return e.restAPI.Account().GetBalance(ctx, currencies...)
}

// GetPositions gets account positions
func (e *OKExExchange) GetPositions(ctx context.Context, symbols ...string) ([]*commontypes.Position, error) {
	return e.restAPI.Account().GetPositions(ctx, symbols...)
}

// PlaceOrder places a new order
func (e *OKExExchange) PlaceOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().PlaceOrder(ctx, req.Symbol, req.Side, req.Type, req.Quantity, req.Price, req.Extra)
}

// CancelOrder cancels an existing order
func (e *OKExExchange) CancelOrder(ctx context.Context, req commontypes.CancelOrderRequest) error {
	return e.restAPI.Trade().CancelOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

// GetOrder gets order details
func (e *OKExExchange) GetOrder(ctx context.Context, req commontypes.GetOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().GetOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

// GetDepositAddress gets deposit address
func (e *OKExExchange) GetDepositAddress(ctx context.Context, currency string) (string, error) {
	return e.restAPI.Funding().GetDepositAddress(ctx, currency)
}

// Withdraw initiates a withdrawal
func (e *OKExExchange) Withdraw(ctx context.Context, req commontypes.WithdrawRequest) (string, error) {
	return e.restAPI.Funding().Withdraw(ctx, req.Currency, req.Amount, req.Address, req.Tag, req.Extra)
}
