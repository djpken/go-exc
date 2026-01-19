package bitmart

import (
	"context"

	"github.com/djpken/go-exc/exchanges/bitmart/rest"
	"github.com/djpken/go-exc/exchanges/bitmart/ws"
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
