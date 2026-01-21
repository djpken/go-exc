package market

// Ticker represents ticker information
type Ticker struct {
	Symbol        string `json:"symbol"`
	LastPrice     string `json:"last"`
	QuoteVolume   string `json:"qv_24h"`
	BaseVolume    string `json:"v_24h"`
	HighPrice     string `json:"high_24h"`
	LowPrice      string `json:"low_24h"`
	OpenPrice     string `json:"open_24h"`
	BestBid       string `json:"bid_px"`
	BestBidSize   string `json:"bid_sz"`
	BestAsk       string `json:"ask_px"`
	BestAskSize   string `json:"ask_sz"`
	Timestamp     int64  `json:"ts,string"`
	Fluctuation   string `json:"fluctuation"`
	PercentChange string `json:"percent_change_24h"`
}

// OrderBookItem represents a single order book entry as [price, amount]
type OrderBookItem []string

// OrderBook represents order book data
type OrderBook struct {
	Timestamp int64           `json:"timestamp"`
	Bids      []OrderBookItem `json:"bids"`
	Asks      []OrderBookItem `json:"asks"`
}

// Trade represents a single trade
type Trade struct {
	TradeID   string `json:"trade_id"`
	Price     string `json:"price"`
	Size      string `json:"size"`
	Side      string `json:"side"` // buy or sell
	Timestamp int64  `json:"timestamp"`
}

// Kline represents candlestick data
type Kline struct {
	Timestamp   int64  `json:"timestamp"`
	Open        string `json:"open"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Close       string `json:"close"`
	Volume      string `json:"volume"`
	QuoteVolume string `json:"quote_volume"`
}

// Symbol represents trading pair information
type Symbol struct {
	Symbol          string `json:"symbol"`
	SymbolID        string `json:"symbol_id"`
	BaseCurrency    string `json:"base_currency"`
	QuoteCurrency   string `json:"quote_currency"`
	QuoteIncrement  string `json:"quote_increment"`
	BaseMinSize     string `json:"base_min_size"`
	PriceMinPrecision int  `json:"price_min_precision"`
	PriceMaxPrecision int  `json:"price_max_precision"`
	ExpirationTime  string `json:"expiration_time"`
	MinBuyAmount    string `json:"min_buy_amount"`
	MinSellAmount   string `json:"min_sell_amount"`
	TradingStatus   string `json:"trade_status"`
}
