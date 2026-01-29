package bitmart

import (
	"fmt"
	"time"

	publicevents "github.com/djpken/go-exc/exchanges/bitmart/events/public"
	"github.com/djpken/go-exc/exchanges/bitmart/ws"
	commontypes "github.com/djpken/go-exc/types"
)

// WebSocketAdapter adapts BitMart WebSocket client to common interface
type WebSocketAdapter struct {
	client         *ws.ClientWs
	converter      *Converter
	tickerChannels map[string]chan *commontypes.TickerUpdate              // symbol -> channel
	candleChannels map[string]map[string]chan *commontypes.CandleUpdate // interval -> symbol -> channel
}

// NewWebSocketAdapter creates a new WebSocket adapter
func NewWebSocketAdapter(client *ws.ClientWs) *WebSocketAdapter {
	return &WebSocketAdapter{
		client:         client,
		converter:      NewConverter(),
		tickerChannels: make(map[string]chan *commontypes.TickerUpdate),
		candleChannels: make(map[string]map[string]chan *commontypes.CandleUpdate),
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

	// Ensure connection
	if err := a.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Subscribe to each symbol
	for _, symbol := range symbols {
		// Create internal channel for this symbol
		internalCh := make(chan *publicevents.TickerEvent, 100)

		// Subscribe to BitMart ticker channel
		if err := a.client.Public.SubscribeTicker(symbol, internalCh); err != nil {
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
func (a *WebSocketAdapter) forwardTickerEvents(symbol string, internalCh chan *publicevents.TickerEvent, userCh chan *commontypes.TickerUpdate) {
	for event := range internalCh {
		// Convert BitMart ticker event to common TickerUpdate
		update := &commontypes.TickerUpdate{
			Symbol:           event.Symbol,
			LastPrice:        commontypes.Decimal(event.LastPrice),
			BidPrice:         commontypes.Decimal(event.BestBid),
			BidSize:          commontypes.Decimal(event.BestBidSize),
			AskPrice:         commontypes.Decimal(event.BestAsk),
			AskSize:          commontypes.Decimal(event.BestAskSize),
			High24h:          commontypes.Decimal(event.HighPrice),
			Low24h:           commontypes.Decimal(event.LowPrice),
			Open24h:          commontypes.Decimal(event.OpenPrice),
			Volume24h:        commontypes.Decimal(event.BaseVolume),
			QuoteVolume24h:   commontypes.Decimal(event.QuoteVolume),
			PriceChange24h:   commontypes.Decimal(event.PriceChange),
			PercentChange24h: commontypes.Decimal(event.PercentChange),
			Timestamp:        commontypes.Timestamp(time.Unix(0, event.Timestamp*int64(time.Millisecond))),
			Extra: map[string]interface{}{
				"close_24h": event.Close,
			},
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

	// Ensure connection
	if err := a.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
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
			Open:        commontypes.Decimal(event.Open),
			High:        commontypes.Decimal(event.High),
			Low:         commontypes.Decimal(event.Low),
			Close:       commontypes.Decimal(event.Close),
			Volume:      commontypes.Decimal(event.Volume),
			QuoteVolume: commontypes.Decimal(event.QuoteVolume),
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
