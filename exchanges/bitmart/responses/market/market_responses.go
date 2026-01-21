package market

import "github.com/djpken/go-exc/exchanges/bitmart/models/market"

// BaseResponse represents the base API response structure
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

// TickerResponse represents ticker API response
type TickerResponse struct {
	BaseResponse
	Data market.Ticker `json:"data"`
}

// TickersResponse represents tickers API response
// BitMart returns array of arrays: [[symbol, last, v_24h, qv_24h, open_24h, high_24h, low_24h, fluctuation, bid_px, bid_sz, ask_px, ask_sz, ts], ...]
type TickersResponse struct {
	BaseResponse
	Data [][]interface{} `json:"data"`
}

// OrderBookResponse represents order book API response
type OrderBookResponse struct {
	BaseResponse
	Data struct {
		Timestamp int64                 `json:"timestamp"`
		Bids      []market.OrderBookItem `json:"bids"`
		Asks      []market.OrderBookItem `json:"asks"`
	} `json:"data"`
}

// TradesResponse represents trades API response
// BitMart returns array of arrays: [[symbol, timestamp, price, size, side], ...]
type TradesResponse struct {
	BaseResponse
	Data [][]interface{} `json:"data"`
}

// KlineResponse represents kline API response
// BitMart returns array of arrays: [[timestamp, open, high, low, close, volume, quote_volume], ...]
type KlineResponse struct {
	BaseResponse
	Data [][]interface{} `json:"data"`
}

// SymbolsResponse represents symbols API response
// For /spot/v1/symbols - returns simple string array
type SymbolsResponse struct {
	BaseResponse
	Data struct {
		Symbols []string `json:"symbols"`
	} `json:"data"`
}

// SymbolDetailResponse represents symbol detail API response
type SymbolDetailResponse struct {
	BaseResponse
	Data market.Symbol `json:"data"`
}

// SymbolsDetailsResponse represents symbols details API response
// For /spot/v1/symbols/details - returns detailed Symbol objects
type SymbolsDetailsResponse struct {
	BaseResponse
	Data struct {
		Symbols []market.Symbol `json:"symbols"`
	} `json:"data"`
}
