package public

// TickerEvent represents ticker WebSocket event
type TickerEvent struct {
	Symbol        string `json:"symbol"`
	LastPrice     string `json:"last_price"`
	QuoteVolume   string `json:"quote_volume_24h"`
	BaseVolume    string `json:"base_volume_24h"`
	HighPrice     string `json:"high_24h"`
	LowPrice      string `json:"low_24h"`
	OpenPrice     string `json:"open_24h"`
	Close         string `json:"close_24h"`
	BestBid       string `json:"best_bid"`
	BestBidSize   string `json:"best_bid_size"`
	BestAsk       string `json:"best_ask"`
	BestAskSize   string `json:"best_ask_size"`
	Timestamp     int64  `json:"timestamp"`
	PriceChange   string `json:"price_change_24h"`
	PercentChange string `json:"percent_change_24h"`
}

// DepthEvent represents order book depth WebSocket event
type DepthEvent struct {
	Symbol    string     `json:"symbol"`
	Timestamp int64      `json:"timestamp"`
	Bids      [][]string `json:"bids"` // [price, amount]
	Asks      [][]string `json:"asks"` // [price, amount]
}

// TradeEvent represents trade WebSocket event
type TradeEvent struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"` // buy or sell
	Timestamp int64  `json:"timestamp"`
	TradeID   string `json:"trade_id"`
}

// KlineEvent represents candlestick WebSocket event
type KlineEvent struct {
	Symbol      string `json:"symbol"`
	Timestamp   int64  `json:"timestamp"`
	Open        string `json:"open"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Close       string `json:"close"`
	Volume      string `json:"volume"`
	QuoteVolume string `json:"quote_volume"`
}

// FuturesTickerEvent represents futures ticker WebSocket event
// This matches the actual BitMart futures ticker response format
type FuturesTickerEvent struct {
	Group string             `json:"group"` // e.g., "futures/ticker:ETHUSDT"
	Data  FuturesTickerData  `json:"data"`
}

// FuturesTickerData represents the data field in futures ticker event
type FuturesTickerData struct {
	Symbol     string `json:"symbol"`      // Contract symbol (e.g., "ETHUSDT")
	LastPrice  string `json:"last_price"`  // Last traded price
	Volume24   string `json:"volume_24"`   // 24h trading volume
	Range      string `json:"range"`       // 24h price change ratio
	MarkPrice  string `json:"mark_price"`  // Mark price
	IndexPrice string `json:"index_price"` // Index price
	AskPrice   string `json:"ask_price"`   // Best ask price
	AskVol     string `json:"ask_vol"`     // Best ask volume
	BidPrice   string `json:"bid_price"`   // Best bid price
	BidVol     string `json:"bid_vol"`     // Best bid volume
}
