package bitmart

import (
	"fmt"
	"time"

	privateevents "github.com/djpken/go-exc/exchanges/bitmart/events/private"
	publicevents "github.com/djpken/go-exc/exchanges/bitmart/events/public"
	"github.com/djpken/go-exc/exchanges/bitmart/ws"
	commontypes "github.com/djpken/go-exc/types"
)

// WebSocketAdapter adapts BitMart WebSocket client to common interface
type WebSocketAdapter struct {
	client           *ws.ClientWs
	converter        *Converter
	tickerChannels   map[string]chan *commontypes.TickerUpdate            // symbol -> channel
	candleChannels   map[string]map[string]chan *commontypes.CandleUpdate // interval -> symbol -> channel
	accountChannels  map[string]chan *commontypes.AccountUpdate           // currency -> channel
	positionChannels map[string]chan *commontypes.PositionUpdate          // "default" -> channel
}

// NewWebSocketAdapter creates a new WebSocket adapter
func NewWebSocketAdapter(client *ws.ClientWs) *WebSocketAdapter {
	return &WebSocketAdapter{
		client:           client,
		converter:        NewConverter(),
		tickerChannels:   make(map[string]chan *commontypes.TickerUpdate),
		candleChannels:   make(map[string]map[string]chan *commontypes.CandleUpdate),
		accountChannels:  make(map[string]chan *commontypes.AccountUpdate),
		positionChannels: make(map[string]chan *commontypes.PositionUpdate),
	}
}

// Connect establishes the WebSocket connection
func (a *WebSocketAdapter) Connect() error {
	return a.client.Connect()
}

// Close closes the WebSocket connection
func (a *WebSocketAdapter) Close() error {
	return a.client.Close()
}

// SubscribeTickers subscribes to ticker updates for specified symbols
func (a *WebSocketAdapter) SubscribeTickers(userCh chan *commontypes.TickerUpdate, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	// Ensure connection (only connect if not already connected)
	if !a.client.IsConnected() {
		if err := a.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
	}

	// Subscribe to each symbol
	for _, symbol := range symbols {
		// Create internal channel for this symbol
		internalCh := make(chan *publicevents.FuturesTickerEvent, 100)

		// Subscribe to BitMart ticker channel
		if err := a.client.Public.SubscribeFuturesTicker(symbol, internalCh); err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", symbol, err)
		}

		// Store the user channel
		a.tickerChannels[symbol] = userCh

		// Start goroutine to convert and forward events
		go a.forwardTickerEvents(symbol, internalCh, userCh)
	}

	return nil
}

// forwardTickerEvents converts BitMart ticker events to common types and forwards them
func (a *WebSocketAdapter) forwardTickerEvents(symbol string, internalCh chan *publicevents.FuturesTickerEvent, userCh chan *commontypes.TickerUpdate) {
	for event := range internalCh {
		event := event.Data
		// Convert BitMart ticker event to common TickerUpdate
		update := &commontypes.TickerUpdate{
			Symbol:    event.Symbol,
			LastPrice: a.converter.stringToDecimal(event.LastPrice),
			BidPrice:  a.converter.stringToDecimal(event.BidPrice),
			BidSize:   a.converter.stringToDecimal(event.BidVol),
			AskPrice:  a.converter.stringToDecimal(event.AskPrice),
			AskSize:   a.converter.stringToDecimal(event.AskVol),
		}

		// Forward to user channel
		select {
		case userCh <- update:
		default:
			// Channel full, drop message
			fmt.Printf("Warning: ticker channel full for %s, dropping update\n", symbol)
		}
	}
}

// UnsubscribeTickers unsubscribes from ticker updates for specified symbols
func (a *WebSocketAdapter) UnsubscribeTickers(symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	for _, symbol := range symbols {
		// Unsubscribe from BitMart ticker channel
		if err := a.client.Public.UnsubscribeTicker(symbol); err != nil {
			return fmt.Errorf("failed to unsubscribe from %s: %w", symbol, err)
		}

		// Remove from tracking
		delete(a.tickerChannels, symbol)
	}

	return nil
}

