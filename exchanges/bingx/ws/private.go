package ws

import (
	"fmt"
	"sync"
	"time"
)

const listenKeyRenewInterval = 30 * time.Minute

// PrivateClientWs manages a BingX private WebSocket connection using listen key auth.
// Private events (ACCOUNT_UPDATE, ORDER_TRADE_UPDATE) are pushed automatically once connected
// — no subscribe message is needed after connection.
type PrivateClientWs struct {
	getListenKey    func() (string, error)
	extendListenKey func(key string) error
	baseURL         string // WebSocket base URL (empty = default)

	mu        sync.Mutex
	client    *ClientWs
	listenKey string

	done    chan struct{}
	started bool
}

// NewPrivateClientWs creates a new private WebSocket client.
// getKey is called to obtain a fresh listen key.
// extendKey is called every 30 minutes to reset the listen key expiry.
// baseURL is the WebSocket base URL (empty = default production URL).
func NewPrivateClientWs(getKey func() (string, error), extendKey func(key string) error, baseURL string) *PrivateClientWs {
	return &PrivateClientWs{
		getListenKey:    getKey,
		extendListenKey: extendKey,
		baseURL:         baseURL,
		done:            make(chan struct{}),
	}
}

// EnsureConnected lazily connects to the private WebSocket on first call.
// Subsequent calls are no-ops when the connection is already active.
func (p *PrivateClientWs) EnsureConnected() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil && p.client.IsConnected() {
		return nil
	}

	key, err := p.getListenKey()
	if err != nil {
		return fmt.Errorf("bingx private ws: get listen key: %w", err)
	}

	p.listenKey = key
	p.client = NewClientWs(p.baseURL, key)

	if err := p.client.Connect(); err != nil {
		return fmt.Errorf("bingx private ws: connect: %w", err)
	}

	if !p.started {
		p.started = true
		go p.renewLoop()
	}
	return nil
}

// RegisterHandler registers a handler for a private event type (e.g., "ACCOUNT_UPDATE").
// Connects to the private WebSocket if not already connected.
func (p *PrivateClientWs) RegisterHandler(eventType string, h Handler) error {
	if err := p.EnsureConnected(); err != nil {
		return err
	}
	p.mu.Lock()
	p.client.RegisterHandler(eventType, h)
	p.mu.Unlock()
	return nil
}

// UnregisterHandler removes the handler for the given event type.
func (p *PrivateClientWs) UnregisterHandler(eventType string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client != nil {
		p.client.UnregisterHandler(eventType)
	}
}

// Close shuts down the private WebSocket connection and stops the renew loop.
func (p *PrivateClientWs) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	select {
	case <-p.done:
	default:
		close(p.done)
	}
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}

// renewLoop sends a keep-alive PUT every 30 minutes to prevent the listen key from expiring.
func (p *PrivateClientWs) renewLoop() {
	ticker := time.NewTicker(listenKeyRenewInterval)
	defer ticker.Stop()
	for {
		select {
		case <-p.done:
			return
		case <-ticker.C:
			p.mu.Lock()
			key := p.listenKey
			p.mu.Unlock()
			if key != "" {
				_ = p.extendListenKey(key)
			}
		}
	}
}
