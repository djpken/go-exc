package types

// BalanceAndPositionUpdate represents balance and position update from WebSocket
type BalanceAndPositionUpdate struct {
	// Balances is the list of currency balances
	Balances []*Balance

	// Positions is the list of positions
	Positions []*Position

	// EventType describes the type of event (e.g., "snapshot", "update")
	EventType string

	// UpdatedAt is the update time
	UpdatedAt Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// AccountUpdate represents account update from WebSocket
type AccountUpdate struct {
	// Balances is the list of currency balances (merged from all pages)
	Balances []*Balance

	// EventType describes the type of event (e.g., "snapshot", "update")
	EventType string

	// UpdatedAt is the update time
	UpdatedAt Timestamp

	// TotalEquity is the total equity in base currency
	TotalEquity Decimal

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// PositionUpdate represents position update from WebSocket
type PositionUpdate struct {
	// Positions is the list of positions (merged from all pages)
	Positions []*Position

	// EventType describes the type of event (e.g., "snapshot", "update")
	EventType string

	// UpdatedAt is the update time
	UpdatedAt Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// OrderUpdate represents order update from WebSocket
type OrderUpdate struct {
	// Orders is the list of orders
	Orders []*Order

	// UpdatedAt is the update time
	UpdatedAt Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// WebSocketSubscribeRequest represents a WebSocket subscription request
type WebSocketSubscribeRequest struct {
	// Channel is the channel to subscribe to
	Channel string

	// Symbols is the list of symbols to subscribe (optional)
	Symbols []string

	// Currencies is the list of currencies to subscribe (optional)
	Currencies []string

	// InstrumentType is the instrument type (spot, futures, swap, etc.)
	InstrumentType InstrumentType

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}

// WebSocketError represents a WebSocket error event
type WebSocketError struct {
	// Event is the event type
	Event string

	// Code is the error code
	Code string

	// Message is the error message
	Message string

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// WebSocketSubscribe represents a WebSocket subscription event
type WebSocketSubscribe struct {
	// Event is the event type
	Event string

	// Channel is the subscribed channel
	Channel string

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// WebSocketUnsubscribe represents a WebSocket unsubscription event
type WebSocketUnsubscribe struct {
	// Event is the event type
	Event string

	// Channel is the unsubscribed channel
	Channel string

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// WebSocketLogin represents a WebSocket login event
type WebSocketLogin struct {
	// Event is the event type
	Event string

	// Code is the response code
	Code string

	// Message is the response message
	Message string

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// WebSocketSuccess represents a WebSocket success event
type WebSocketSuccess struct {
	// Code is the response code
	Code int

	// Message is the response message
	Message string

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// WebSocketSystemMessage represents system-level messages from WebSocket client
type WebSocketSystemMessage struct {
	// Type is the message type (e.g., "connection", "reconnection", "subscription")
	Type string

	// Message is the human-readable message
	Message string

	// Private indicates whether this is for private or public connection
	Private bool

	// Timestamp is when the message was generated
	Timestamp Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// WebSocketSystemError represents system-level errors from WebSocket client
type WebSocketSystemError struct {
	// Type is the error type (e.g., "connection", "sender", "receiver", "subscription")
	Type string

	// Error is the error message
	Error string

	// Private indicates whether this is for private or public connection
	Private bool

	// Timestamp is when the error occurred
	Timestamp Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}
