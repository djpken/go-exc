package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/djpken/go-exc/exchanges/bitmart/utils"
)

// ClientRest represents the BitMart REST API client
type ClientRest struct {
	ctx        context.Context
	httpClient *http.Client
	apiKey     string
	secretKey  string
	memo       string
	baseURL    string

	// API endpoints
	Market  *Market
	Account *Account
	Trade   *Trade
	Funding *Funding
}

// Config interface for BitMart configuration
type Config interface {
	GetBaseURL() string
	Validate() error
}

// BitMartConfig represents BitMart configuration
type BitMartConfig struct {
	APIKey    string
	SecretKey string
	Memo      string
	BaseURL   string
}

// GetBaseURL returns the base URL
func (c *BitMartConfig) GetBaseURL() string {
	return c.BaseURL
}

// Validate validates the configuration
func (c *BitMartConfig) Validate() error {
	return nil
}

// NewClientRest creates a new BitMart REST API client
func NewClientRest(ctx context.Context, cfg Config) (*ClientRest, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Type assertion to get actual config values
	bmConfig, ok := cfg.(*BitMartConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type")
	}

	client := &ClientRest{
		ctx: ctx,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:    bmConfig.APIKey,
		secretKey: bmConfig.SecretKey,
		memo:      bmConfig.Memo,
		baseURL:   cfg.GetBaseURL(),
	}

	// Initialize API endpoints
	client.Market = NewMarket(client)
	client.Account = NewAccount(client)
	client.Trade = NewTrade(client)
	client.Funding = NewFunding(client)

	return client, nil
}

// doRequest performs HTTP request with authentication
func (c *ClientRest) doRequest(method, endpoint string, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(c.ctx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	timestamp := utils.GetTimestamp()
	sign := utils.GenerateSignature(timestamp, string(reqBody), c.secretKey)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BM-KEY", c.apiKey)
	req.Header.Set("X-BM-SIGN", sign)
	req.Header.Set("X-BM-TIMESTAMP", timestamp)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	// Parse response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// GET performs a GET request
func (c *ClientRest) GET(endpoint string, result interface{}) error {
	return c.doRequest(http.MethodGet, endpoint, nil, result)
}

// POST performs a POST request
func (c *ClientRest) POST(endpoint string, body interface{}, result interface{}) error {
	return c.doRequest(http.MethodPost, endpoint, body, result)
}

// DELETE performs a DELETE request
func (c *ClientRest) DELETE(endpoint string, body interface{}, result interface{}) error {
	return c.doRequest(http.MethodDelete, endpoint, body, result)
}
