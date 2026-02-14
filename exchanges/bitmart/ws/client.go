package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/djpken/go-exc/exchanges/bitmart/utils"
	commontypes "github.com/djpken/go-exc/types"
	"github.com/gorilla/websocket"
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

	// Event channels for system events
	errCh       chan *commontypes.WebSocketError
	subCh       chan *commontypes.WebSocketSubscribe
	unsubCh     chan *commontypes.WebSocketUnsubscribe
	loginCh     chan *commontypes.WebSocketLogin
	successCh   chan *commontypes.WebSocketSuccess
	systemMsgCh chan *commontypes.WebSocketSystemMessage
	systemErrCh chan *commontypes.WebSocketSystemError

	// API endpoints
	Public  *Public
	Private *Private

	// Reconnection support
	reconnecting    bool
	reconnectMu     sync.Mutex
	subscriptions   []string // Track subscribed channels for resubscription
	subscriptionsMu sync.RWMutex
	connCtx         context.Context
	connCancel      context.CancelFunc
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

	if c.isConnected {
		c.mu.Unlock()
		return errors.New("already connected")
	}
	c.mu.Unlock()

	// Emit without holding lock
	c.emitSystemMessage("connection", "Connecting to WebSocket...", false)

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	conn, _, err := dialer.Dial(c.wsURL, nil)
	if err != nil {
		c.emitSystemError("connection", fmt.Sprintf("Failed to connect: %v", err), false)
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.isConnected = true

	// Create connection-specific context
	connCtx, connCancel := context.WithCancel(c.ctx)
	c.connCtx = connCtx
	c.connCancel = connCancel
	c.mu.Unlock()

	c.emitSystemMessage("connection", "WebSocket connected successfully", false)

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
// Format: {"action":"access","args":["<API_KEY>","<timestamp>","<sign>","<dev>"]}
// Sign: HmacSHA256(timestamp + "#" + memo + "#" + "bitmart.WebSocket", secretKey)
func (c *ClientWs) Login() error {
	c.mu.RLock()
	if !c.isConnected {
		c.mu.RUnlock()
		return errors.New("not connected")
	}
	c.mu.RUnlock()

	timestamp := utils.GetTimestamp()
	// Generate signature: timestamp + "#" + memo + "#" + "bitmart.WebSocket"
	sign := utils.GenerateSignature(timestamp, c.memo, "bitmart.WebSocket", c.secretKey)

	// BitMart uses "action": "access" for login
	loginMsg := map[string]interface{}{
		"action": "access",
		"args":   []string{c.apiKey, timestamp, sign, "web"},
	}

	// fmt.Printf("[WS] Logging in with timestamp: %s\n", timestamp)

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.conn.WriteJSON(loginMsg); err != nil {
		return fmt.Errorf("failed to send login: %w", err)
	}

	// Note: Authentication status will be set to true when receiving success response
	return nil
}

// Subscribe subscribes to a channel
func (c *ClientWs) Subscribe(channel string) error {
	c.mu.RLock()
	if !c.isConnected {
		c.mu.RUnlock()
		return errors.New("not connected")
	}
	conn := c.conn
	c.mu.RUnlock()

	// BitMart v2 API uses "action" instead of "op"
	subscribeMsg := map[string]interface{}{
		"action": "subscribe",
		"args":   []string{channel},
	}

	err := conn.WriteJSON(subscribeMsg)
	if err == nil {
		// Track subscription for resubscription on reconnect
		c.subscriptionsMu.Lock()
		c.subscriptions = append(c.subscriptions, channel)
		c.subscriptionsMu.Unlock()

		c.emitSubscribeEvent(channel)
		c.emitSystemMessage("subscription", fmt.Sprintf("Subscribed to %s", channel), false)
	} else {
		c.emitSystemError("subscription", fmt.Sprintf("Failed to subscribe to %s: %v", channel, err), false)
	}

	return err
}

// Unsubscribe unsubscribes from a channel
func (c *ClientWs) Unsubscribe(channel string) error {
	c.mu.RLock()
	if !c.isConnected {
		c.mu.RUnlock()
		return errors.New("not connected")
	}
	conn := c.conn
	c.mu.RUnlock()

	// BitMart v2 API uses "action" instead of "op"
	unsubscribeMsg := map[string]interface{}{
		"action": "unsubscribe",
		"args":   []string{channel},
	}

	err := conn.WriteJSON(unsubscribeMsg)
	if err == nil {
		// Remove from subscription tracking
		c.subscriptionsMu.Lock()
		filtered := make([]string, 0)
		for _, sub := range c.subscriptions {
			if sub != channel {
				filtered = append(filtered, sub)
			}
		}
		c.subscriptions = filtered
		c.subscriptionsMu.Unlock()

		c.emitUnsubscribeEvent(channel)
		c.emitSystemMessage("subscription", fmt.Sprintf("Unsubscribed from %s", channel), false)
	} else {
		c.emitSystemError("subscription", fmt.Sprintf("Failed to unsubscribe from %s: %v", channel, err), false)
	}

	return err
}

// startHeartbeat sends ping messages periodically
func (c *ClientWs) startHeartbeat() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.connCtx.Done():
			return
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			if !c.isConnected {
				c.mu.Unlock()
				return
			}
			// Send ping using BitMart v2 API format
			pingMsg := map[string]string{"action": "ping"}
			if err := c.conn.WriteJSON(pingMsg); err != nil {
				fmt.Printf("[WS ERROR] Failed to send ping: %v\n", err)
				c.emitSystemError("heartbeat", fmt.Sprintf("Failed to send ping: %v", err), false)
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
		c.emitSystemMessage("connection", "WebSocket connection closed", false)
	}()

	for {
		select {
		case <-c.connCtx.Done():
			c.emitSystemMessage("connection", "WebSocket closed by connection context cancellation", false)
			return
		case <-c.ctx.Done():
			c.emitSystemMessage("connection", "WebSocket closed by context cancellation", false)
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				// Log all errors for debugging
				errMsg := fmt.Sprintf("WebSocket read error: %v (type: %T)", err, err)
				fmt.Printf("[WS ERROR] %s\n", errMsg)
				c.emitSystemError("receiver", errMsg, false)

				// Also emit detailed close error
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					c.emitSystemMessage("connection", "WebSocket closed normally", false)
				} else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.emitSystemError("connection", fmt.Sprintf("Unexpected close: %v", err), false)
					// Attempt to reconnect on unexpected close
					go c.reconnect()
				} else {
					// Network error or other issue
					c.emitSystemError("connection", fmt.Sprintf("Connection error: %v", err), false)
					// Attempt to reconnect on error
					go c.reconnect()
				}
				return
			}

			// Process message
			c.processMessage(message)
		}
	}
}

