package okex

import (
	"context"
	"strconv"

	okexEvents "github.com/djpken/go-exc/exchanges/okex/events"
	privateEvents "github.com/djpken/go-exc/exchanges/okex/events/private"
	privateWs "github.com/djpken/go-exc/exchanges/okex/requests/ws/private"
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

// GetTickers gets ticker information for multiple instruments
func (e *OKExExchange) GetTickers(ctx context.Context, req commontypes.GetTickersRequest) ([]*commontypes.Ticker, error) {
	return e.restAPI.Market().GetTickers(ctx, req)
}

// GetInstruments gets trading instrument information
func (e *OKExExchange) GetInstruments(ctx context.Context, req commontypes.GetInstrumentsRequest) ([]*commontypes.Instrument, error) {
	return e.restAPI.Market().GetInstruments(ctx, req)
}

// GetCandles gets historical candlestick/kline data
func (e *OKExExchange) GetCandles(ctx context.Context, req commontypes.GetCandlesRequest) ([]*commontypes.Candle, error) {
	return e.restAPI.Market().GetCandles(ctx, req)
}

// GetBalance gets account balance
func (e *OKExExchange) GetBalance(ctx context.Context, typee string, currencies ...string) (*commontypes.AccountBalance, error) {
	return e.restAPI.Account().GetBalance(ctx, currencies...)
}

// GetPositions gets account positions
func (e *OKExExchange) GetPositions(ctx context.Context, symbols ...string) ([]*commontypes.Position, error) {
	return e.restAPI.Account().GetPositions(ctx, symbols...)
}

// GetLeverage gets leverage configuration
func (e *OKExExchange) GetLeverage(ctx context.Context, req commontypes.GetLeverageRequest) ([]*commontypes.Leverage, error) {
	return e.restAPI.Account().GetLeverage(ctx, req)
}

// SetLeverage sets leverage for a trading pair
func (e *OKExExchange) SetLeverage(ctx context.Context, req commontypes.SetLeverageRequest) (*commontypes.Leverage, error) {
	return e.restAPI.Account().SetLeverage(ctx, req)
}

// PlaceOrder places a new order
func (e *OKExExchange) PlaceOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().PlaceOrder(ctx, req)
}

// CancelOrder cancels an existing order
func (e *OKExExchange) CancelOrder(ctx context.Context, req commontypes.CancelOrderRequest) error {
	return e.restAPI.Trade().CancelOrder(ctx, req.Symbol, req.OrderID, req.Extra)
}

// GetOrderDetail gets order details
func (e *OKExExchange) GetOrderDetail(ctx context.Context, req commontypes.GetOrderRequest) (*commontypes.Order, error) {
	return e.restAPI.Trade().GetOrderDetail(ctx, req)
}

// ========== WebSocket Subscription Methods ==========

// SubscribeTickers subscribes to ticker updates for specified symbols via WebSocket
func (e *OKExExchange) SubscribeTickers(ch chan *commontypes.TickerUpdate, symbols ...string) error {
	return e.wsAPI.SubscribeTickers(ch, symbols...)
}

// UnsubscribeTickers unsubscribes from ticker updates for specified symbols
func (e *OKExExchange) UnsubscribeTickers(symbols ...string) error {
	return e.wsAPI.UnsubscribeTickers(symbols...)
}

// SubscribeCandles subscribes to candlestick/kline updates for specified symbols via WebSocket
func (e *OKExExchange) SubscribeCandles(ch chan *commontypes.CandleUpdate, interval string, symbols ...string) error {
	return e.wsAPI.SubscribeCandles(ch, interval, symbols...)
}

// UnsubscribeCandles unsubscribes from candlestick updates for specified symbols
func (e *OKExExchange) UnsubscribeCandles(interval string, symbols ...string) error {
	return e.wsAPI.UnsubscribeCandles(interval, symbols...)
}

// SubscribeBalanceAndPosition subscribes to balance and position updates via WebSocket
func (e *OKExExchange) SubscribeBalanceAndPosition(ch chan *commontypes.BalanceAndPositionUpdate) error {
	// Create internal channel for native OKEx events
	nativeCh := make(chan *privateEvents.BalanceAndPosition, 100)

	// Start goroutine to convert events
	go func() {
		converter := NewConverter()
		for event := range nativeCh {
			if len(event.BalanceAndPositions) > 0 {
				if converted := converter.ConvertBalanceAndPosition(event.BalanceAndPositions[0]); converted != nil {
					ch <- converted
				}
			}
		}
	}()

	// Subscribe using native client
	return e.client.Ws.Private.BalanceAndPosition(nativeCh)
}

