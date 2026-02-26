package bingx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/djpken/go-exc/exchanges/bingx/rest"
	commontypes "github.com/djpken/go-exc/types"
)

// Converter converts between BingX-specific types and common types
type Converter struct{}

func NewConverter() *Converter { return &Converter{} }

func (c *Converter) str(s string) commontypes.Decimal {
	if s == "" {
		return commontypes.ZeroDecimal
	}
	d, err := commontypes.NewDecimal(s)
	if err != nil {
		return commontypes.ZeroDecimal
	}
	return d
}

// ConvertTicker converts a BingX TickerData to the common Ticker type
func (c *Converter) ConvertTicker(t *rest.TickerData) *commontypes.Ticker {
	if t == nil {
		return nil
	}
	return &commontypes.Ticker{
		Symbol:    t.Symbol,
		LastPrice: c.str(t.LastPrice),
		BidPrice:  c.str(t.BidPrice),
		AskPrice:  c.str(t.AskPrice),
		High24h:   c.str(t.HighPrice),
		Low24h:    c.str(t.LowPrice),
		Volume24h: c.str(t.Volume),
		Timestamp: commontypes.Timestamp(time.UnixMilli(t.CloseTime)),
		Extra: map[string]interface{}{
			"quoteVolume":        t.QuoteVolume,
			"openPrice":          t.OpenPrice,
			"priceChange":        t.PriceChange,
			"priceChangePercent": t.PriceChangePercent,
			"count":              t.Count,
		},
	}
}

// ConvertOrderBook converts BingX order book data to the common OrderBook type
func (c *Converter) ConvertOrderBook(ob *rest.OrderBookData, symbol string) *commontypes.OrderBook {
	if ob == nil {
		return nil
	}
	bids := make([]commontypes.OrderBookLevel, 0, len(ob.Bids))
	for _, b := range ob.Bids {
		if len(b) >= 2 {
			bids = append(bids, commontypes.OrderBookLevel{
				Price:    c.str(b[0]),
				Quantity: c.str(b[1]),
			})
		}
	}
	asks := make([]commontypes.OrderBookLevel, 0, len(ob.Asks))
	for _, a := range ob.Asks {
		if len(a) >= 2 {
			asks = append(asks, commontypes.OrderBookLevel{
				Price:    c.str(a[0]),
				Quantity: c.str(a[1]),
			})
		}
	}
	return &commontypes.OrderBook{
		Symbol:    symbol,
		Bids:      bids,
		Asks:      asks,
		Timestamp: commontypes.Timestamp(time.UnixMilli(ob.T)),
	}
}

// ConvertKline converts a BingX KlineEntry to the common Candle type
func (c *Converter) ConvertKline(k *rest.KlineEntry, symbol, interval string) *commontypes.Candle {
	if k == nil {
		return nil
	}
	return &commontypes.Candle{
		Symbol:    symbol,
		Interval:  interval,
		Open:      c.str(k.Open),
		High:      c.str(k.High),
		Low:       c.str(k.Low),
		Close:     c.str(k.Close),
		Volume:    c.str(k.Volume),
		Timestamp: commontypes.Timestamp(time.UnixMilli(k.Time)),
		Confirmed: true,
		Extra:     map[string]interface{}{},
	}
}

// ConvertBalance converts BingX BalanceAsset to the common AccountBalance type
func (c *Converter) ConvertBalance(b *rest.BalanceAsset) *commontypes.AccountBalance {
	if b == nil {
		return nil
	}
	available := c.str(b.AvailableMargin)
	frozen := c.str(b.UsedMargin)
	total := c.str(b.Balance)

	return &commontypes.AccountBalance{
		Balances: []*commontypes.Balance{
			{
				Currency:  b.Asset,
				Available: available,
				Frozen:    frozen,
				Total:     total,
				Extra: map[string]interface{}{
					"equity":           b.Equity,
					"unrealizedProfit": b.UnrealizedProfit,
					"realisedProfit":   b.RealisedProfit,
					"usedMargin":       b.UsedMargin,
					"freezedMargin":    b.FreezedMargin,
				},
			},
		},
		TotalEquity: c.str(b.Equity),
	}
}

// ConvertPosition converts BingX PositionData to the common Position type
func (c *Converter) ConvertPosition(p *rest.PositionData) *commontypes.Position {
	if p == nil {
		return nil
	}

	var posSide commontypes.PositionSide
	switch p.PositionSide {
	case "LONG":
		posSide = commontypes.PositionSideLong
	case "SHORT":
		posSide = commontypes.PositionSideShort
	default:
		posSide = commontypes.PositionSideNet
	}

	var marginMode commontypes.MarginMode
	if p.Isolated {
		marginMode = commontypes.MarginModeIsolated
	} else {
		marginMode = commontypes.MarginModeCross
	}

	return &commontypes.Position{
		Symbol:        p.Symbol,
		PosSide:       posSide,
		Quantity:      c.str(p.PositionAmt),
		AvgPrice:      c.str(p.AvgPrice),
		Leverage:      p.Leverage,
		MarginMode:    marginMode,
		UnrealizedPnL: c.str(p.UnrealizedProfit),
		RealizedPnL:   c.str(p.RealisedProfit),
		Extra: map[string]interface{}{
			"positionId":    p.PositionID,
			"availableAmt":  p.AvailableAmt,
			"initialMargin": p.InitialMargin,
		},
	}
}

