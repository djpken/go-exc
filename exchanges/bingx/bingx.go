// Package bingx implements the Exchange interface for the BingX exchange.
package bingx

import (
	"context"

	"github.com/djpken/go-exc/exchanges/bingx/rest"
	"github.com/djpken/go-exc/exchanges/bingx/ws"
	commontypes "github.com/djpken/go-exc/types"
)

const (
	defaultRESTURL = "https://open-api.bingx.com"
	testRESTURL    = "https://open-api.bingx.com" // BingX simulation trading uses the same URL with demo account credentials

	defaultWSURL = "wss://open-api-swap.bingx.com/swap-market"
	testWSURL    = "wss://open-api-swap.bingx.com/swap-market" // BingX simulation trading uses the same WS URL with demo account credentials
)

// BingXExchange implements the Exchange interface for BingX
type BingXExchange struct {
	restClient *rest.ClientRest
	wsClient   *ws.ClientWs
	privateWS  *ws.PrivateClientWs
	restAPI    *RESTAdapter
	wsAPI      *WebSocketAdapter
	ctx        context.Context
	testMode   bool
}

// NewBingXExchange creates a new BingX exchange instance.
// testMode=true uses the BingX simulation trading environment (demo accounts).
func NewBingXExchange(ctx context.Context, apiKey, secretKey string, testMode bool) (*BingXExchange, error) {
	restURL := defaultRESTURL
	wsURL := defaultWSURL
	if testMode {
		restURL = testRESTURL
		wsURL = testWSURL
	}

	restClient := rest.NewClientRest(ctx, apiKey, secretKey, restURL)
	wsClient := ws.NewClientWs(wsURL, "") // public WebSocket; no auth needed

	// Private WebSocket uses listen key obtained via REST
	privateWS := ws.NewPrivateClientWs(
		restClient.CreateListenKey,
		restClient.ExtendListenKey,
		wsURL,
	)

	restAdapter := NewRESTAdapter(restClient)
	wsAdapter := NewWebSocketAdapter(wsClient, privateWS)

	return &BingXExchange{
		restClient: restClient,
		wsClient:   wsClient,
		privateWS:  privateWS,
		restAPI:    restAdapter,
		wsAPI:      wsAdapter,
		ctx:        ctx,
		testMode:   testMode,
	}, nil
}

// ─── Basic Methods ────────────────────────────────────────────────────────────

func (e *BingXExchange) Name() string {
	if e.testMode {
		return "BingX (Test)"
	}
	return "BingX"
}
func (e *BingXExchange) REST() interface{}      { return e.restAPI }
func (e *BingXExchange) WebSocket() interface{} { return e.wsAPI }

func (e *BingXExchange) Close() error {
	if e.wsClient != nil {
		_ = e.wsClient.Close()
	}
	if e.privateWS != nil {
		_ = e.privateWS.Close()
	}
	return nil
}

// ─── Market Data ─────────────────────────────────────────────────────────────

func (e *BingXExchange) GetTicker(ctx context.Context, symbol string) (*commontypes.Ticker, error) {
	return e.restAPI.Market().GetTicker(ctx, symbol)
}

func (e *BingXExchange) GetTickers(ctx context.Context, _ commontypes.GetTickersRequest) ([]*commontypes.Ticker, error) {
	return e.restAPI.Market().GetTickers(ctx)
}

func (e *BingXExchange) GetInstruments(ctx context.Context, _ commontypes.GetInstrumentsRequest) ([]*commontypes.Instrument, error) {
	return e.restAPI.Market().GetInstruments(ctx)
}

func (e *BingXExchange) GetOrderBook(ctx context.Context, symbol string, depth int) (*commontypes.OrderBook, error) {
	return e.restAPI.Market().GetOrderBook(ctx, symbol, depth)
}

func (e *BingXExchange) GetCandles(ctx context.Context, req commontypes.GetCandlesRequest) ([]*commontypes.Candle, error) {
	return e.restAPI.Market().GetCandles(ctx, req)
}

// ─── Account ─────────────────────────────────────────────────────────────────

// GetConfig is not supported by BingX
func (e *BingXExchange) GetConfig(_ context.Context) (*commontypes.AccountConfig, error) {
	return nil, commontypes.ErrNotSupported
}

func (e *BingXExchange) GetBalance(ctx context.Context, _ string, currencies ...string) (*commontypes.AccountBalance, error) {
	return e.restAPI.Account().GetBalance(ctx, currencies...)
}

