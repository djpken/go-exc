package bitmart

import (
	"context"
	"errors"

	"github.com/djpken/go-exc/exchanges/bitmart/rest"
	"github.com/djpken/go-exc/exchanges/bitmart/ws"
)

// Client represents the BitMart native client
type Client struct {
	// Rest is the REST API client
	Rest *rest.ClientRest

	// Ws is the WebSocket API client
	Ws *ws.ClientWs

	ctx    context.Context
	config *Config
}

// NewClient creates a new BitMart native client
func NewClient(ctx context.Context, apiKey, secretKey, memo string, testMode bool) (*Client, error) {
	config := NewDefaultConfig(apiKey, secretKey, memo, testMode)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Create REST client config
	restConfig := &rest.BitMartConfig{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		BaseURL:   config.GetBaseURL(),
	}
	restClient, err := rest.NewClientRest(ctx, restConfig)
	if err != nil {
		return nil, err
	}

	// Create WebSocket client config
	wsConfig := &ws.BitMartConfig{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		WSBaseURL: config.GetWSBaseURL(),
	}
	wsClient, err := ws.NewClientWs(ctx, wsConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		Rest:   restClient,
		Ws:     wsClient,
		ctx:    ctx,
		config: config,
	}, nil
}

// Close closes all connections
func (c *Client) Close() error {
	var errs []error

	if c.Ws != nil {
		if err := c.Ws.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
