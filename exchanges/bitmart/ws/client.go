package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/djpken/go-exc/exchanges/bitmart/utils"
)

// MessageHandler is a function that handles WebSocket messages
type MessageHandler func([]byte)

// ClientWs represents the BitMart WebSocket API client
type ClientWs struct {
	ctx       context.Context
	conn      *websocket.Conn
	apiKey    string
	secretKey string
	memo      string
	wsURL     string

	mu              sync.RWMutex
	isConnected     bool
	isAuthenticated bool
	handlers        map[string]MessageHandler

	// API endpoints
	Public  *Public
	Private *Private
}

// Config interface for BitMart WebSocket configuration
type Config interface {
	GetWSBaseURL() string
	Validate() error
}

// BitMartConfig represents BitMart WebSocket configuration
type BitMartConfig struct {
	APIKey    string
	SecretKey string
	Memo      string
	WSBaseURL string
}

// GetWSBaseURL returns the WebSocket base URL
func (c *BitMartConfig) GetWSBaseURL() string {
	return c.WSBaseURL
}

// Validate validates the configuration
func (c *BitMartConfig) Validate() error {
	return nil
}

// NewClientWs creates a new BitMart WebSocket API client
func NewClientWs(ctx context.Context, cfg Config) (*ClientWs, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Type assertion to get actual config values
	bmConfig, ok := cfg.(*BitMartConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type")
	}

	client := &ClientWs{
		ctx:             ctx,
		apiKey:          bmConfig.APIKey,
		secretKey:       bmConfig.SecretKey,
		memo:            bmConfig.Memo,
		wsURL:           cfg.GetWSBaseURL(),
		isConnected:     false,
		isAuthenticated: false,
		handlers:        make(map[string]MessageHandler),
	}

	// Initialize API endpoints
	client.Public = NewPublic(client)
	client.Private = NewPrivate(client)

	return client, nil
}

// Connect establishes WebSocket connection
func (c *ClientWs) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return errors.New("already connected")
	}

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	conn, _, err := dialer.Dial(c.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.conn = conn
	c.isConnected = true

	// Start message reader
	go c.readMessages()
	// Start ping/pong heartbeat
	go c.startHeartbeat()

	return nil
}

// Close closes the WebSocket connection
func (c *ClientWs) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return nil
	}

	c.isConnected = false
	c.isAuthenticated = false

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// Login authenticates the WebSocket connection for private channels
func (c *ClientWs) Login() error {
	c.mu.RLock()
	if !c.isConnected {
		c.mu.RUnlock()
		return errors.New("not connected")
	}
	c.mu.RUnlock()

	timestamp := utils.GetTimestamp()
	message := timestamp + "#" + c.memo + "#" + "bitmart.WebSocket"
	sign := utils.GenerateSignature(timestamp, message, c.secretKey)

	loginMsg := map[string]interface{}{
		"op":   "login",
		"args": []string{c.apiKey, timestamp, sign},
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.conn.WriteJSON(loginMsg); err != nil {
		return fmt.Errorf("failed to send login: %w", err)
	}

	c.isAuthenticated = true
	return nil
}

// Subscribe subscribes to a channel
func (c *ClientWs) Subscribe(channel string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected {
		return errors.New("not connected")
	}

	subscribeMsg := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{channel},
	}

	return c.conn.WriteJSON(subscribeMsg)
}

// Unsubscribe unsubscribes from a channel
func (c *ClientWs) Unsubscribe(channel string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected {
		return errors.New("not connected")
	}

	unsubscribeMsg := map[string]interface{}{
		"op":   "unsubscribe",
		"args": []string{channel},
	}

	return c.conn.WriteJSON(unsubscribeMsg)
}

// startHeartbeat sends ping messages periodically
func (c *ClientWs) startHeartbeat() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			if !c.isConnected {
				c.mu.Unlock()
				return
			}
			// Send ping
			pingMsg := map[string]string{"op": "ping"}
			if err := c.conn.WriteJSON(pingMsg); err != nil {
				fmt.Printf("Failed to send ping: %v\n", err)
			}
			c.mu.Unlock()
		}
	}
}

// readMessages reads messages from WebSocket
func (c *ClientWs) readMessages() {
	defer func() {
		c.mu.Lock()
		c.isConnected = false
		c.isAuthenticated = false
		c.mu.Unlock()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket read error: %v\n", err)
				}
				return
			}

			// Process message
			c.processMessage(message)
		}
	}
}

// RegisterHandler registers a message handler for a channel
func (c *ClientWs) RegisterHandler(channel string, handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[channel] = handler
}

// UnregisterHandler unregisters a message handler for a channel
func (c *ClientWs) UnregisterHandler(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.handlers, channel)
}

// processMessage processes incoming WebSocket message
func (c *ClientWs) processMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Printf("Failed to unmarshal message: %v\n", err)
		return
	}

	// Handle pong response
	if event, ok := msg["event"].(string); ok && event == "pong" {
		return
	}

	// Handle login response
	if event, ok := msg["event"].(string); ok && event == "login" {
		if success, ok := msg["success"].(bool); ok && success {
			fmt.Println("WebSocket authenticated successfully")
		}
		return
	}

	// Route data messages to registered handlers
	if table, ok := msg["table"].(string); ok {
		c.mu.RLock()
		handler, exists := c.handlers[table]
		c.mu.RUnlock()

		if exists {
			handler(message)
		}
	}
}

// IsConnected returns connection status
func (c *ClientWs) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// IsAuthenticated returns authentication status
func (c *ClientWs) IsAuthenticated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isAuthenticated
}
