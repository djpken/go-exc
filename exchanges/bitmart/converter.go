package bitmart

import (
	"fmt"
	"strconv"
	"time"

	accountmodels "github.com/djpken/go-exc/exchanges/bitmart/models/account"
	"github.com/djpken/go-exc/exchanges/bitmart/models/contract"
	marketmodels "github.com/djpken/go-exc/exchanges/bitmart/models/market"
	trademodels "github.com/djpken/go-exc/exchanges/bitmart/models/trade"
	bitmarttypes "github.com/djpken/go-exc/exchanges/bitmart/types"
	commontypes "github.com/djpken/go-exc/types"
)

// Converter converts between BitMart-specific types and common types
type Converter struct{}

// NewConverter creates a new converter instance
func NewConverter() *Converter {
	return &Converter{}
}

// ConvertOrder converts BitMart order to common order type
func (c *Converter) ConvertOrder(order *trademodels.Order) *commontypes.Order {
	if order == nil {
		return nil
	}

	return &commontypes.Order{
		ID:             order.OrderID,
		Symbol:         order.Symbol,
		Side:           order.Side,
		Type:           order.Type,
		Quantity:       commontypes.Decimal(order.Size),
		Price:          commontypes.Decimal(order.Price),
		FilledQuantity: commontypes.Decimal(order.FilledSize),
		Status:         c.ConvertOrderStatus(order.Status),
		CreatedAt:      commontypes.Timestamp(time.Unix(0, order.CreateTime*int64(time.Millisecond))),
		UpdatedAt:      commontypes.Timestamp(time.Unix(0, order.UpdateTime*int64(time.Millisecond))),
		Extra:          map[string]interface{}{"notional": order.Notional},
	}
}

// ConvertOrderDetail converts BitMart order detail to common order type
func (c *Converter) ConvertOrderDetail(detail *trademodels.OrderDetail) *commontypes.Order {
	if detail == nil {
		return nil
	}

	order := c.ConvertOrder(&detail.Order)
	if order != nil {
		order.Fee = commontypes.Decimal(detail.Fee)
		order.FeeCurrency = detail.FeeCurrency
		if order.Extra == nil {
			order.Extra = make(map[string]interface{})
		}
		order.Extra["avg_price"] = detail.AvgPrice
		order.Extra["client_order_id"] = detail.ClientOrderID
	}
	return order
}

// ConvertBalance converts BitMart balance to common balance type
func (c *Converter) ConvertBalance(balance *accountmodels.Balance) *commontypes.Balance {
	if balance == nil {
		return nil
	}

	// Calculate total from available and unavailable if Total is not set
	total := ""
	if balance.Available != "" && balance.UnAvailable != "" {
		// Total = Available + UnAvailable (all frozen amounts)
		availFloat, _ := commontypes.Decimal(balance.Available).Float64()
		unavailFloat, _ := commontypes.Decimal(balance.UnAvailable).Float64()
		total = strconv.FormatFloat(availFloat+unavailFloat, 'f', -1, 64)
	}

	return &commontypes.Balance{
		Currency:  balance.Currency,
		Available: commontypes.Decimal(balance.Available),
		Frozen:    commontypes.Decimal(balance.Frozen),
		Total:     commontypes.Decimal(total),
		Extra: map[string]interface{}{
			"name":                    balance.Name,
			"available_usd_valuation": balance.AvailableUsdValuation,
			"unAvailable":             balance.UnAvailable,
		},
	}
}

