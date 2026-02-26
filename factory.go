package exc

import (
	"context"

	"github.com/djpken/go-exc/exchanges/bingx"
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
	case Bitmart:
		return newBitMartExchange(ctx, cfg, false)
	case BitmartTest:
		return newBitMartExchange(ctx, cfg, true)
	case BingX:
		return newBingXExchange(ctx, cfg, false)
	case BingXTest:
		return newBingXExchange(ctx, cfg, true)
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

// newBingXExchange creates a new BingX exchange instance
func newBingXExchange(ctx context.Context, cfg Config, testMode bool) (Exchange, error) {
	client, err := bingx.NewBingXExchange(ctx, cfg.APIKey, cfg.SecretKey, testMode)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// newBitMartExchange creates a new Bitmart exchange instance
func newBitMartExchange(ctx context.Context, cfg Config, testMode bool) (Exchange, error) {
	// Bitmart requires memo which should be in Extra["memo"]
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
