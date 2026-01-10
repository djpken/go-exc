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