// Subscribe subscribes to a channel
// Note: This is a simplified implementation
// Full implementation would need to map common channel types to BitMart-specific channels
func (a *WebSocketAdapter) Subscribe(channelType string, params map[string]interface{}) error {
	// TODO: Implement channel type mapping and subscription
	return nil
}

// SetEventHandler sets the event handler for WebSocket events
// Note: This is a simplified implementation
// Full implementation would need to convert BitMart events to common event types
func (a *WebSocketAdapter) SetEventHandler(handler func(event interface{})) {
	// TODO: Implement event handler wrapping and type conversion
}

// SubscribeCandles subscribes to candlestick/kline updates for specified symbols
func (a *WebSocketAdapter) SubscribeCandles(userCh chan *commontypes.CandleUpdate, interval string, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	// Ensure connection (only connect if not already connected)
	if !a.client.IsConnected() {
		if err := a.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
	}

	// Initialize interval map if needed
	if a.candleChannels[interval] == nil {
		a.candleChannels[interval] = make(map[string]chan *commontypes.CandleUpdate)
	}

	// Subscribe to each symbol
	for _, symbol := range symbols {
		// Create internal channel for this symbol
		internalCh := make(chan *publicevents.KlineEvent, 100)

		// Subscribe to BitMart kline channel
		// BitMart step format: "1m", "3m", "5m", "15m", "30m", "1H", "2H", "4H", "1D", "1W", "1M"
		if err := a.client.Public.SubscribeKline(symbol, interval, internalCh); err != nil {
			return fmt.Errorf("failed to subscribe to %s %s: %w", symbol, interval, err)
		}

		// Store the user channel
		a.candleChannels[interval][symbol] = userCh

		// Start goroutine to convert and forward events
		go a.forwardCandleEvents(symbol, interval, internalCh, userCh)
	}

	return nil
}

// forwardCandleEvents converts BitMart kline events to common types and forwards them
func (a *WebSocketAdapter) forwardCandleEvents(symbol, interval string, internalCh chan *publicevents.KlineEvent, userCh chan *commontypes.CandleUpdate) {
	for event := range internalCh {
		// Convert BitMart kline event to common CandleUpdate
		update := &commontypes.CandleUpdate{
			Symbol:      event.Symbol,
			Interval:    interval,
			Open:        a.converter.stringToDecimal(event.Open),
			High:        a.converter.stringToDecimal(event.High),
			Low:         a.converter.stringToDecimal(event.Low),
			Close:       a.converter.stringToDecimal(event.Close),
			Volume:      a.converter.stringToDecimal(event.Volume),
			QuoteVolume: a.converter.stringToDecimal(event.QuoteVolume),
			Timestamp:   commontypes.Timestamp(time.Unix(0, event.Timestamp*int64(time.Millisecond))),
			Confirmed:   false, // BitMart doesn't provide confirmation status
			Extra:       make(map[string]interface{}),
		}

		// Forward to user channel
		select {
		case userCh <- update:
		default:
			// Channel full, drop message
			fmt.Printf("Warning: candle channel full for %s %s, dropping update\n", symbol, interval)
		}
	}
}

// UnsubscribeCandles unsubscribes from candlestick updates for specified symbols
func (a *WebSocketAdapter) UnsubscribeCandles(interval string, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	for _, symbol := range symbols {
		// Unsubscribe from BitMart kline channel
		if err := a.client.Public.UnsubscribeKline(symbol, interval); err != nil {
			return fmt.Errorf("failed to unsubscribe from %s %s: %w", symbol, interval, err)
		}

		// Remove from tracking
		if a.candleChannels[interval] != nil {
			delete(a.candleChannels[interval], symbol)
		}
	}

	return nil
}

