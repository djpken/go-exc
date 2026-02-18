package bitmart

import (
	"fmt"
	"strconv"
	"time"

	accountmodels "github.com/djpken/go-exc/exchanges/bitmart/models/account"
	"github.com/djpken/go-exc/exchanges/bitmart/models/contract"
	marketmodels "github.com/djpken/go-exc/exchanges/bitmart/models/market"
	trademodels "github.com/djpken/go-exc/exchanges/bitmart/models/trade"
	contractresponses "github.com/djpken/go-exc/exchanges/bitmart/responses/contract"
	bitmarttypes "github.com/djpken/go-exc/exchanges/bitmart/types"
	commontypes "github.com/djpken/go-exc/types"
)

// Converter converts between BitMart-specific types and common types
type Converter struct{}

// NewConverter creates a new converter instance
func NewConverter() *Converter {
	return &Converter{}
}

// stringToDecimal converts a string to Decimal
// Returns ZeroDecimal if the string is empty or invalid
func (c *Converter) stringToDecimal(s string) commontypes.Decimal {
	if s == "" {
		return commontypes.ZeroDecimal
	}
	d, err := commontypes.NewDecimal(s)
	if err != nil {
		// If parsing fails, return zero (for robustness)
		return commontypes.ZeroDecimal
	}
	return d
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
		Quantity:       c.stringToDecimal(order.Size),
		Price:          c.stringToDecimal(order.Price),
		FilledQuantity: c.stringToDecimal(order.FilledSize),
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
		order.Fee = c.stringToDecimal(detail.Fee)
		order.FeeCurrency = detail.FeeCurrency
		if order.Extra == nil {
			order.Extra = make(map[string]interface{})
		}
		order.Extra["avg_price"] = detail.AvgPrice
		order.Extra["client_order_id"] = detail.ClientOrderID
	}
	return order
}

// ConvertContractTrades converts BitMart contract trades to common order type
// This aggregates trade executions into an Order object
func (c *Converter) ConvertContractTrades(trades []contractresponses.ContractTrade) *commontypes.Order {
	if len(trades) == 0 {
		return nil
	}

	// Use the first trade for basic info
	firstTrade := trades[0]

	// Calculate totals from all trades
	var totalVolume, totalFees, totalRealisedProfit float64
	var avgPrice float64
	var latestTime int64

	for _, trade := range trades {
		vol, _ := strconv.ParseFloat(trade.Vol, 64)
		fee, _ := strconv.ParseFloat(trade.PaidFees, 64)
		profit, _ := strconv.ParseFloat(trade.RealisedProfit, 64)
		price, _ := strconv.ParseFloat(trade.Price, 64)

		totalVolume += vol
		totalFees += fee
		totalRealisedProfit += profit
		avgPrice += price * vol // Weighted by volume

		if trade.CreateTime > latestTime {
			latestTime = trade.CreateTime
		}
	}

	// Calculate weighted average price
	if totalVolume > 0 {
		avgPrice = avgPrice / totalVolume
	}

	// Convert side (contract side is different from spot)
	side := "buy"
	if firstTrade.Side == 3 || firstTrade.Side == 4 {
		side = "sell"
	}

	return &commontypes.Order{
		ID:             firstTrade.OrderID,
		Symbol:         firstTrade.Symbol,
		Side:           side,
		Type:           "limit", // Default, actual type not provided in trades API
		Quantity:       c.stringToDecimal(strconv.FormatFloat(totalVolume, 'f', -1, 64)),
		Price:          c.stringToDecimal(strconv.FormatFloat(avgPrice, 'f', -1, 64)),
		FilledQuantity: c.stringToDecimal(strconv.FormatFloat(totalVolume, 'f', -1, 64)),
		Status:         "filled", // If trades exist, order is at least partially filled
		Fee:            c.stringToDecimal(strconv.FormatFloat(totalFees, 'f', -1, 64)),
		FeeCurrency:    "USDT", // Default, actual currency not provided
		CreatedAt:      commontypes.Timestamp(time.UnixMilli(firstTrade.CreateTime)),
		UpdatedAt:      commontypes.Timestamp(time.UnixMilli(latestTime)),
		Extra: map[string]interface{}{
			"account":         firstTrade.Account,
			"realised_profit": strconv.FormatFloat(totalRealisedProfit, 'f', -1, 64),
			"trade_count":     len(trades),
			"contract_side":   firstTrade.Side, // Original contract side value
			"avg_price":       strconv.FormatFloat(avgPrice, 'f', -1, 64),
		},
	}
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
		availFloat, _ := c.stringToDecimal(balance.Available).Float64()
		unavailFloat, _ := c.stringToDecimal(balance.UnAvailable).Float64()
		total = strconv.FormatFloat(availFloat+unavailFloat, 'f', -1, 64)
	}

	return &commontypes.Balance{
		Currency:  balance.Currency,
		Available: c.stringToDecimal(balance.Available),
		Frozen:    c.stringToDecimal(balance.Frozen),
		Total:     c.stringToDecimal(total),
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
			availFloat, _ := c.stringToDecimal(bal.Available).Float64()
			unavailFloat, _ := c.stringToDecimal(bal.UnAvailable).Float64()
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
			Available: c.stringToDecimal(bal.Available),
			Frozen:    c.stringToDecimal(bal.Frozen),
			Total:     c.stringToDecimal(total),
			Extra: map[string]interface{}{
				"name":                    bal.Name,
				"available_usd_valuation": bal.AvailableUsdValuation,
				"unAvailable":             bal.UnAvailable,
			},
		})
	}

	return &commontypes.AccountBalance{
		Balances:    commonBalances,
		TotalEquity: c.stringToDecimal(strconv.FormatFloat(totalEquity, 'f', -1, 64)),
	}
}

