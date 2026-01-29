package okex

import (
	"fmt"

	okexconstants "github.com/djpken/go-exc/exchanges/okex/constants"
	publicevents "github.com/djpken/go-exc/exchanges/okex/events/public"
	publicrequests "github.com/djpken/go-exc/exchanges/okex/requests/ws/public"
	"github.com/djpken/go-exc/exchanges/okex/ws"
	commontypes "github.com/djpken/go-exc/types"
)

// WebSocketAdapter adapts OKEx WebSocket client to common interface
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
	// OKEx WebSocket connects automatically
	// No explicit connect needed
	return nil
}

// Close closes the WebSocket connection
func (a *WebSocketAdapter) Close() error {
	// OKEx WebSocket client doesn't have a public Close method
	// The connection management is handled internally
	return nil
}

// Subscribe subscribes to a channel
func (a *WebSocketAdapter) Subscribe(channel string, params map[string]interface{}) error {
	// TODO: Implement channel subscription mapping
	// This requires mapping common channel types to OKEx specific channels
	return nil
}

// SetEventHandler sets the event handler for WebSocket events
func (a *WebSocketAdapter) SetEventHandler(handler func(event interface{})) {
	// TODO: Implement event handler mapping
	// This requires wrapping OKEx events into common event format
}

// GetNativeClient returns the native WebSocket client for advanced usage
func (a *WebSocketAdapter) GetNativeClient() *ws.ClientWs {
	return a.client
}

// Public returns the public channel API
func (a *WebSocketAdapter) Public() *ws.Public {
	return a.client.Public
}

// Private returns the private channel API
func (a *WebSocketAdapter) Private() *ws.Private {
	return a.client.Private
}

// Trade returns the trade API
func (a *WebSocketAdapter) Trade() *ws.Trade {
	return a.client.Trade
}

// SubscribeTickers subscribes to ticker updates for specified symbols
func (a *WebSocketAdapter) SubscribeTickers(userCh chan *commontypes.TickerUpdate, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	// Subscribe to each symbol
	for _, symbol := range symbols {
		// Create internal channel for this symbol's ticker events
		internalCh := make(chan *publicevents.Tickers, 100)

		// Subscribe to OKEx tickers channel
		req := publicrequests.Tickers{
			InstID: symbol,
		}
		if err := a.client.Public.Tickers(req, internalCh); err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", symbol, err)
		}

		// Store the user channel
		a.tickerChannels[symbol] = userCh

		// Start goroutine to convert and forward events
		go a.forwardTickerEvents(symbol, internalCh, userCh)
	}

	return nil
}

// forwardTickerEvents converts OKEx ticker events to common types and forwards them
func (a *WebSocketAdapter) forwardTickerEvents(symbol string, internalCh chan *publicevents.Tickers, userCh chan *commontypes.TickerUpdate) {
	for event := range internalCh {
		// OKEx Tickers event contains multiple ticker updates
		for _, ticker := range event.Tickers {
			// Convert OKEx ticker to common TickerUpdate
			update := a.converter.ConvertTickerToUpdate(ticker)
			if update == nil {
				continue
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
}

// UnsubscribeTickers unsubscribes from ticker updates for specified symbols
func (a *WebSocketAdapter) UnsubscribeTickers(symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	for _, symbol := range symbols {
		// Unsubscribe from OKEx tickers channel
		req := publicrequests.Tickers{
			InstID: symbol,
		}
		if err := a.client.Public.UTickers(req, true); err != nil {
			return fmt.Errorf("failed to unsubscribe from %s: %w", symbol, err)
		}

		// Remove from tracking
		delete(a.tickerChannels, symbol)
	}

	return nil
}

// SubscribeCandles subscribes to candlestick updates for specified symbols
func (a *WebSocketAdapter) SubscribeCandles(userCh chan *commontypes.CandleUpdate, interval string, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	// Convert common interval format to OKEx format
	// Common: "1m", "5m", "1H", "1D"
	// OKEx:   "candle1m", "candle5m", "candle1H", "candle1D"
	okexInterval := "candle" + interval

	// Initialize interval map if needed
	if a.candleChannels[interval] == nil {
		a.candleChannels[interval] = make(map[string]chan *commontypes.CandleUpdate)
	}

	// Subscribe to each symbol
	for _, symbol := range symbols {
		// Create internal channel for this symbol's candle events
		internalCh := make(chan *publicevents.Candlesticks, 100)

		// Subscribe to OKEx candlesticks channel
		req := publicrequests.Candlesticks{
			InstID:  symbol,
			Channel: okexconstants.CandleStickWsBarSize(okexInterval),
		}
		if err := a.client.Public.Candlesticks(req, internalCh); err != nil {
			return fmt.Errorf("failed to subscribe to %s %s: %w", symbol, interval, err)
		}

		// Store the user channel
		a.candleChannels[interval][symbol] = userCh

		// Start goroutine to convert and forward events
		go a.forwardCandleEvents(symbol, interval, internalCh, userCh)
	}

	return nil
}

// forwardCandleEvents converts OKEx candle events to common types and forwards them
func (a *WebSocketAdapter) forwardCandleEvents(symbol, interval string, internalCh chan *publicevents.Candlesticks, userCh chan *commontypes.CandleUpdate) {
	for event := range internalCh {
		// OKEx Candlesticks event contains multiple candle updates
		for _, candle := range event.Candles {
			// Convert OKEx candle to common CandleUpdate
			update := a.converter.ConvertCandleToUpdate(symbol, interval, candle)
			if update == nil {
				continue
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
}

// UnsubscribeCandles unsubscribes from candlestick updates for specified symbols
func (a *WebSocketAdapter) UnsubscribeCandles(interval string, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols specified")
	}

	// Convert common interval format to OKEx format
	okexInterval := "candle" + interval

	for _, symbol := range symbols {
		// Unsubscribe from OKEx candlesticks channel
		req := publicrequests.Candlesticks{
			InstID:  symbol,
			Channel: okexconstants.CandleStickWsBarSize(okexInterval),
		}
		if err := a.client.Public.UCandlesticks(req, true); err != nil {
			return fmt.Errorf("failed to unsubscribe from %s %s: %w", symbol, interval, err)
		}

		// Remove from tracking
		if a.candleChannels[interval] != nil {
			delete(a.candleChannels[interval], symbol)
		}
	}

	return nil
}