// SubscribeAccount subscribes to account/balance updates
// BitMart requires authentication before subscribing to private channels
func (a *WebSocketAdapter) SubscribeAccount(userCh chan *commontypes.AccountUpdate, currencies ...string) error {
	if len(currencies) == 0 {
		return fmt.Errorf("no currencies specified")
	}

	// Ensure connection
	if !a.client.IsConnected() {
		if err := a.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
	}

	// Authenticate if not already authenticated
	if !a.client.IsAuthenticated() {
		if err := a.client.Login(); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		// Wait a bit for authentication to complete
		time.Sleep(500 * time.Millisecond)
	}

	// Create internal channel for futures asset events
	internalCh := make(chan *privateevents.FuturesAssetEvent, 100)

	// Subscribe to BitMart futures asset channels
	if err := a.client.Private.SubscribeFuturesAsset(internalCh, currencies...); err != nil {
		return fmt.Errorf("failed to subscribe to futures asset: %w", err)
	}

	// Store the user channel for each currency
	for _, currency := range currencies {
		a.accountChannels[currency] = userCh
	}

	// Start goroutine to convert and forward events
	go a.forwardAccountEvents(internalCh, userCh)

	return nil
}

// forwardAccountEvents converts BitMart futures asset events to common types and forwards them
func (a *WebSocketAdapter) forwardAccountEvents(internalCh chan *privateevents.FuturesAssetEvent, userCh chan *commontypes.AccountUpdate) {
	for event := range internalCh {
		data := event.Data

		// Convert to common Balance type
		available := a.converter.stringToDecimal(data.AvailableBalance)
		frozen := a.converter.stringToDecimal(data.FrozenBalance)
		total, _ := available.Add(frozen)

		balance := &commontypes.Balance{
			Currency:        data.Currency,
			Available:       available,
			Frozen:          frozen,
			PositionDeposit: a.converter.stringToDecimal(data.PositionDeposit),
			Total:           total,
		}

		// Create AccountUpdate
		update := &commontypes.AccountUpdate{
			Balances:    []*commontypes.Balance{balance},
			EventType:   "update",
			UpdatedAt:   commontypes.Timestamp(time.Now()),
			TotalEquity: commontypes.ZeroDecimal,
			Extra: map[string]interface{}{
				"group": event.Group,
			},
		}

		// Forward to user channel
		select {
		case userCh <- update:
		default:
			// Channel full, drop message
			fmt.Printf("Warning: account channel full for %s, dropping update\n", data.Currency)
		}
	}
}

// UnsubscribeAccount unsubscribes from account updates
func (a *WebSocketAdapter) UnsubscribeAccount(currencies ...string) error {
	if len(currencies) == 0 {
		return fmt.Errorf("no currencies specified")
	}

	// Unsubscribe from BitMart futures asset channels
	if err := a.client.Private.UnsubscribeFuturesAsset(currencies...); err != nil {
		return fmt.Errorf("failed to unsubscribe from futures asset: %w", err)
	}

	// Remove from tracking
	for _, currency := range currencies {
		delete(a.accountChannels, currency)
	}

	return nil
}

// SubscribePosition subscribes to position updates
// BitMart requires authentication before subscribing to private channels
func (a *WebSocketAdapter) SubscribePosition(userCh chan *commontypes.PositionUpdate, req commontypes.WebSocketSubscribeRequest) error {
	// Ensure connection
	if !a.client.IsConnected() {
		if err := a.Connect(); err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
	}

	// Authenticate if not already authenticated
	if !a.client.IsAuthenticated() {
		if err := a.client.Login(); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		// Wait a bit for authentication to complete
		time.Sleep(500 * time.Millisecond)
	}

	// Create internal channel for futures position events
	internalCh := make(chan *privateevents.FuturesPositionEvent, 100)

	// Subscribe to BitMart futures position channel
	if err := a.client.Private.SubscribeFuturesPosition(internalCh); err != nil {
		return fmt.Errorf("failed to subscribe to futures position: %w", err)
	}

	// Store the user channel
	a.positionChannels["default"] = userCh

	// Start goroutine to convert and forward events
	go a.forwardPositionEvents(internalCh, userCh)

	return nil
}

