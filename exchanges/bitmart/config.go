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
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return string(types.ProductionServer)
}

// GetWSBaseURL returns the base URL for WebSocket API
func (c *Config) GetWSBaseURL() string {
	if c.WSBaseURL != "" {
		return c.WSBaseURL
	}
	return string(types.ProductionWSServer)
}

// NewDefaultConfig creates a new default configuration
func NewDefaultConfig(apiKey, secretKey, memo string) *Config {
	return &Config{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		BaseURL:   string(types.ProductionServer),
		WSBaseURL: string(types.ProductionWSServer),
	}
}