// ConvertAccountBalance converts BitMart balances to common account balance
func (c *Converter) ConvertAccountBalance(balances []accountmodels.Balance) *commontypes.AccountBalance {
	if balances == nil {
		return nil
	}

	commonBalances := make([]*commontypes.Balance, 0, len(balances))
	totalEquity := 0.0

	for _, bal := range balances {
		// Calculate total from available and unavailable if Total is not set
		total := ""
		if bal.Available != "" && bal.UnAvailable != "" {
			// Total = Available + UnAvailable (all frozen amounts)
			availFloat, _ := commontypes.Decimal(bal.Available).Float64()
			unavailFloat, _ := commontypes.Decimal(bal.UnAvailable).Float64()
			total = strconv.FormatFloat(availFloat+unavailFloat, 'f', -1, 64)
		}

		// Calculate total equity from USD valuations
		if bal.AvailableUsdValuation != "" {
			usdValue, err := strconv.ParseFloat(bal.AvailableUsdValuation, 64)
			if err == nil {
				totalEquity += usdValue
			}
		}

		commonBalances = append(commonBalances, &commontypes.Balance{
			Currency:  bal.Currency,
			Available: commontypes.Decimal(bal.Available),
			Frozen:    commontypes.Decimal(bal.Frozen),
			Total:     commontypes.Decimal(total),
			Extra: map[string]interface{}{
				"name":                    bal.Name,
				"available_usd_valuation": bal.AvailableUsdValuation,
				"unAvailable":             bal.UnAvailable,
			},
		})
	}

	return &commontypes.AccountBalance{
		Balances:    commonBalances,
		TotalEquity: commontypes.Decimal(strconv.FormatFloat(totalEquity, 'f', -1, 64)),
	}
}

// ConvertTicker converts BitMart ticker to common ticker type
func (c *Converter) ConvertTicker(ticker *marketmodels.Ticker) *commontypes.Ticker {
	if ticker == nil {
		return nil
	}

	return &commontypes.Ticker{
		Symbol:    ticker.Symbol,
		LastPrice: commontypes.Decimal(ticker.LastPrice),
		High24h:   commontypes.Decimal(ticker.HighPrice),
		Low24h:    commontypes.Decimal(ticker.LowPrice),
		Volume24h: commontypes.Decimal(ticker.BaseVolume),
		Timestamp: commontypes.Timestamp(time.Unix(0, ticker.Timestamp*int64(time.Millisecond))),
		Extra: map[string]interface{}{
			"quote_volume":   ticker.QuoteVolume,
			"open":           ticker.OpenPrice,
			"best_bid":       ticker.BestBid,
			"best_ask":       ticker.BestAsk,
			"best_bid_size":  ticker.BestBidSize,
			"best_ask_size":  ticker.BestAskSize,
			"fluctuation":    ticker.Fluctuation,
			"percent_change": ticker.PercentChange,
		},
	}
}

// ConvertOrderBook converts BitMart order book to common order book type
func (c *Converter) ConvertOrderBook(ob interface{}, symbol string) *commontypes.OrderBook {
	if ob == nil {
		return nil
	}

	// Handle different order book types
	type orderBookData struct {
		Timestamp int64
		Bids      []marketmodels.OrderBookItem
		Asks      []marketmodels.OrderBookItem
	}

	var data orderBookData
	switch v := ob.(type) {
	case *marketmodels.OrderBook:
		data.Timestamp = v.Timestamp
		data.Bids = v.Bids
		data.Asks = v.Asks
	case *struct {
		Timestamp int64                        `json:"timestamp"`
		Bids      []marketmodels.OrderBookItem `json:"bids"`
		Asks      []marketmodels.OrderBookItem `json:"asks"`
	}:
		data.Timestamp = v.Timestamp
		data.Bids = v.Bids
		data.Asks = v.Asks
	default:
		return nil
	}

	// Convert bids
	bids := make([]commontypes.OrderBookLevel, len(data.Bids))
	for i, bid := range data.Bids {
		if len(bid) >= 2 {
			bids[i] = commontypes.OrderBookLevel{
				Price:    commontypes.Decimal(bid[0]),
				Quantity: commontypes.Decimal(bid[1]),
			}
		}
	}

	// Convert asks
	asks := make([]commontypes.OrderBookLevel, len(data.Asks))
	for i, ask := range data.Asks {
		if len(ask) >= 2 {
			asks[i] = commontypes.OrderBookLevel{
				Price:    commontypes.Decimal(ask[0]),
				Quantity: commontypes.Decimal(ask[1]),
			}
		}
	}

	return &commontypes.OrderBook{
		Symbol:    symbol,
		Bids:      bids,
		Asks:      asks,
		Timestamp: commontypes.Timestamp(time.Unix(0, data.Timestamp*int64(time.Millisecond))),
	}
}