// forwardPositionEvents converts BitMart futures position events to common types and forwards them
func (a *WebSocketAdapter) forwardPositionEvents(internalCh chan *privateevents.FuturesPositionEvent, userCh chan *commontypes.PositionUpdate) {
	for event := range internalCh {
		// Convert each position to common Position type
		positions := make([]*commontypes.Position, 0, len(event.Data))

		for _, data := range event.Data {
			// Convert position side
			var posSide commontypes.PositionSide
			if data.PositionMode == "hedge_mode" {
				if data.PositionType == 1 {
					posSide = commontypes.PositionSideLong
				} else {
					posSide = commontypes.PositionSideShort
				}
			} else {
				posSide = commontypes.PositionSideNet
			}

			// Convert margin mode
			var marginMode commontypes.MarginMode
			if data.OpenType == 1 {
				marginMode = commontypes.MarginModeIsolated
			} else {
				marginMode = commontypes.MarginModeCross
			}

			position := &commontypes.Position{
				Symbol:           data.Symbol,
				PosSide:          posSide,
				Quantity:         a.converter.stringToDecimal(data.HoldVolume),
				AvgPrice:         a.converter.stringToDecimal(data.OpenAvgPrice),
				MarkPrice:        commontypes.ZeroDecimal, // Not provided in BitMart position update
				LiquidationPrice: a.converter.stringToDecimal(data.LiquidatePrice),
				Leverage:         0, // Not provided in BitMart position update
				MarginMode:       marginMode,
				UnrealizedPnL:    commontypes.ZeroDecimal, // Not provided directly
				RealizedPnL:      commontypes.ZeroDecimal, // Not provided directly
				CreatedAt:        commontypes.Timestamp(time.UnixMilli(data.CreateTime)),
				UpdatedAt:        commontypes.Timestamp(time.UnixMilli(data.UpdateTime)),
				Extra: map[string]interface{}{
					"hold_avg_price":  data.HoldAvgPrice,
					"close_avg_price": data.CloseAvgPrice,
					"frozen_volume":   data.FrozenVolume,
					"close_volume":    data.CloseVolume,
					"position_mode":   data.PositionMode,
				},
			}

			positions = append(positions, position)
		}

		// Create PositionUpdate
		update := &commontypes.PositionUpdate{
			Positions: positions,
			EventType: "update",
			UpdatedAt: commontypes.Timestamp(time.Now()),
			Extra: map[string]interface{}{
				"group": event.Group,
			},
		}

		// Forward to user channel
		select {
		case userCh <- update:
		default:
			// Channel full, drop message
			fmt.Printf("Warning: position channel full, dropping update\n")
		}
	}
}

// UnsubscribePosition unsubscribes from position updates
func (a *WebSocketAdapter) UnsubscribePosition(req commontypes.WebSocketSubscribeRequest) error {
	// Unsubscribe from BitMart futures position channel
	if err := a.client.Private.UnsubscribeFuturesPosition(); err != nil {
		return fmt.Errorf("failed to unsubscribe from futures position: %w", err)
	}

	// Remove from tracking
	delete(a.positionChannels, "default")

	return nil
}

// SetChannels sets channels for receiving WebSocket events
// This allows users to receive notifications about connection events, errors, subscriptions, etc.
func (a *WebSocketAdapter) SetChannels(
	errCh chan *commontypes.WebSocketError,
	subCh chan *commontypes.WebSocketSubscribe,
	unsubCh chan *commontypes.WebSocketUnsubscribe,
	loginCh chan *commontypes.WebSocketLogin,
	successCh chan *commontypes.WebSocketSuccess,
	systemMsgCh chan *commontypes.WebSocketSystemMessage,
	systemErrCh chan *commontypes.WebSocketSystemError,
) error {
	// Delegate to the underlying WebSocket client
	a.client.SetChannels(errCh, subCh, unsubCh, loginCh, successCh, systemMsgCh, systemErrCh)
	return nil
}