func (e *BingXExchange) GetPositions(ctx context.Context, symbols ...string) ([]*commontypes.Position, error) {
	return e.restAPI.Account().GetPositions(ctx, symbols...)
}

func (e *BingXExchange) GetLeverage(ctx context.Context, req commontypes.GetLeverageRequest) ([]*commontypes.Leverage, error) {
	return e.restAPI.Account().GetLeverage(ctx, req.Symbols)
}

func (e *BingXExchange) SetLeverage(ctx context.Context, req commontypes.SetLeverageRequest) (*commontypes.Leverage, error) {
	return e.restAPI.Account().SetLeverage(ctx, req)
}

// ─── Trading ─────────────────────────────────────────────────────────────────

func (e *BingXExchange) PlaceOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().PlaceOrder(ctx, req)
}

// PlaceSingleOrder places exactly one order and returns a PlaceOrderResult.
func (e *BingXExchange) PlaceSingleOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.PlaceOrderResult, error) {
	return e.restAPI.Trade().PlaceSingleOrder(ctx, req)
}

// PlaceMultiOrder places multiple orders and returns per-order results.
func (e *BingXExchange) PlaceMultiOrder(ctx context.Context, reqs []commontypes.PlaceOrderRequest) ([]*commontypes.PlaceOrderResult, error) {
	return e.restAPI.Trade().PlaceMultiOrder(ctx, reqs)
}

func (e *BingXExchange) CancelOrder(ctx context.Context, req commontypes.CancelOrderRequest) error {
	return e.restAPI.Trade().CancelOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

func (e *BingXExchange) GetOrderDetail(ctx context.Context, req commontypes.GetOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().GetOrderDetail(ctx, req)
}

// ─── WebSocket ────────────────────────────────────────────────────────────────

func (e *BingXExchange) SubscribeTickers(ch chan *commontypes.TickerUpdate, symbols ...string) error {
	return e.wsAPI.SubscribeTickers(ch, symbols...)
}

func (e *BingXExchange) UnsubscribeTickers(symbols ...string) error {
	return e.wsAPI.UnsubscribeTickers(symbols...)
}

func (e *BingXExchange) SubscribeCandles(ch chan *commontypes.CandleUpdate, interval string, symbols ...string) error {
	return e.wsAPI.SubscribeCandles(ch, interval, symbols...)
}

func (e *BingXExchange) UnsubscribeCandles(interval string, symbols ...string) error {
	return e.wsAPI.UnsubscribeCandles(interval, symbols...)
}

// SubscribeBalanceAndPosition is not supported by BingX (use SubscribeAccount + SubscribeOrders).
func (e *BingXExchange) SubscribeBalanceAndPosition(_ chan *commontypes.BalanceAndPositionUpdate) error {
	return commontypes.ErrNotSupported
}

func (e *BingXExchange) UnsubscribeBalanceAndPosition() error {
	return commontypes.ErrNotSupported
}

// SubscribeAccount subscribes to ACCOUNT_UPDATE events via the private WebSocket.
// A listen key is obtained automatically via the REST API on first call.
func (e *BingXExchange) SubscribeAccount(ch chan *commontypes.AccountUpdate, currencies ...string) error {
	return e.wsAPI.SubscribeAccount(ch, currencies...)
}

func (e *BingXExchange) UnsubscribeAccount(currencies ...string) error {
	return e.wsAPI.UnsubscribeAccount(currencies...)
}

// SubscribePosition subscribes to position updates via the private WebSocket.
// BingX carries position changes inside ACCOUNT_UPDATE events (the "P" array).
// A listen key is obtained automatically via the REST API on first call.
func (e *BingXExchange) SubscribePosition(ch chan *commontypes.PositionUpdate, req commontypes.WebSocketSubscribeRequest) error {
	return e.wsAPI.SubscribePosition(ch, req)
}

func (e *BingXExchange) UnsubscribePosition(req commontypes.WebSocketSubscribeRequest) error {
	return e.wsAPI.UnsubscribePosition(req)
}

// SubscribeOrders subscribes to ORDER_TRADE_UPDATE events via the private WebSocket.
// A listen key is obtained automatically via the REST API on first call.
func (e *BingXExchange) SubscribeOrders(ch chan *commontypes.OrderUpdate, req commontypes.WebSocketSubscribeRequest) error {
	return e.wsAPI.SubscribeOrders(ch, req)
}

func (e *BingXExchange) UnsubscribeOrders(req commontypes.WebSocketSubscribeRequest) error {
	return e.wsAPI.UnsubscribeOrders(req)
}

func (e *BingXExchange) SetChannels(
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
