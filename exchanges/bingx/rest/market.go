package rest

import "fmt"

// Market provides BingX market data endpoints
type Market struct {
	client *ClientRest
}

func NewMarket(c *ClientRest) *Market { return &Market{client: c} }

// TickerData is the data field for a single ticker response
type TickerData struct {
	Symbol             string `json:"symbol"`
	LastPrice          string `json:"lastPrice"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenPrice          string `json:"openPrice"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	Count              int64  `json:"count"`
}

// TickerResponse is the full API response for a single ticker
type TickerResponse struct {
	Code int        `json:"code"`
	Data TickerData `json:"data"`
}

// TickersResponse is the full API response for all tickers
type TickersResponse struct {
	Code int          `json:"code"`
	Data []TickerData `json:"data"`
}

// GetTicker retrieves 24hr ticker statistics for a symbol
// GET /openApi/swap/v2/quote/ticker
func (m *Market) GetTicker(symbol string) (*TickerResponse, error) {
	var result TickerResponse
	params := map[string]string{"symbol": symbol}
	if err := m.client.GETPublic("/openApi/swap/v2/quote/ticker", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTickers retrieves 24hr ticker statistics for all symbols
// GET /openApi/swap/v2/quote/ticker (no symbol param)
func (m *Market) GetTickers() (*TickersResponse, error) {
	var result TickersResponse
	if err := m.client.GETPublic("/openApi/swap/v2/quote/ticker", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OrderBookData is the depth data
type OrderBookData struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
	T    int64      `json:"T"`
}

// OrderBookResponse is the full API response for order book depth
type OrderBookResponse struct {
	Code int           `json:"code"`
	Data OrderBookData `json:"data"`
}

// GetOrderBook retrieves order book depth
// GET /openApi/swap/v2/quote/depth
func (m *Market) GetOrderBook(symbol string, limit int) (*OrderBookResponse, error) {
	var result OrderBookResponse
	params := map[string]string{"symbol": symbol}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	if err := m.client.GETPublic("/openApi/swap/v2/quote/depth", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// KlineEntry represents a single candlestick
type KlineEntry struct {
	Open   string `json:"open"`
	Close  string `json:"close"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Volume string `json:"volume"`
	Time   int64  `json:"time"`
}

// KlinesResponse is the full API response for klines
type KlinesResponse struct {
	Code int          `json:"code"`
	Data []KlineEntry `json:"data"`
}

// GetKlines retrieves candlestick/kline data (v3)
// GET /openApi/swap/v3/quote/klines
// timeZone: 0 = UTC+0, 8 = UTC+8 (default 8)
func (m *Market) GetKlines(symbol, interval string, startTime, endTime int64, limit int, timeZone int32) (*KlinesResponse, error) {
	var result KlinesResponse
	params := map[string]string{
		"symbol":   symbol,
		"interval": interval,
		"timeZone": fmt.Sprintf("%d", timeZone),
	}
	if startTime > 0 {
		params["startTime"] = fmt.Sprintf("%d", startTime)
	}
	if endTime > 0 {
		params["endTime"] = fmt.Sprintf("%d", endTime)
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	if err := m.client.GETPublic("/openApi/swap/v3/quote/klines", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ContractInfo holds information about a single perpetual swap contract
type ContractInfo struct {
	ContractID        string  `json:"contractId"`
	Symbol            string  `json:"symbol"`
	Size              string  `json:"size"`
	QuantityPrecision int     `json:"quantityPrecision"`
	PricePrecision    int     `json:"pricePrecision"`
	FeeRate           float64 `json:"feeRate"`
	MakerFeeRate      float64 `json:"makerFeeRate"`
	TakerFeeRate      float64 `json:"takerFeeRate"`
	TradeMinLimit     float64 `json:"tradeMinLimit"`
	TradeMinQuantity  float64 `json:"tradeMinQuantity"`
	TradeMinUSDT      float64 `json:"tradeMinUSDT"`
	Currency          string  `json:"currency"`
	Asset             string  `json:"asset"`
	Status            int     `json:"status"`
	APIStateOpen      string  `json:"apiStateOpen"`
	APIStateClose     string  `json:"apiStateClose"`
	EnsureTrigger     bool    `json:"ensureTrigger"`
	TriggerFeeRate    string  `json:"triggerFeeRate"`
	BrokerState       bool    `json:"brokerState"`
	LaunchTime        int64   `json:"launchTime"`
	MaintainTime      int64   `json:"maintainTime"`
	OffTime           int64   `json:"offTime"`
	DisplayName       string  `json:"displayName"`
}

// ContractsResponse is the full API response for contract info
type ContractsResponse struct {
	Code int            `json:"code"`
	Data []ContractInfo `json:"data"`
}

// GetContracts retrieves all perpetual swap contract specifications
// GET /openApi/swap/v2/quote/contracts
func (m *Market) GetContracts() (*ContractsResponse, error) {
	var result ContractsResponse
	if err := m.client.GETPublic("/openApi/swap/v2/quote/contracts", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