// reconnect handles the reconnection logic when the connection is lost
func (c *ClientWs) reconnect() {
	c.reconnectMu.Lock()
	// Check if already reconnecting
	if c.reconnecting {
		c.reconnectMu.Unlock()
		return
	}
	c.reconnecting = true
	c.reconnectMu.Unlock()

	// Cancel old connection context to stop heartbeat goroutine
	c.mu.Lock()
	if c.connCancel != nil {
		c.connCancel()
	}

	// Close existing connection if any
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}

	c.isConnected = false
	c.isAuthenticated = false
	c.mu.Unlock()

	// Wait a bit for old goroutines to exit
	time.Sleep(100 * time.Millisecond)

	// Attempt to reconnect with exponential backoff
	c.emitSystemMessage("reconnection", "Attempting to reconnect...", false)

	maxRetries := 10
	retryInterval := 2 * time.Second
	maxInterval := 30 * time.Second

	for retry := 0; retry < maxRetries; retry++ {
		err := c.Connect()
		if err == nil {
			c.emitSystemMessage("reconnection", "Successfully reconnected", false)

			// Re-authenticate if it was a private connection
			if c.apiKey != "" && c.secretKey != "" && c.memo != "" {
				time.Sleep(500 * time.Millisecond) // Give connection time to stabilize
				if err := c.Login(); err != nil {
					c.emitSystemError("reconnection", fmt.Sprintf("Failed to re-authenticate: %v", err), false)
				} else {
					// Wait for authentication
					time.Sleep(500 * time.Millisecond)
				}
			}

			// Re-subscribe to all previous subscriptions
			c.resubscribe()

			c.reconnectMu.Lock()
			c.reconnecting = false
			c.reconnectMu.Unlock()
			return
		}

		c.emitSystemError("reconnection", fmt.Sprintf("Reconnection attempt %d failed: %v", retry+1, err), false)

		// Exponential backoff
		time.Sleep(retryInterval)
		retryInterval = time.Duration(float64(retryInterval) * 1.5)
		if retryInterval > maxInterval {
			retryInterval = maxInterval
		}
	}

	c.reconnectMu.Lock()
	c.reconnecting = false
	c.reconnectMu.Unlock()

	c.emitSystemError("reconnection", fmt.Sprintf("Failed to reconnect after %d attempts", maxRetries), false)
}

