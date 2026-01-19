package ws

// SubscribeRequest represents WebSocket subscribe request
type SubscribeRequest struct {
	Op   string   `json:"op"`   // subscribe or unsubscribe
	Args []string `json:"args"` // channel list
}

// LoginRequest represents WebSocket login request for private channels
type LoginRequest struct {
	Op   string   `json:"op"`   // login
	Args []string `json:"args"` // [api_key, timestamp, sign]
}

// Public channel request types

// TickerRequest represents ticker channel subscription
type TickerRequest struct {
	Symbol string `json:"symbol"` // e.g., "BTC_USDT"
}

// DepthRequest represents depth channel subscription
type DepthRequest struct {
	Symbol string `json:"symbol"` // e.g., "BTC_USDT"
	Depth  int    `json:"depth"`  // 5, 20, 50
}

// TradeRequest represents trade channel subscription
type TradeRequest struct {
	Symbol string `json:"symbol"` // e.g., "BTC_USDT"
}

// KlineRequest represents kline channel subscription
type KlineRequest struct {
	Symbol string `json:"symbol"` // e.g., "BTC_USDT"
	Step   int    `json:"step"`   // 1, 3, 5, 15, 30, 60, etc.
}

// Private channel request types (require authentication)

// OrderRequest represents order channel subscription
type OrderRequest struct {
	// Private channel - no additional params needed
}

// BalanceRequest represents balance channel subscription
type BalanceRequest struct {
	// Private channel - no additional params needed
}

// TradeUpdateRequest represents trade update channel subscription
type TradeUpdateRequest struct {
	// Private channel - no additional params needed
}