// ConvertOrderStatus converts BitMart order status to common status
func (c *Converter) ConvertOrderStatus(status string) string {
	switch bitmarttypes.OrderStatus(status) {
	case bitmarttypes.OrderStatusNew:
		return string(commontypes.OrderStatusOpen)
	case bitmarttypes.OrderStatusPartiallyFilled:
		return string(commontypes.OrderStatusPartiallyFilled)
	case bitmarttypes.OrderStatusFilled:
		return string(commontypes.OrderStatusFilled)
	case bitmarttypes.OrderStatusCanceled:
		return string(commontypes.OrderStatusCanceled)
	case bitmarttypes.OrderStatusPendingCancel:
		return "canceling"
	default:
		return status
	}
}

// ConvertOrderSide converts common side to BitMart side
func (c *Converter) ConvertOrderSide(side string) string {
	switch side {
	case "buy", "BUY":
		return string(bitmarttypes.OrderSideBuy)
	case "sell", "SELL":
		return string(bitmarttypes.OrderSideSell)
	default:
		return side
	}
}

// ConvertOrderType converts common order type to BitMart order type
func (c *Converter) ConvertOrderType(orderType string) string {
	switch orderType {
	case "limit", "LIMIT":
		return string(bitmarttypes.OrderTypeLimit)
	case "market", "MARKET":
		return string(bitmarttypes.OrderTypeMarket)
	case "limit_maker", "LIMIT_MAKER":
		return string(bitmarttypes.OrderTypeLimitMaker)
	case "ioc", "IOC":
		return string(bitmarttypes.OrderTypeIOC)
	default:
		return orderType
	}
}

// formatFloat converts float64 to string
func (c *Converter) formatFloat(f float64) string {
	return fmt.Sprintf("%f", f)
}

// formatFloatPrecision converts float64 to string with specific precision
func (c *Converter) formatFloatPrecision(f float64, precision int) string {
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, f)
}

// ConvertInstrument converts BitMart symbol to common Instrument type
func (c *Converter) ConvertInstrument(symbol *contract.ContractDetail) *commontypes.Instrument {
	if symbol == nil {
		return nil
	}

	// BitMart is primarily spot exchange
	instType := commontypes.InstrumentSwap

	maxLever, err := strconv.Atoi(symbol.MaxLeverage)
	if err != nil {
		return nil
	}

	return &commontypes.Instrument{
		Symbol:            symbol.Symbol,
		BaseCurrency:      symbol.BaseCurrency,
		QuoteCurrency:     symbol.QuoteCurrency,
		CtVal:             commontypes.Decimal(symbol.ContractSize),
		InstrumentType:    instType,
		Status:            symbol.Status,
		MinOrderSize:      commontypes.Decimal(symbol.MinVolume),
		MaxOrderSize:      commontypes.Decimal(symbol.MaxVolume),
		PricePrecision:    commontypes.Decimal(symbol.PricePrecision),
		QuantityPrecision: commontypes.Decimal(symbol.VolPrecision),
		LastPrice:         commontypes.Decimal(symbol.LastPrice),
		MaxLever:          maxLever,
	}
}

