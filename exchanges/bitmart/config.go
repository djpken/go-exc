package bitmart

import (
	"errors"

	"github.com/djpken/go-exc/exchanges/bitmart/types"
)

// Config contains BitMart-specific configuration
type Config struct {
	// APIKey is the API key
	APIKey string

	// SecretKey is the secret key
	SecretKey string

	// Memo is the API memo (required by BitMart)
	Memo string

	// BaseURL is the base URL for REST API
	BaseURL string

	// WSBaseURL is the base URL for WebSocket API
	WSBaseURL string
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("APIKey is required")
	}
	if c.SecretKey == "" {
		return errors.New("SecretKey is required")
	}
	if c.Memo == "" {
		return errors.New("Memo is required")
	}
	return nil
}

// GetBaseURL returns the base URL for REST API
func (c *Config) GetBaseURL() string {
	return c.BaseURL
}

// GetWSBaseURL returns the base URL for WebSocket API
func (c *Config) GetWSBaseURL() string {
	return c.WSBaseURL
}

// NewDefaultConfig creates a new default configuration
func NewDefaultConfig(apiKey, secretKey, memo string, testMode bool) *Config {
	baseURL := string(types.ProductionSwapServer)
	if testMode {
		baseURL = string(types.DemoSwapServer)
	}
	return &Config{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		BaseURL:   baseURL,
		WSBaseURL: string(types.ProductionAPIWSServer),
	}
}