// resubscribe re-subscribes to all saved subscriptions after reconnection
func (c *ClientWs) resubscribe() {
	c.subscriptionsMu.RLock()
	subs := make([]string, len(c.subscriptions))
	copy(subs, c.subscriptions)
	c.subscriptionsMu.RUnlock()

	if len(subs) == 0 {
		return
	}

	c.emitSystemMessage("subscription", fmt.Sprintf("Re-subscribing to %d channel(s)...", len(subs)), false)

	// Clear subscriptions before re-subscribing to avoid duplication
	c.subscriptionsMu.Lock()
	c.subscriptions = []string{}
	c.subscriptionsMu.Unlock()

	// Re-subscribe to each saved subscription
	for _, channel := range subs {
		time.Sleep(100 * time.Millisecond) // Rate limit subscriptions
		err := c.Subscribe(channel)
		if err != nil {
			c.emitSystemError("subscription", fmt.Sprintf("Failed to re-subscribe to %s: %v", channel, err), false)
		}
	}

	c.emitSystemMessage("subscription", "Re-subscription complete", false)
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

	// Debug: print raw message to understand format
	// fmt.Printf("\n[WS RAW] %s\n", string(message))

	// Handle pong response
	if event, ok := msg["event"].(string); ok && event == "pong" {
		return
	}

	// Also check for "action": "pong" format (BitMart v2 API might use this)
	if action, ok := msg["action"].(string); ok && action == "pong" {
		return
	}

	// Handle login response (action: "access", success: true)
	if action, ok := msg["action"].(string); ok && action == "access" {
		if success, ok := msg["success"].(bool); ok {
			if success {
				c.mu.Lock()
				c.isAuthenticated = true
				c.mu.Unlock()
				c.emitLoginEvent(true, "WebSocket authenticated successfully")
				c.emitSystemMessage("login", "Authentication successful", true)
			} else {
				c.emitLoginEvent(false, "WebSocket authentication failed")
				c.emitSystemError("login", "Authentication failed", true)
			}
		}
		return
	}

	// Debug: print message structure
	// fmt.Printf("[WS] Message keys: %v\n", getMapKeys(msg))

	// Route data messages to registered handlers
	// BitMart uses "table" for public channels and "group" for private channels
	var channelKey string
	if table, ok := msg["table"].(string); ok {
		channelKey = table
	} else if group, ok := msg["group"].(string); ok {
		channelKey = group
	}

	if channelKey != "" {
		// fmt.Printf("[WS] Routing to handler: %s\n", channelKey)
		c.mu.RLock()
		handler, exists := c.handlers[channelKey]
		c.mu.RUnlock()

		if exists {
			// fmt.Printf("[WS] Handler found, executing\n")
			handler(message)
		} else {
			// fmt.Printf("[WS] No handler for channel: %s\n", channelKey)
			// fmt.Printf("[WS] Available handlers: %v\n", c.getHandlerKeys())
		}
	} else {
		// fmt.Printf("[WS] No 'table' or 'group' field in message\n")
	}
}

