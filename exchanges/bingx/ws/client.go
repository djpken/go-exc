package ws

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	swapPublicURL     = "wss://open-api-swap.bingx.com/swap-market"
	swapPublicTestURL = "wss://open-api-swap.bingx.com/swap-market" // BingX simulation trading uses the same WS URL with demo account credentials
	pingInterval      = 20 * time.Second
	reconnectDelay    = 3 * time.Second
)

// subRequest is the JSON message sent to subscribe/unsubscribe
type subRequest struct {
	ID       string `json:"id"`
	ReqType  string `json:"reqType"` // "sub" or "unsub"
	DataType string `json:"dataType"`
}

// Handler is called with the decompressed JSON message bytes for a dataType
type Handler func(data []byte)

// ClientWs manages the BingX WebSocket connection
type ClientWs struct {
	url       string
	listenKey string // non-empty for private connections

	mu       sync.RWMutex
	conn     *websocket.Conn
	handlers map[string]Handler // dataType -> handler

	done   chan struct{}
	closed bool
}

// NewClientWs creates a new public WebSocket client.
// baseURL is the WebSocket base URL (e.g. swapPublicURL or swapPublicTestURL).
// listenKey is non-empty for private connections.
func NewClientWs(baseURL, listenKey string) *ClientWs {
	if baseURL == "" {
		baseURL = swapPublicURL
	}
	url := baseURL
	if listenKey != "" {
		url += "?listenKey=" + listenKey
	}
	return &ClientWs{
		url:       url,
		listenKey: listenKey,
		handlers:  make(map[string]Handler),
		done:      make(chan struct{}),
	}
}

// Connect establishes the WebSocket connection and starts the read loop
func (c *ClientWs) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("bingx ws: connect: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.closed = false
	c.mu.Unlock()

	go c.readLoop()
	go c.pingLoop()
	return nil
}

// IsConnected returns true if the connection is active
func (c *ClientWs) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && !c.closed
}

// Close shuts down the WebSocket connection
func (c *ClientWs) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	close(c.done)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// RegisterHandler registers a handler for a specific dataType channel
func (c *ClientWs) RegisterHandler(dataType string, h Handler) {
	c.mu.Lock()
	c.handlers[dataType] = h
	c.mu.Unlock()
}

// UnregisterHandler removes a handler for a dataType
func (c *ClientWs) UnregisterHandler(dataType string) {
	c.mu.Lock()
	delete(c.handlers, dataType)
	c.mu.Unlock()
}

// Subscribe sends a subscription request for the given dataType
func (c *ClientWs) Subscribe(dataType string) error {
	return c.sendMsg("sub", dataType)
}

// Unsubscribe sends an unsubscription request for the given dataType
func (c *ClientWs) Unsubscribe(dataType string) error {
	return c.sendMsg("unsub", dataType)
}

func (c *ClientWs) sendMsg(reqType, dataType string) error {
	msg := subRequest{
		ID:       fmt.Sprintf("go-exc-%d", time.Now().UnixNano()),
		ReqType:  reqType,
		DataType: dataType,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil || c.closed {
		return fmt.Errorf("bingx ws: not connected")
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *ClientWs) readLoop() {
	for {
		select {
		case <-c.done:
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()
		if conn == nil {
			return
		}

		_, raw, err := conn.ReadMessage()
		if err != nil {
			select {
			case <-c.done:
				return
			default:
				// connection error; attempt reconnect
				c.reconnect()
				return
			}
		}

		// BingX compresses with GZIP
		decompressed, err := decompressGzip(raw)
		if err != nil {
			// Server ack messages may not be compressed; try as-is
			decompressed = raw
		}

		c.dispatch(decompressed)
	}
}

func (c *ClientWs) dispatch(data []byte) {
	// Parse the dataType field from the message
	var envelope struct {
		DataType string `json:"dataType"`
		E        string `json:"e"` // private event type field
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return
	}

	key := envelope.DataType
	if key == "" {
		key = envelope.E
	}

	c.mu.RLock()
	h, ok := c.handlers[key]
	c.mu.RUnlock()

	if ok {
		h(data)
	}
}

func (c *ClientWs) pingLoop() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.mu.Lock()
			conn := c.conn
			closed := c.closed
			c.mu.Unlock()
			if conn == nil || closed {
				return
			}
			_ = conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}

func (c *ClientWs) reconnect() {
	time.Sleep(reconnectDelay)
	select {
	case <-c.done:
		return
	default:
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return
	}
	c.mu.Lock()
	c.conn = conn
	c.closed = false
	c.mu.Unlock()

	// Re-subscribe to all registered channels
	c.mu.RLock()
	dataTypes := make([]string, 0, len(c.handlers))
	for dt := range c.handlers {
		dataTypes = append(dataTypes, dt)
	}
	c.mu.RUnlock()

	for _, dt := range dataTypes {
		_ = c.Subscribe(dt)
	}

	go c.readLoop()
}

func decompressGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
