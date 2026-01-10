package okex

import (
	"github.com/djpken/go-exc/exchanges/okex/ws"
)

// WebSocketAdapter adapts OKEx WebSocket client to common interface
type WebSocketAdapter struct {
	client    *ws.ClientWs
	converter *Converter
}

// NewWebSocketAdapter creates a new WebSocket adapter
func NewWebSocketAdapter(client *ws.ClientWs) *WebSocketAdapter {
	return &WebSocketAdapter{
		client:    client,
		converter: NewConverter(),
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