// getHandlerKeys returns all registered handler keys for debugging
func (c *ClientWs) getHandlerKeys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.handlers))
	for k := range c.handlers {
		keys = append(keys, k)
	}
	return keys
}

// getMapKeys returns all keys in a map for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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

// SetChannels sets channels for receiving WebSocket events
// This allows users to receive notifications about connection events, errors, subscriptions, etc.
func (c *ClientWs) SetChannels(
	errCh chan *commontypes.WebSocketError,
	subCh chan *commontypes.WebSocketSubscribe,
	unsubCh chan *commontypes.WebSocketUnsubscribe,
	loginCh chan *commontypes.WebSocketLogin,
	successCh chan *commontypes.WebSocketSuccess,
	systemMsgCh chan *commontypes.WebSocketSystemMessage,
	systemErrCh chan *commontypes.WebSocketSystemError,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.errCh = errCh
	c.subCh = subCh
	c.unsubCh = unsubCh
	c.loginCh = loginCh
	c.successCh = successCh
	c.systemMsgCh = systemMsgCh
	c.systemErrCh = systemErrCh
}

// emitSystemMessage sends a system message event if the channel is set
// Note: This function is safe to call without holding locks
func (c *ClientWs) emitSystemMessage(msgType, message string, private bool) {
	// Read channel pointer without lock - it's safe since SetChannels is called once at setup
	ch := c.systemMsgCh

	if ch != nil {
		select {
		case ch <- &commontypes.WebSocketSystemMessage{
			Type:      msgType,
			Message:   message,
			Private:   private,
			Timestamp: commontypes.Timestamp(time.Now()),
			Extra:     make(map[string]interface{}),
		}:
		default:
			// Channel full, drop message
		}
	}
}

// emitSystemError sends a system error event if the channel is set
// Note: This function is safe to call without holding locks
func (c *ClientWs) emitSystemError(errType, errMsg string, private bool) {
	ch := c.systemErrCh

	if ch != nil {
		select {
		case ch <- &commontypes.WebSocketSystemError{
			Type:      errType,
			Error:     errMsg,
			Private:   private,
			Timestamp: commontypes.Timestamp(time.Now()),
			Extra:     make(map[string]interface{}),
		}:
		default:
			// Channel full, drop message
		}
	}
}

// emitLoginEvent sends a login event if the channel is set
// Note: This function is safe to call without holding locks
func (c *ClientWs) emitLoginEvent(success bool, message string) {
	ch := c.loginCh

	if ch != nil {
		code := "0"
		if !success {
			code = "1"
		}

		select {
		case ch <- &commontypes.WebSocketLogin{
			Event:   "login",
			Code:    code,
			Message: message,
			Extra:   make(map[string]interface{}),
		}:
		default:
			// Channel full, drop message
		}
	}
}

// emitSubscribeEvent sends a subscription event if the channel is set
// Note: This function is safe to call without holding locks
func (c *ClientWs) emitSubscribeEvent(channel string) {
	ch := c.subCh

	if ch != nil {
		select {
		case ch <- &commontypes.WebSocketSubscribe{
			Event:   "subscribe",
			Channel: channel,
			Extra:   make(map[string]interface{}),
		}:
		default:
			// Channel full, drop message
		}
	}
}

// emitUnsubscribeEvent sends an unsubscription event if the channel is set
// Note: This function is safe to call without holding locks
func (c *ClientWs) emitUnsubscribeEvent(channel string) {
	ch := c.unsubCh

	if ch != nil {
		select {
		case ch <- &commontypes.WebSocketUnsubscribe{
			Event:   "unsubscribe",
			Channel: channel,
			Extra:   make(map[string]interface{}),
		}:
		default:
			// Channel full, drop message
		}
	}
}

// emitErrorEvent sends an error event if the channel is set
// Note: This function is safe to call without holding locks
func (c *ClientWs) emitErrorEvent(code, message string) {
	ch := c.errCh

	if ch != nil {
		select {
		case ch <- &commontypes.WebSocketError{
			Event:   "error",
			Code:    code,
			Message: message,
			Extra:   make(map[string]interface{}),
		}:
		default:
			// Channel full, drop message
		}
	}
}
