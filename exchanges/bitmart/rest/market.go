package rest

import (
	"fmt"
	"strconv"

	"github.com/djpken/go-exc/exchanges/bitmart/requests/rest/market"
	responses "github.com/djpken/go-exc/exchanges/bitmart/responses/market"
)

// Market provides access to BitMart market data API
type Market struct {
	client *ClientRest
}

// NewMarket creates a new Market API instance
func NewMarket(c *ClientRest) *Market {
	return &Market{client: c}
}

// GetTicker retrieves ticker information for a specific symbol
//
// API: GET /spot/quotation/v3/ticker
func (m *Market) GetTicker(req market.GetTickerRequest) (*responses.TickerResponse, error) {
	endpoint := fmt.Sprintf("/spot/quotation/v3/ticker?symbol=%s", req.Symbol)

	var result responses.TickerResponse
	if err := m.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTickers retrieves ticker information for all symbols
//
// API: GET /spot/quotation/v3/tickers
func (m *Market) GetTickers() (*responses.TickersResponse, error) {
	endpoint := "/spot/quotation/v3/tickers"

	var result responses.TickersResponse
	if err := m.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrderBook retrieves order book for a specific symbol
//
// API: GET /spot/quotation/v3/books
func (m *Market) GetOrderBook(req market.GetOrderBookRequest) (*responses.OrderBookResponse, error) {
	endpoint := fmt.Sprintf("/spot/quotation/v3/books?symbol=%s", req.Symbol)

	if req.Depth > 0 {
		endpoint += fmt.Sprintf("&depth=%d", req.Depth)
	}

	var result responses.OrderBookResponse
	if err := m.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTrades retrieves recent trades for a specific symbol
//
// API: GET /spot/quotation/v3/trades
func (m *Market) GetTrades(req market.GetTradesRequest) (*responses.TradesResponse, error) {
	endpoint := fmt.Sprintf("/spot/quotation/v3/trades?symbol=%s", req.Symbol)

	if req.Limit > 0 {
		endpoint += fmt.Sprintf("&limit=%d", req.Limit)
	}

	var result responses.TradesResponse
	if err := m.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetKlines retrieves candlestick/kline data for a specific symbol
//
// API: GET /spot/quotation/v3/klines
func (m *Market) GetKlines(req market.GetKlineRequest) (*responses.KlineResponse, error) {
	endpoint := fmt.Sprintf("/spot/quotation/v3/klines?symbol=%s&step=%d",
		req.Symbol, req.Step)

	if req.FromTime > 0 {
		endpoint += fmt.Sprintf("&from=%s", strconv.FormatInt(req.FromTime, 10))
	}

	if req.ToTime > 0 {
		endpoint += fmt.Sprintf("&to=%s", strconv.FormatInt(req.ToTime, 10))
	}

	if req.Limit > 0 {
		endpoint += fmt.Sprintf("&limit=%d", req.Limit)
	}

	var result responses.KlineResponse
	if err := m.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSymbols retrieves all available trading pairs
//
// API: GET /spot/v1/symbols
func (m *Market) GetSymbols() (*responses.SymbolsResponse, error) {
	endpoint := "/spot/v1/symbols"

	var result responses.SymbolsResponse
	if err := m.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSymbolDetail retrieves details for a specific trading pair
//
// API: GET /spot/v1/symbols/details
func (m *Market) GetSymbolDetail(req market.GetSymbolDetailRequest) (*responses.SymbolDetailResponse, error) {
	endpoint := fmt.Sprintf("/spot/v1/symbols/details?symbol=%s", req.Symbol)

	var result responses.SymbolDetailResponse
	if err := m.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
