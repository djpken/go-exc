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
