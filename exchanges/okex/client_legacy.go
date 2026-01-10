package okex

import (
	"context"

	"github.com/djpken/go-exc/exchanges/okex/rest"
	"github.com/djpken/go-exc/exchanges/okex/ws"
)

// Client is the main api wrapper of okex
type Client struct {
	Rest *rest.ClientRest
	Ws   *ws.ClientWs
	ctx  context.Context
}

// NewClient returns a pointer to a fresh Client
func NewClient(ctx context.Context, apiKey, secretKey, passphrase string, destination Destination) (*Client, error) {
	restURL := RestURL
	wsPubURL := PublicWsURL
	wsPriURL := PrivateWsURL
	switch destination {
	case AwsServer:
		restURL = AwsRestURL
		wsPubURL = AwsPublicWsURL
		wsPriURL = AwsPrivateWsURL
	case DemoServer:
		restURL = DemoRestURL
		wsPubURL = DemoPublicWsURL
		wsPriURL = DemoPrivateWsURL
	}

	r := rest.NewClient(apiKey, secretKey, passphrase, restURL, destination)
	c := ws.NewClient(ctx, apiKey, secretKey, passphrase, map[bool]BaseURL{true: wsPriURL, false: wsPubURL})

	return &Client{r, c, ctx}, nil
}