// UnsubscribeBalanceAndPosition unsubscribes from balance and position updates
func (e *OKExExchange) UnsubscribeBalanceAndPosition() error {
	return e.client.Ws.Private.UBalanceAndPosition(true)
}

// SubscribeAccount subscribes to account balance updates via WebSocket
func (e *OKExExchange) SubscribeAccount(ch chan *commontypes.AccountUpdate, currencies ...string) error {
	// Create internal channel for native OKEx events
	nativeCh := make(chan *privateEvents.Account, 100)

	// Start goroutine to handle pagination and convert events
	go func() {
		converter := NewConverter()
		// Store pages data for merging
		pageData := make(map[int]*commontypes.AccountUpdate)
		var lastPageNum int

		for event := range nativeCh {
			// Convert current page
			converted := converter.ConvertAccountEvent(event.Balances, string(event.EventType))
			if converted == nil {
				continue
			}

			// For non-snapshot events, send immediately without pagination handling
			if event.EventType != "snapshot" {
				ch <- converted
				continue
			}

			// Store this page's data
			pageData[event.CurPage] = converted

			// If this is the last page, record it
			if event.LastPage {
				lastPageNum = event.CurPage
			}

			// Check if we have all pages (from 1 to lastPageNum)
			if lastPageNum > 0 {
				allPagesReceived := true
				for i := 1; i <= lastPageNum; i++ {
					if _, exists := pageData[i]; !exists {
						allPagesReceived = false
						break
					}
				}

				// If all pages received, merge and send
				if allPagesReceived {
					// Merge all balances from all pages
					mergedBalances := make([]*commontypes.Balance, 0)
					var totalEquity commontypes.Decimal
					var updateTime commontypes.Timestamp

					for page := 1; page <= lastPageNum; page++ {
						pageUpdate := pageData[page]
						mergedBalances = append(mergedBalances, pageUpdate.Balances...)
						if page == lastPageNum {
							totalEquity = pageUpdate.TotalEquity
							updateTime = pageUpdate.UpdatedAt
						}
					}

					// Send merged result
					ch <- &commontypes.AccountUpdate{
						Balances:    mergedBalances,
						EventType:   string(event.EventType),
						UpdatedAt:   updateTime,
						TotalEquity: totalEquity,
						Extra:       map[string]interface{}{},
					}

					// Clear page data for next snapshot
					pageData = make(map[int]*commontypes.AccountUpdate)
					lastPageNum = 0
				}
			}
		}
	}()

	// Prepare request
	req := privateWs.Account{}
	if len(currencies) > 0 {
		req.Ccy = currencies[0] // OKEx supports single currency
	}

	// Subscribe using native client
	return e.client.Ws.Private.Account(req, nativeCh)
}

// UnsubscribeAccount unsubscribes from account balance updates
func (e *OKExExchange) UnsubscribeAccount(currencies ...string) error {
	req := privateWs.Account{}
	if len(currencies) > 0 {
		req.Ccy = currencies[0]
	}
	return e.client.Ws.Private.UAccount(req, true)
}

