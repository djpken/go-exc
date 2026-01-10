package exc

import (
	"context"

	"github.com/djpken/go-exc/exchanges/okex"
)

// NewExchange creates a new exchange instance based on the exchange type
func NewExchange(ctx context.Context, exchangeType ExchangeType, cfg Config) (Exchange, error) {
	switch exchangeType {
	case OKEx:
		return newOKExExchange(ctx, cfg)
	case Binance:
		return nil, ErrNotImplemented
	case BitMart:
		return nil, ErrNotImplemented
	case BingX:
		return nil, ErrNotImplemented
	default:
		return nil, ErrInvalidExchange
	}
}

// newOKExExchange creates a new OKEx exchange instance
// Note: Currently returns the native OKEx client directly
// TODO: Implement full Exchange interface wrapper
func newOKExExchange(ctx context.Context, cfg Config) (Exchange, error) {
	_, err := okex.NewOKExExchange(ctx, cfg.APIKey, cfg.SecretKey, cfg.Passphrase, cfg.TestMode)
	if err != nil {
		return nil, err
	}

	// TODO: Wrap in Exchange interface adapter
	return nil, ErrNotImplemented
}

// NewOKExClient creates a native OKEx client directly
// This is the recommended way to use OKEx for now
func NewOKExClient(ctx context.Context, apiKey, secretKey, passphrase string, testMode bool) (*okex.OKExExchange, error) {
	return okex.NewOKExExchange(ctx, apiKey, secretKey, passphrase, testMode)
}