// ConvertLeverage converts BingX LeverageData to the common Leverage type
func (c *Converter) ConvertLeverage(symbol string, d *rest.LeverageData) []*commontypes.Leverage {
	if d == nil {
		return nil
	}
	return []*commontypes.Leverage{
		{
			Symbol:   symbol,
			Leverage: int(d.LongLeverage),
			PosSide:  commontypes.PositionSideLong,
			Extra: map[string]interface{}{
				"maxLongLeverage":  d.MaxLongLeverage,
				"maxShortLeverage": d.MaxShortLeverage,
			},
		},
		{
			Symbol:   symbol,
			Leverage: int(d.ShortLeverage),
			PosSide:  commontypes.PositionSideShort,
			Extra: map[string]interface{}{
				"maxLongLeverage":  d.MaxLongLeverage,
				"maxShortLeverage": d.MaxShortLeverage,
			},
		},
	}
}

// ConvertOrder converts BingX OrderData to the common Order type
func (c *Converter) ConvertOrder(o *rest.OrderData) *commontypes.Order {
	if o == nil {
		return nil
	}
	return &commontypes.Order{
		ID:             strconv.FormatInt(o.OrderID, 10),
		Symbol:         o.Symbol,
		Side:           o.Side,
		Type:           o.Type,
		Price:          c.str(o.Price),
		Quantity:       c.str(o.OrigQty),
		FilledQuantity: c.str(o.ExecutedQty),
		Status:         c.ConvertOrderStatus(o.Status),
		ClientOrderID:  o.ClientOrderID,
		CreatedAt:      commontypes.Timestamp(time.UnixMilli(o.Time)),
		UpdatedAt:      commontypes.Timestamp(time.UnixMilli(o.UpdateTime)),
		Extra: map[string]interface{}{
			"positionSide": o.PositionSide,
			"avgPrice":     o.AvgPrice,
			"cumQuote":     o.CumQuote,
			"stopPrice":    o.StopPrice,
			"profit":       o.Profit,
			"commission":   o.Commission,
			"workingType":  o.WorkingType,
		},
	}
}

// ConvertOrderSide converts a BingX order side string to common form
func (c *Converter) ConvertOrderSide(s string) string {
	switch s {
	case "BUY":
		return string(commontypes.OrderSideBuy)
	case "SELL":
		return string(commontypes.OrderSideSell)
	default:
		return s
	}
}

// ConvertOrderType converts a BingX order type string to common form
func (c *Converter) ConvertOrderType(s string) string {
	switch s {
	case "MARKET":
		return string(commontypes.OrderTypeMarket)
	case "LIMIT":
		return string(commontypes.OrderTypeLimit)
	default:
		return s
	}
}

// ConvertOrderStatus converts BingX order status to common status
func (c *Converter) ConvertOrderStatus(s string) string {
	switch s {
	case "NEW", "PENDING":
		return string(commontypes.OrderStatusOpen)
	case "PARTIALLY_FILLED":
		return string(commontypes.OrderStatusPartiallyFilled)
	case "FILLED":
		return string(commontypes.OrderStatusFilled)
	case "CANCELED", "CANCELLED":
		return string(commontypes.OrderStatusCanceled)
	default:
		return s
	}
}

// ConvertInstrument converts BingX ContractInfo to the common Instrument type
func (c *Converter) ConvertInstrument(ci *rest.ContractInfo) *commontypes.Instrument {
	if ci == nil {
		return nil
	}
	return &commontypes.Instrument{
		Symbol:            ci.Symbol,
		BaseCurrency:      ci.Asset,
		QuoteCurrency:     ci.Currency,
		InstrumentType:    commontypes.InstrumentSwap,
		MaxLever:          ci.MaxLongLeverage,
		CtVal:             c.str(ci.Size),
		PricePrecision:    commontypes.NewDecimalFromInt(int64(ci.PricePrecision)),
		QuantityPrecision: commontypes.NewDecimalFromInt(int64(ci.QuantityPrecision)),
	}
}

// ConvertIntervalToWS maps common interval strings to BingX WebSocket kline interval format
func (c *Converter) ConvertIntervalToWS(interval string) (string, error) {
	m := map[string]string{
		"1m":  "1min",
		"3m":  "3min",
		"5m":  "5min",
		"15m": "15min",
		"30m": "30min",
		"1H":  "1hour",
		"2H":  "2hour",
		"4H":  "4hour",
		"6H":  "6hour",
		"8H":  "8hour",
		"12H": "12hour",
		"1D":  "1day",
		"3D":  "3day",
		"1W":  "1week",
		"1M":  "1month",
	}
	v, ok := m[interval]
	if !ok {
		return "", fmt.Errorf("bingx: unsupported interval: %s", interval)
	}
	return v, nil
}

// ConvertIntervalToREST maps common interval strings to BingX REST API kline interval format
func (c *Converter) ConvertIntervalToREST(interval string) (string, error) {
	m := map[string]string{
		"1m":  "1m",
		"3m":  "3m",
		"5m":  "5m",
		"15m": "15m",
		"30m": "30m",
		"1H":  "1h",
		"2H":  "2h",
		"4H":  "4h",
		"6H":  "6h",
		"8H":  "8h",
		"12H": "12h",
		"1D":  "1d",
		"3D":  "3d",
		"1W":  "1w",
		"1M":  "1M",
	}
	v, ok := m[interval]
	if !ok {
		return "", fmt.Errorf("bingx: unsupported interval: %s", interval)
	}
	return v, nil
}