// SubscribePosition subscribes to position updates via WebSocket
func (e *OKExExchange) SubscribePosition(ch chan *commontypes.PositionUpdate, req commontypes.WebSocketSubscribeRequest) error {
	// Create internal channel for native OKEx events
	nativeCh := make(chan *privateEvents.Position, 100)

	// Start goroutine to handle pagination and convert events
	go func() {
		converter := NewConverter()
		// Store pages data for merging
		pageData := make(map[int]*commontypes.PositionUpdate)
		var lastPageNum int

		for event := range nativeCh {
			// Convert current page
			converted := converter.ConvertPositionEvent(event.Positions, string(event.EventType))
			if converted == nil {
				continue
			}

			// For non-snapshot events, send immediately without pagination handling
			if event.EventType != "snapshot" {
				ch <- converted
				continue
			}

			// Store this page's data
			pageData[event.CurPage] = converted

			// If this is the last page, record it
			if event.LastPage {
				lastPageNum = event.CurPage
			}

			// Check if we have all pages (from 1 to lastPageNum)
			if lastPageNum > 0 {
				allPagesReceived := true
				for i := 1; i <= lastPageNum; i++ {
					if _, exists := pageData[i]; !exists {
						allPagesReceived = false
						break
					}
				}

				// If all pages received, merge and send
				if allPagesReceived {
					// Merge all positions from all pages
					mergedPositions := make([]*commontypes.Position, 0)
					var updateTime commontypes.Timestamp

					for page := 1; page <= lastPageNum; page++ {
						pageUpdate := pageData[page]
						mergedPositions = append(mergedPositions, pageUpdate.Positions...)
						if page == lastPageNum {
							updateTime = pageUpdate.UpdatedAt
						}
					}

					// Send merged result
					ch <- &commontypes.PositionUpdate{
						Positions: mergedPositions,
						EventType: string(event.EventType),
						UpdatedAt: updateTime,
						Extra:     map[string]interface{}{},
					}

					// Clear page data for next snapshot
					pageData = make(map[int]*commontypes.PositionUpdate)
					lastPageNum = 0
				}
			}
		}
	}()

	// Prepare request
	converter := NewConverter()
	posReq := privateWs.Position{
		InstType: converter.ConvertInstrumentType(req.InstrumentType),
	}
	if len(req.Symbols) > 0 {
		posReq.InstID = req.Symbols[0] // OKEx supports single symbol
	}

	// Subscribe using native client
	return e.client.Ws.Private.Position(posReq, nativeCh)
}

// UnsubscribePosition unsubscribes from position updates
func (e *OKExExchange) UnsubscribePosition(req commontypes.WebSocketSubscribeRequest) error {
	converter := NewConverter()
	posReq := privateWs.Position{
		InstType: converter.ConvertInstrumentType(req.InstrumentType),
	}
	if len(req.Symbols) > 0 {
		posReq.InstID = req.Symbols[0]
	}
	return e.client.Ws.Private.UPosition(posReq, true)
}

// SubscribeOrders subscribes to order updates via WebSocket
func (e *OKExExchange) SubscribeOrders(ch chan *commontypes.OrderUpdate, req commontypes.WebSocketSubscribeRequest) error {
	// Create internal channel for native OKEx events
	nativeCh := make(chan *privateEvents.Order, 100)

	// Start goroutine to convert events
	go func() {
		converter := NewConverter()
		for event := range nativeCh {
			if converted := converter.ConvertOrderEvent(event.Orders); converted != nil {
				ch <- converted
			}
		}
	}()

	// Prepare request
	converter := NewConverter()
	orderReq := privateWs.Order{
		InstType: converter.ConvertInstrumentType(req.InstrumentType),
	}
	if len(req.Symbols) > 0 {
		orderReq.InstID = req.Symbols[0] // OKEx supports single symbol
	}

	// Subscribe using native client
	return e.client.Ws.Private.Order(orderReq, nativeCh)
}

// UnsubscribeOrders unsubscribes from order updates
func (e *OKExExchange) UnsubscribeOrders(req commontypes.WebSocketSubscribeRequest) error {
	converter := NewConverter()
	orderReq := privateWs.Order{
		InstType: converter.ConvertInstrumentType(req.InstrumentType),
	}
	if len(req.Symbols) > 0 {
		orderReq.InstID = req.Symbols[0]
	}
	return e.client.Ws.Private.UOrder(orderReq, true)
}

