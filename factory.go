package exc

import (
	"context"

	"github.com/djpken/go-exc/exchanges/bitmart"
	"github.com/djpken/go-exc/exchanges/okex"
)

// NewExchange creates a new exchange instance based on the exchange type
func NewExchange(ctx context.Context, exchangeType ExchangeType, cfg Config) (Exchange, error) {
	switch exchangeType {
	case OKX:
		return newOKExExchange(ctx, cfg, false)
	case OKXTest:
		return newOKExExchange(ctx, cfg, true)
	case BitMart:
		return newBitMartExchange(ctx, cfg, false)
	case BitMartTest:
		return newBitMartExchange(ctx, cfg, true)
	case BingX:
		return nil, ErrNotImplemented
	default:
		return nil, ErrInvalidExchange
	}
}

// newOKExExchange creates a new OKEx exchange instance
func newOKExExchange(ctx context.Context, cfg Config, testMode bool) (Exchange, error) {
	client, err := okex.NewOKExExchange(ctx, cfg.APIKey, cfg.SecretKey, cfg.Passphrase, testMode)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// newBitMartExchange creates a new BitMart exchange instance
func newBitMartExchange(ctx context.Context, cfg Config, testMode bool) (Exchange, error) {
	// BitMart requires memo which should be in Extra["memo"]
	memo := ""
	if cfg.Extra != nil {
		if m, ok := cfg.Extra["memo"].(string); ok {
			memo = m
		}
	}

	client, err := bitmart.NewBitMartExchange(ctx, cfg.APIKey, cfg.SecretKey, memo, testMode)
	if err != nil {
		return nil, err
	}
	return client, nil
}
