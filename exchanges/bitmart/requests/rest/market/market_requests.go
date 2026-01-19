package market

// GetTickerRequest represents request for getting ticker
type GetTickerRequest struct {
	Symbol string `json:"symbol"`
}

// GetTickersRequest represents request for getting all tickers
type GetTickersRequest struct {
	// Empty struct - gets all tickers
}

// GetOrderBookRequest represents request for getting order book
type GetOrderBookRequest struct {
	Symbol string `json:"symbol"`
	Depth  int    `json:"depth,omitempty"` // Optional: 5, 20, 50, 100, 200
}

// GetTradesRequest represents request for getting recent trades
type GetTradesRequest struct {
	Symbol string `json:"symbol"`
	Limit  int    `json:"limit,omitempty"` // Optional: default 50, max 50
}

// GetKlineRequest represents request for getting candlestick data
type GetKlineRequest struct {
	Symbol    string `json:"symbol"`
	FromTime  int64  `json:"from,omitempty"`  // Start timestamp (seconds)
	ToTime    int64  `json:"to,omitempty"`    // End timestamp (seconds)
	Step      int    `json:"step"`            // Kline step (1, 3, 5, 15, 30, 45, 60, 120, 180, 240, 1440, 10080, 43200)
	Limit     int    `json:"limit,omitempty"` // Optional: max 500
}

// GetSymbolsRequest represents request for getting trading pairs
type GetSymbolsRequest struct {
	// Empty struct - gets all symbols
}

// GetSymbolDetailRequest represents request for getting symbol details
type GetSymbolDetailRequest struct {
	Symbol string `json:"symbol"`
}
