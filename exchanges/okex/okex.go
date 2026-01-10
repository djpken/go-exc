package okex

import (
	"context"

	"github.com/djpken/go-exc/exchanges/okex/rest"
	"github.com/djpken/go-exc/exchanges/okex/ws"
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

// GetNativeClient returns the native OKEx client for advanced usage
func (e *OKExExchange) GetNativeClient() *Client {
	return e.client
}

// GetNativeRest returns the native REST client for advanced usage
func (e *OKExExchange) GetNativeRest() *rest.ClientRest {
	return e.client.Rest
}

// GetNativeWs returns the native WebSocket client for advanced usage
func (e *OKExExchange) GetNativeWs() *ws.ClientWs {
	return e.client.Ws
}