// ConvertTicker converts BitMart ticker to common ticker type
func (c *Converter) ConvertTicker(ticker *marketmodels.Ticker) *commontypes.Ticker {
	if ticker == nil {
		return nil
	}

	return &commontypes.Ticker{
		Symbol:    ticker.Symbol,
		LastPrice: c.stringToDecimal(ticker.LastPrice),
		High24h:   c.stringToDecimal(ticker.HighPrice),
		Low24h:    c.stringToDecimal(ticker.LowPrice),
		Volume24h: c.stringToDecimal(ticker.BaseVolume),
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
				Price:    c.stringToDecimal(bid[0]),
				Quantity: c.stringToDecimal(bid[1]),
			}
		}
	}

	// Convert asks
	asks := make([]commontypes.OrderBookLevel, len(data.Asks))
	for i, ask := range data.Asks {
		if len(ask) >= 2 {
			asks[i] = commontypes.OrderBookLevel{
				Price:    c.stringToDecimal(ask[0]),
				Quantity: c.stringToDecimal(ask[1]),
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
		CtVal:             c.stringToDecimal(symbol.ContractSize),
		InstrumentType:    instType,
		Status:            symbol.Status,
		MinOrderSize:      c.stringToDecimal(symbol.MinVolume),
		MaxOrderSize:      c.stringToDecimal(symbol.MaxVolume),
		PricePrecision:    c.stringToDecimal(symbol.PricePrecision),
		QuantityPrecision: c.stringToDecimal(symbol.VolPrecision),
		LastPrice:         c.stringToDecimal(symbol.LastPrice),
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
			LastPrice: c.stringToDecimal(last),
			BidPrice:  c.stringToDecimal(bidPx),
			AskPrice:  c.stringToDecimal(askPx),
			High24h:   c.stringToDecimal(high24h),
			Low24h:    c.stringToDecimal(low24h),
			Volume24h: c.stringToDecimal(volume24h),
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
		Open:        c.stringToDecimal(open),
		High:        c.stringToDecimal(high),
		Low:         c.stringToDecimal(low),
		Close:       c.stringToDecimal(close_),
		Volume:      c.stringToDecimal(volume),
		QuoteVolume: c.stringToDecimal(quoteVolume),
		Timestamp:   commontypes.Timestamp(time.Unix(int64(timestamp), 0)),
		Confirmed:   true, // Historical candles are always confirmed
		Extra:       make(map[string]interface{}),
	}
}

// ConvertContractKlineToCandle converts BitMart contract kline data to common Candle type
// BitMart contract format: {timestamp, open_price, close_price, high_price, low_price, volume}
func (c *Converter) ConvertContractKlineToCandle(klineData *contractresponses.ContractKlineData, symbol, interval string) *commontypes.Candle {
	if klineData == nil {
		return nil
	}

	return &commontypes.Candle{
		Symbol:      symbol,
		Interval:    interval,
		Open:        c.stringToDecimal(klineData.OpenPrice),
		High:        c.stringToDecimal(klineData.HighPrice),
		Low:         c.stringToDecimal(klineData.LowPrice),
		Close:       c.stringToDecimal(klineData.ClosePrice),
		Volume:      c.stringToDecimal(klineData.Volume),
		QuoteVolume: commontypes.ZeroDecimal, // Contract kline doesn't provide quote volume
		Timestamp:   commontypes.Timestamp(time.Unix(klineData.Timestamp, 0)),
		Confirmed:   true, // Historical candles are always confirmed
		Extra: map[string]interface{}{
			"type": "contract", // Mark as contract kline data
		},
	}
}

// ConvertPositionV2ToPosition converts BitMart PositionV2 to common Position type
func (c *Converter) ConvertPositionV2ToPosition(position *contractresponses.PositionV2) *commontypes.Position {
	if position == nil {
		return nil
	}

	// Skip positions with zero quantity (no actual position)
	currentAmount := c.stringToDecimal(position.CurrentAmount)
	if currentAmount.IsZero() {
		return nil
	}

	// Convert open_type to MarginMode
	var marginMode commontypes.MarginMode
	switch position.OpenType {
	case "cross":
		marginMode = commontypes.MarginModeCross
	case "isolated":
		marginMode = commontypes.MarginModeIsolated
	default:
		marginMode = commontypes.MarginModeIsolated
	}

	// Convert position_side to PositionSide
	var posSide commontypes.PositionSide
	switch position.PositionSide {
	case "long":
		posSide = commontypes.PositionSideLong
	case "short":
		posSide = commontypes.PositionSideShort
	case "both":
		// In one-way mode, position_side is "both"
		// Need to determine direction from position_amount sign
		positionAmount := c.stringToDecimal(position.PositionAmount)
		if positionAmount.IsPositive() {
			posSide = commontypes.PositionSideLong
		} else if positionAmount.IsNegative() {
			posSide = commontypes.PositionSideShort
		} else {
			posSide = commontypes.PositionSideNet
		}
	default:
		posSide = commontypes.PositionSideNet
	}

	// Parse leverage
	leverageInt, _ := strconv.Atoi(position.Leverage)

	// Get absolute value for quantity
	absQuantity, _ := currentAmount.Abs()

	return &commontypes.Position{
		Symbol:           position.Symbol,
		PosSide:          posSide,
		Quantity:         absQuantity,
		AvgPrice:         c.stringToDecimal(position.EntryPrice),
		MarkPrice:        c.stringToDecimal(position.MarkPrice),
		LiquidationPrice: c.stringToDecimal(position.LiquidationPrice),
		Leverage:         leverageInt,
		MarginMode:       marginMode,
		UnrealizedPnL:    c.stringToDecimal(position.UnrealizedPnl),
		RealizedPnL:      c.stringToDecimal(position.RealizedValue),
		CreatedAt:        commontypes.Timestamp(time.UnixMilli(position.OpenTimestamp)),
		UpdatedAt:        commontypes.Timestamp(time.UnixMilli(position.Timestamp)),
		Extra: map[string]interface{}{
			"account":              position.Account,
			"position_value":       position.PositionValue,
			"position_cross":       position.PositionCross,
			"initial_margin":       position.InitialMargin,
			"maintenance_margin":   position.MaintenanceMargin,
			"current_fee":          position.CurrentFee,
			"close_vol":            position.CloseVol,
			"close_avg_price":      position.CloseAvgPrice,
			"open_avg_price":       position.OpenAvgPrice,
			"max_notional_value":   position.MaxNotionalValue,
			"position_amount":      position.PositionAmount,
			"current_amount":       position.CurrentAmount,
			"mark_value":           position.MarkValue,
			"current_value":        position.CurrentValue,
		},
	}
}

// ConvertPositionV2ToLeverage converts BitMart PositionV2 to common Leverage type
func (c *Converter) ConvertPositionV2ToLeverage(position *contractresponses.PositionV2) *commontypes.Leverage {
	if position == nil {
		return nil
	}

	// Convert leverage string to int
	leverage, err := strconv.Atoi(position.Leverage)
	if err != nil {
		// If conversion fails, default to 1
		leverage = 1
	}

	// Convert open_type to MarginMode
	var marginMode commontypes.MarginMode
	switch position.OpenType {
	case "cross":
		marginMode = commontypes.MarginModeCross
	case "isolated":
		marginMode = commontypes.MarginModeIsolated
	default:
		marginMode = commontypes.MarginModeIsolated
	}

	// Convert position_side to PositionSide
	var posSide commontypes.PositionSide
	switch position.PositionSide {
	case "long":
		posSide = commontypes.PositionSideLong
	case "short":
		posSide = commontypes.PositionSideShort
	case "both":
		// In one-way mode, position_side is "both"
		posSide = commontypes.PositionSideNet
	default:
		posSide = commontypes.PositionSideNet
	}

	return &commontypes.Leverage{
		Symbol:     position.Symbol,
		Leverage:   leverage,
		MarginMode: marginMode,
		PosSide:    posSide,
		Extra: map[string]interface{}{
			"account":            position.Account,
			"mark_price":         position.MarkPrice,
			"position_amount":    position.PositionAmount,
			"current_amount":     position.CurrentAmount,
			"entry_price":        position.EntryPrice,
			"liquidation_price":  position.LiquidationPrice,
			"unrealized_pnl":     position.UnrealizedPnl,
			"max_notional_value": position.MaxNotionalValue,
			"initial_margin":     position.InitialMargin,
			"maintenance_margin": position.MaintenanceMargin,
			"timestamp":          position.Timestamp,
		},
	}
}

// ConvertSubmitLeverageResponse converts BitMart SubmitLeverageResponse to common Leverage type
func (c *Converter) ConvertSubmitLeverageResponse(resp *contractresponses.SubmitLeverageResponse) *commontypes.Leverage {
	if resp == nil {
		return nil
	}

	// Convert leverage string to int
	leverage, err := strconv.Atoi(resp.Data.Leverage)
	if err != nil {
		// If conversion fails, default to 1
		leverage = 1
	}

	// Convert open_type to MarginMode
	var marginMode commontypes.MarginMode
	switch resp.Data.OpenType {
	case "cross":
		marginMode = commontypes.MarginModeCross
	case "isolated":
		marginMode = commontypes.MarginModeIsolated
	default:
		marginMode = commontypes.MarginModeIsolated
	}

	return &commontypes.Leverage{
		Symbol:     resp.Data.Symbol,
		Leverage:   leverage,
		MarginMode: marginMode,
		PosSide:    commontypes.PositionSideNet, // BitMart doesn't specify position side in leverage response
		Extra: map[string]interface{}{
			"max_value": resp.Data.MaxValue,
		},
	}
}

// ConvertToContractSide converts common side and position side to BitMart contract side
// BitMart contract side values:
//   - For dual position mode (hedge mode):
//     1 = Open long
//     2 = Close short
//     3 = Close long
//     4 = Open short
//   - For single position mode (one-way mode):
//     1 = Buy
//     4 = Sell
func (c *Converter) ConvertToContractSide(side commontypes.OrderSide, posSide commontypes.PositionSide) int {
	// For hedge mode (long/short specified)
	if posSide == commontypes.PositionSideLong {
		if side == "buy" {
			return 1 // Open long
		}
		return 3 // Close long
	}
	if posSide == commontypes.PositionSideShort {
		if side == "buy" {
			return 2 // Close short
		}
		return 4 // Open short
	}

	// For one-way mode or unspecified position side
	if side == "buy" {
		return 1 // Buy
	}
	return 4 // Sell
}

// ConvertContractOrder converts BitMart contract order response to common Order type
func (c *Converter) ConvertContractOrder(resp *contractresponses.SubmitOrderResponse) *commontypes.Order {
	if resp == nil {
		return nil
	}

	orderID := fmt.Sprintf("%d", resp.Data.OrderID)

	return &commontypes.Order{
		ID:     orderID,
		Price:  c.stringToDecimal(resp.Data.Price),
		Status: "new", // New order, status pending confirmation
		Extra: map[string]interface{}{
			"order_id": resp.Data.OrderID,
		},
	}
}