// ConvertTickersResponse converts BitMart tickers array response to common Ticker types
// BitMart returns array of arrays: [[symbol, last, v_24h, qv_24h, open_24h, high_24h, low_24h, ...], ...]
func (c *Converter) ConvertTickersResponse(tickersData [][]interface{}) []*commontypes.Ticker {
	tickers := make([]*commontypes.Ticker, 0, len(tickersData))

	for _, tickerArray := range tickersData {
		if len(tickerArray) < 13 {
			continue
		}

		// Extract fields from array
		// Format: [symbol, last, v_24h, qv_24h, open_24h, high_24h, low_24h, fluctuation, bid_px, bid_sz, ask_px, ask_sz, ts]
		symbol, _ := tickerArray[0].(string)
		last, _ := tickerArray[1].(string)
		volume24h, _ := tickerArray[2].(string)
		quoteVolume24h, _ := tickerArray[3].(string)
		open24h, _ := tickerArray[4].(string)
		high24h, _ := tickerArray[5].(string)
		low24h, _ := tickerArray[6].(string)
		fluctuation, _ := tickerArray[7].(string)
		bidPx, _ := tickerArray[8].(string)
		// bidSz, _ := tickerArray[9].(string)
		askPx, _ := tickerArray[10].(string)
		// askSz, _ := tickerArray[11].(string)
		tsFloat, _ := tickerArray[12].(float64)

		ticker := &commontypes.Ticker{
			Symbol:    symbol,
			LastPrice: commontypes.Decimal(last),
			BidPrice:  commontypes.Decimal(bidPx),
			AskPrice:  commontypes.Decimal(askPx),
			High24h:   commontypes.Decimal(high24h),
			Low24h:    commontypes.Decimal(low24h),
			Volume24h: commontypes.Decimal(volume24h),
			Timestamp: commontypes.Timestamp(time.Unix(0, int64(tsFloat)*int64(time.Millisecond))),
			Extra: map[string]interface{}{
				"quote_volume_24h": quoteVolume24h,
				"open_24h":         open24h,
				"fluctuation":      fluctuation,
			},
		}

		tickers = append(tickers, ticker)
	}

	return tickers
}

// ConvertIntervalToStep converts common interval format to BitMart step (minutes)
func (c *Converter) ConvertIntervalToStep(interval string) (int, error) {
	// Common intervals: "1m", "5m", "15m", "30m", "1H", "4H", "1D", "1W", "1M"
	// BitMart steps (in minutes): 1, 3, 5, 15, 30, 45, 60, 120, 180, 240, 1440, 10080, 43200

	intervalMap := map[string]int{
		"1m":  1,
		"3m":  3,
		"5m":  5,
		"15m": 15,
		"30m": 30,
		"45m": 45,
		"1H":  60,
		"2H":  120,
		"3H":  180,
		"4H":  240,
		"1D":  1440,
		"1W":  10080,
		"1M":  43200,
	}

	step, ok := intervalMap[interval]
	if !ok {
		return 0, fmt.Errorf("unsupported interval: %s", interval)
	}

	return step, nil
}

// ConvertKlineArrayToCandle converts BitMart kline array to common Candle type
// BitMart format: [timestamp, open, high, low, close, volume, quote_volume]
func (c *Converter) ConvertKlineArrayToCandle(klineData []interface{}, symbol, interval string) *commontypes.Candle {
	if len(klineData) < 7 {
		return nil
	}

	// Parse timestamp
	timestamp, ok := klineData[0].(float64)
	if !ok {
		return nil
	}

	// Parse OHLCV data
	open, _ := klineData[1].(string)
	high, _ := klineData[2].(string)
	low, _ := klineData[3].(string)
	close_, _ := klineData[4].(string)
	volume, _ := klineData[5].(string)
	quoteVolume, _ := klineData[6].(string)

	return &commontypes.Candle{
		Symbol:      symbol,
		Interval:    interval,
		Open:        commontypes.Decimal(open),
		High:        commontypes.Decimal(high),
		Low:         commontypes.Decimal(low),
		Close:       commontypes.Decimal(close_),
		Volume:      commontypes.Decimal(volume),
		QuoteVolume: commontypes.Decimal(quoteVolume),
		Timestamp:   commontypes.Timestamp(time.Unix(int64(timestamp), 0)),
		Confirmed:   true, // Historical candles are always confirmed
		Extra:       make(map[string]interface{}),
	}
}
