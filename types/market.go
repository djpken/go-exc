package types

// Ticker represents ticker information
type Ticker struct {
	// Symbol is the trading symbol
	Symbol string

	// LastPrice is the last traded price
	LastPrice Decimal

	// BidPrice is the best bid price
	BidPrice Decimal

	// AskPrice is the best ask price
	AskPrice Decimal

	// High24h is the 24h high price
	High24h Decimal

	// Low24h is the 24h low price
	Low24h Decimal

	// Volume24h is the 24h trading volume
	Volume24h Decimal

	// Timestamp is the ticker timestamp
	Timestamp Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// TickerUpdate represents a ticker update event from WebSocket
type TickerUpdate struct {
	// Symbol is the trading symbol
	Symbol string

	// LastPrice is the last traded price
	LastPrice Decimal

	// BidPrice is the best bid price
	BidPrice Decimal

	// BidSize is the best bid size
	BidSize Decimal

	// AskPrice is the best ask price
	AskPrice Decimal

	// AskSize is the best ask size
	AskSize Decimal

	// High24h is the 24h high price
	High24h Decimal

	// Low24h is the 24h low price
	Low24h Decimal

	// Open24h is the 24h opening price
	Open24h Decimal

	// Volume24h is the 24h trading volume (base currency)
	Volume24h Decimal

	// QuoteVolume24h is the 24h quote volume
	QuoteVolume24h Decimal

	// PriceChange24h is the 24h price change
	PriceChange24h Decimal

	// PercentChange24h is the 24h percent change
	PercentChange24h Decimal

	// Timestamp is the ticker timestamp
	Timestamp Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// CandleUpdate represents a candlestick/kline update event from WebSocket
type CandleUpdate struct {
	// Symbol is the trading symbol
	Symbol string

	// Interval is the candlestick interval (e.g., "1m", "5m", "1H", "1D")
	Interval string

	// Open is the opening price
	Open Decimal

	// High is the highest price
	High Decimal

	// Low is the lowest price
	Low Decimal

	// Close is the closing price
	Close Decimal

	// Volume is the trading volume (base currency)
	Volume Decimal

	// QuoteVolume is the quote currency volume
	QuoteVolume Decimal

	// Timestamp is the candle start time
	Timestamp Timestamp

	// Confirmed indicates if the candle is confirmed (closed)
	// false means the candle is still being formed
	Confirmed bool

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// OrderBook represents an order book
type OrderBook struct {
	// Symbol is the trading symbol
	Symbol string

	// Bids is the list of bid orders (price, quantity)
	Bids []OrderBookLevel

	// Asks is the list of ask orders (price, quantity)
	Asks []OrderBookLevel

	// Timestamp is the order book timestamp
	Timestamp Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// OrderBookLevel represents a single order book level
type OrderBookLevel struct {
	// Price is the price level
	Price Decimal

	// Quantity is the total quantity at this price level
	Quantity Decimal
}

// Candle represents a candlestick/kline (for REST API historical data)
type Candle struct {
	// Symbol is the trading symbol
	Symbol string

	// Interval is the candlestick interval (e.g., "1m", "5m", "1H", "1D")
	Interval string

	// Open is the opening price
	Open Decimal

	// High is the highest price
	High Decimal

	// Low is the lowest price
	Low Decimal

	// Close is the closing price
	Close Decimal

	// Volume is the trading volume (base currency)
	Volume Decimal

	// QuoteVolume is the quote currency volume
	QuoteVolume Decimal

	// Timestamp is the candle start time
	Timestamp Timestamp

	// Confirmed indicates if the candle is confirmed (closed)
	// false means the candle is still being formed (for latest candle)
	Confirmed bool

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// Trade represents a trade
type Trade struct {
	// ID is the trade ID
	ID string

	// Symbol is the trading symbol
	Symbol string

	// Side is the trade side (buy/sell)
	Side string

	// Price is the trade price
	Price Decimal

	// Quantity is the trade quantity
	Quantity Decimal

	// Timestamp is the trade timestamp
	Timestamp Timestamp

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// Kline represents a candlestick
type Kline struct {
	// Symbol is the trading symbol
	Symbol string

	// Interval is the kline interval
	Interval string

	// OpenTime is the opening time
	OpenTime Timestamp

	// CloseTime is the closing time
	CloseTime Timestamp

	// Open is the opening price
	Open Decimal

	// High is the highest price
	High Decimal

	// Low is the lowest price
	Low Decimal

	// Close is the closing price
	Close Decimal

	// Volume is the trading volume
	Volume Decimal

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// Instrument represents a trading instrument/symbol
type Instrument struct {
	// Symbol is the trading symbol (e.g., "BTC_USDT", "BTC-USDT")
	Symbol string

	// BaseCurrency is the base currency (e.g., "BTC" in "BTC_USDT")
	BaseCurrency string

	// QuoteCurrency is the quote currency (e.g., "USDT" in "BTC_USDT")
	QuoteCurrency string

	// InstrumentType is the type of instrument (spot, futures, swap, etc.)
	InstrumentType InstrumentType

	// Status is the trading status (trading, halt, etc.)
	Status string

	// MinOrderSize is the minimum order size
	MinOrderSize Decimal

	// MaxOrderSize is the maximum order size
	MaxOrderSize Decimal

	// PricePrecision is the price decimal precision
	PricePrecision Decimal

	// QuantityPrecision is the quantity decimal precision
	QuantityPrecision Decimal

	// LastPrice is the last trade price
	LastPrice Decimal
	// Extra contains exchange-specific fields
	Extra    map[string]interface{}
	MaxLever int
	CtVal    Decimal
}
