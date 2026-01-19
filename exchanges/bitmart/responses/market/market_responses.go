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
type TickersResponse struct {
	BaseResponse
	Data []market.Ticker `json:"data"`
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
type TradesResponse struct {
	BaseResponse
	Data []market.Trade `json:"data"`
}

// KlineResponse represents kline API response
type KlineResponse struct {
	BaseResponse
	Data []market.Kline `json:"data"`
}

// SymbolsResponse represents symbols API response
type SymbolsResponse struct {
	BaseResponse
	Data struct {
		Symbols []market.Symbol `json:"symbols"`
	} `json:"data"`
}

// SymbolDetailResponse represents symbol detail API response
type SymbolDetailResponse struct {
	BaseResponse
	Data market.Symbol `json:"data"`
}
