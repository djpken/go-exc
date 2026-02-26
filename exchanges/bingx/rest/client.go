package rest

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://open-api.bingx.com"
	testBaseURL    = "https://open-api.bingx.com" // BingX simulation trading uses the same URL with demo account credentials
)

// ClientRest is the BingX REST API client
type ClientRest struct {
	ctx        context.Context
	httpClient *http.Client
	apiKey     string
	secretKey  string
	baseURL    string

	Market  *Market
	Account *Account
	Trade   *Trade
}

// Response is the standard BingX API response envelope
type Response[T any] struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      T      `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

// NewClientRest creates a new BingX REST client
func NewClientRest(ctx context.Context, apiKey, secretKey, baseURL string) *ClientRest {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	c := &ClientRest{
		ctx:       ctx,
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	c.Market = NewMarket(c)
	c.Account = NewAccount(c)
	c.Trade = NewTrade(c)
	return c
}

// sign computes the HMAC-SHA256 hex signature for the given query string
func (c *ClientRest) sign(queryString string) string {
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(queryString))
	return hex.EncodeToString(mac.Sum(nil))
}

// timestamp returns current Unix millisecond timestamp as string
func timestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// GET performs an authenticated GET request.
// params should be passed as key=value pairs without timestamp/signature.
func (c *ClientRest) GET(path string, params map[string]string, result interface{}) error {
	return c.do(http.MethodGet, path, params, result)
}

// POST performs an authenticated POST request.
func (c *ClientRest) POST(path string, params map[string]string, result interface{}) error {
	return c.do(http.MethodPost, path, params, result)
}

// PUT performs an authenticated PUT request (params sent as URL query string).
func (c *ClientRest) PUT(path string, params map[string]string, result interface{}) error {
	return c.do(http.MethodPut, path, params, result)
}

// DELETE performs an authenticated DELETE request.
func (c *ClientRest) DELETE(path string, params map[string]string, result interface{}) error {
	return c.do(http.MethodDelete, path, params, result)
}

// GETPublic performs an unauthenticated GET request (for public endpoints)
func (c *ClientRest) GETPublic(path string, params map[string]string, result interface{}) error {
	qs := buildQueryString(params)
	url := c.baseURL + path
	if qs != "" {
		url += "?" + qs
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("bingx: create request: %w", err)
	}

	return c.executeAndDecode(req, result)
}

func (c *ClientRest) do(method, path string, params map[string]string, result interface{}) error {
	if params == nil {
		params = make(map[string]string)
	}
	params["timestamp"] = timestamp()

	qs := buildQueryString(params)
	params["signature"] = c.sign(qs)
	qs = buildQueryString(params)

	var req *http.Request
	var err error

	url := c.baseURL + path
	if method == http.MethodGet || method == http.MethodDelete || method == http.MethodPut {
		req, err = http.NewRequestWithContext(c.ctx, method, url+"?"+qs, nil)
	} else {
		// POST: send params as URL-encoded body
		req, err = http.NewRequestWithContext(c.ctx, method, url, strings.NewReader(qs))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	if err != nil {
		return fmt.Errorf("bingx: create request: %w", err)
	}

	req.Header.Set("X-BX-APIKEY", c.apiKey)

	return c.executeAndDecode(req, result)
}

func (c *ClientRest) executeAndDecode(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("bingx: http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("bingx: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bingx: http %d: %s", resp.StatusCode, body)
	}

	// Check business error code
	var envelope struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("bingx: decode envelope: %w", err)
	}
	if envelope.Code != 0 {
		return fmt.Errorf("bingx: api error %d: %s", envelope.Code, envelope.Msg)
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("bingx: decode response: %w", err)
		}
	}
	return nil
}

// buildQueryString builds a URL query string from a map (no URL encoding for signature calculation)
func buildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	parts := make([]string, 0, len(params))
	for k, v := range params {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, "&")
}