// SetChannels sets channels for receiving WebSocket events
func (e *OKExExchange) SetChannels(
	errCh chan *commontypes.WebSocketError,
	subCh chan *commontypes.WebSocketSubscribe,
	unsubCh chan *commontypes.WebSocketUnsubscribe,
	loginCh chan *commontypes.WebSocketLogin,
	successCh chan *commontypes.WebSocketSuccess,
	systemMsgCh chan *commontypes.WebSocketSystemMessage,
	systemErrCh chan *commontypes.WebSocketSystemError,
) error {
	// Create native OKEx event channels
	var nativeErrCh chan *okexEvents.Error
	var nativeSubCh chan *okexEvents.Subscribe
	var nativeUnsubCh chan *okexEvents.Unsubscribe
	var nativeLoginCh chan *okexEvents.Login
	var nativeSuccessCh chan *okexEvents.Success

	// Setup error channel and converter
	if errCh != nil {
		nativeErrCh = make(chan *okexEvents.Error, 100)
		go func() {
			for event := range nativeErrCh {
				errCh <- &commontypes.WebSocketError{
					Event:   event.Event,
					Code:    strconv.FormatInt(int64(event.Code), 10),
					Message: event.Msg,
					Extra: map[string]interface{}{
						"op":   event.Op,
						"args": event.Args,
						"arg":  event.Arg,
						"data": event.Data,
						"id":   event.ID,
					},
				}
			}
		}()
	}

	// Setup subscribe channel and converter
	if subCh != nil {
		nativeSubCh = make(chan *okexEvents.Subscribe, 100)
		go func() {
			for event := range nativeSubCh {
				channel := ""
				if event.Arg != nil {
					if ch, ok := event.Arg.Get("channel"); ok {
						channel = ch.(string)
					}
				}
				subCh <- &commontypes.WebSocketSubscribe{
					Event:   event.Event,
					Channel: channel,
					Extra: map[string]interface{}{
						"arg": event.Arg,
					},
				}
			}
		}()
	}

	// Setup unsubscribe channel and converter
	if unsubCh != nil {
		nativeUnsubCh = make(chan *okexEvents.Unsubscribe, 100)
		go func() {
			for event := range nativeUnsubCh {
				channel := ""
				if event.Arg != nil {
					if ch, ok := event.Arg.Get("channel"); ok {
						channel = ch.(string)
					}
				}
				unsubCh <- &commontypes.WebSocketUnsubscribe{
					Event:   event.Event,
					Channel: channel,
					Extra: map[string]interface{}{
						"arg": event.Arg,
					},
				}
			}
		}()
	}

	// Setup login channel and converter
	if loginCh != nil {
		nativeLoginCh = make(chan *okexEvents.Login, 100)
		go func() {
			for event := range nativeLoginCh {
				loginCh <- &commontypes.WebSocketLogin{
					Event:   event.Event,
					Code:    event.Code,
					Message: event.Msg,
					Extra:   map[string]interface{}{},
				}
			}
		}()
	}

	// Setup success channel and converter
	if successCh != nil {
		nativeSuccessCh = make(chan *okexEvents.Success, 100)
		go func() {
			for event := range nativeSuccessCh {
				successCh <- &commontypes.WebSocketSuccess{
					Code:    event.Code,
					Message: event.Msg,
					Extra: map[string]interface{}{
						"id":   event.ID,
						"op":   event.Op,
						"data": event.Data,
					},
				}
			}
		}()
	}

	// Setup system message channel and converter
	var nativeSystemMsgCh chan *ws.SystemMessage
	if systemMsgCh != nil {
		nativeSystemMsgCh = make(chan *ws.SystemMessage, 100)
		go func() {
			for event := range nativeSystemMsgCh {
				systemMsgCh <- &commontypes.WebSocketSystemMessage{
					Type:      event.Type,
					Message:   event.Message,
					Private:   event.Private,
					Timestamp: commontypes.Timestamp(event.Timestamp),
					Extra:     map[string]interface{}{},
				}
			}
		}()
	}

	// Setup system error channel and converter
	var nativeSystemErrCh chan *ws.SystemError
	if systemErrCh != nil {
		nativeSystemErrCh = make(chan *ws.SystemError, 100)
		go func() {
			for event := range nativeSystemErrCh {
				errMsg := ""
				if event.Error != nil {
					errMsg = event.Error.Error()
				}
				systemErrCh <- &commontypes.WebSocketSystemError{
					Type:      event.Type,
					Error:     errMsg,
					Private:   event.Private,
					Timestamp: commontypes.Timestamp(event.Timestamp),
					Extra:     map[string]interface{}{},
				}
			}
		}()
	}

	// Set channels on native WebSocket client
	e.client.Ws.SetChannels(nativeErrCh, nativeSubCh, nativeUnsubCh, nativeLoginCh, nativeSuccessCh)
	e.client.Ws.SetSystemChannels(nativeSystemMsgCh, nativeSystemErrCh)

	return nil
}
