package okex

import (
	"strconv"
	"time"

	okexconstants "github.com/djpken/go-exc/exchanges/okex/constants"
	"github.com/djpken/go-exc/exchanges/okex/models/account"
	"github.com/djpken/go-exc/exchanges/okex/models/market"
	"github.com/djpken/go-exc/exchanges/okex/models/publicdata"
	"github.com/djpken/go-exc/exchanges/okex/models/trade"
	commontypes "github.com/djpken/go-exc/types"
)

// Converter handles conversion between OKEx types and common types
type Converter struct {
	constantsConverter *ConstantsConverter
}

// NewConverter creates a new Converter
func NewConverter() *Converter {
	return &Converter{
		constantsConverter: NewConstantsConverter(),
	}
}

// ConvertOrder converts OKEx order to common Order type
func (c *Converter) ConvertOrder(okexOrder *trade.Order) *commontypes.Order {
	if okexOrder == nil {
		return nil
	}

	price := float64(okexOrder.AvgPx)
	quantity := float64(okexOrder.Sz)
	filledQty := float64(okexOrder.AccFillSz)
	fee := float64(okexOrder.Fee)

	return &commontypes.Order{
		ID:                okexOrder.OrdID,
		Symbol:            okexOrder.InstID,
		Side:              string(okexOrder.Side),
		Type:              string(okexOrder.OrdType),
		Status:            string(okexOrder.State),
		Price:             commontypes.Decimal(strconv.FormatFloat(price, 'f', -1, 64)),
		Quantity:          commontypes.Decimal(strconv.FormatFloat(quantity, 'f', -1, 64)),
		FilledQuantity:    commontypes.Decimal(strconv.FormatFloat(filledQty, 'f', -1, 64)),
		RemainingQuantity: commontypes.Decimal(strconv.FormatFloat(quantity-filledQty, 'f', -1, 64)),
		Fee:               commontypes.Decimal(strconv.FormatFloat(fee, 'f', -1, 64)),
		FeeCurrency:       okexOrder.FeeCcy,
		CreatedAt:         commontypes.Timestamp(okexOrder.CTime),
		UpdatedAt:         commontypes.Timestamp(time.Time(okexOrder.UTime)),
		Extra: map[string]interface{}{
			"clOrdID":  okexOrder.ClOrdID,
			"tag":      okexOrder.Tag,
			"category": okexOrder.Category,
			"tdMode":   okexOrder.TdMode,
		},
	}
}

// ConvertBalance converts OKEx balance to common AccountBalance type
func (c *Converter) ConvertBalance(okexBalance *account.Balance) *commontypes.AccountBalance {
	if okexBalance == nil {
		return nil
	}

	balances := make([]*commontypes.Balance, 0, len(okexBalance.Details))
	for _, detail := range okexBalance.Details {
		balances = append(balances, &commontypes.Balance{
			Currency:  detail.Ccy,
			Available: commontypes.Decimal(strconv.FormatFloat(float64(detail.AvailBal), 'f', -1, 64)),
			Frozen:    commontypes.Decimal(strconv.FormatFloat(float64(detail.FrozenBal), 'f', -1, 64)),
			Total:     commontypes.Decimal(strconv.FormatFloat(float64(detail.Eq), 'f', -1, 64)),
			Extra: map[string]interface{}{
				"cashBal":   detail.CashBal,
				"upl":       detail.Upl,
				"liab":      detail.Liab,
				"interest":  detail.Interest,
				"ordFrozen": detail.OrdFrozen,
				"isoEq":     detail.IsoEq,
				"availEq":   detail.AvailEq,
				"disEq":     detail.DisEq,
			},
		})
	}

	return &commontypes.AccountBalance{
		Balances:    balances,
		TotalEquity: commontypes.Decimal(strconv.FormatFloat(float64(okexBalance.TotalEq), 'f', -1, 64)),
		UpdatedAt:   commontypes.Timestamp(okexBalance.UTime),
		Extra: map[string]interface{}{
			"isoEq":    okexBalance.IsoEq,
			"adjEq":    okexBalance.AdjEq,
			"ordFroz":  okexBalance.OrdFroz,
			"imr":      okexBalance.Imr,
			"mmr":      okexBalance.Mmr,
			"mgnRatio": okexBalance.MgnRatio,
		},
	}
}

// ConvertPosition converts OKEx position to common Position type
func (c *Converter) ConvertPosition(okexPos *account.Position) *commontypes.Position {
	if okexPos == nil {
		return nil
	}

	return &commontypes.Position{
		Symbol:           okexPos.InstID,
		Side:             string(okexPos.PosSide),
		Quantity:         commontypes.Decimal(strconv.FormatFloat(float64(okexPos.Pos), 'f', -1, 64)),
		AvgPrice:         commontypes.Decimal(strconv.FormatFloat(float64(okexPos.AvgPx), 'f', -1, 64)),
		MarkPrice:        commontypes.Decimal(strconv.FormatFloat(float64(okexPos.Last), 'f', -1, 64)), // Using Last as mark price
		LiquidationPrice: commontypes.Decimal(strconv.FormatFloat(float64(okexPos.LiqPx), 'f', -1, 64)),
		UnrealizedPnL:    commontypes.Decimal(strconv.FormatFloat(float64(okexPos.Upl), 'f', -1, 64)),
		RealizedPnL:      commontypes.Decimal(strconv.FormatFloat(float64(okexPos.RealizedPnl), 'f', -1, 64)),
		Leverage:         commontypes.Decimal(strconv.FormatFloat(float64(okexPos.Lever), 'f', -1, 64)),
		MarginMode:       string(okexPos.MgnMode),
		CreatedAt:        commontypes.Timestamp(time.Time(okexPos.CTime)),
		UpdatedAt:        commontypes.Timestamp(time.Time(okexPos.UTime)),
		Extra: map[string]interface{}{
			"posID":       okexPos.PosID,
			"tradeID":     okexPos.TradeID,
			"instType":    okexPos.InstType,
			"margin":      okexPos.Margin,
			"imr":         okexPos.Imr,
			"mmr":         okexPos.Mmr,
			"liab":        okexPos.Liab,
			"interest":    okexPos.Interest,
			"notionalUsd": okexPos.NotionalUsd,
			"adl":         okexPos.ADL,
		},
	}
}

// ConvertOrderSide converts common order side string to OKEx order side
// Deprecated: Use ConstantsConverter.ToOKExOrderSide with typed constants instead
func (c *Converter) ConvertOrderSide(side string) okexconstants.OrderSide {
	// Convert string to typed constant
	var typedSide commontypes.OrderSide
	switch side {
	case "buy":
		typedSide = commontypes.OrderSideBuy
	case "sell":
		typedSide = commontypes.OrderSideSell
	default:
		typedSide = commontypes.OrderSideBuy
	}
	return c.constantsConverter.ToOKExOrderSide(typedSide)
}

// ConvertOrderType converts common order type string to OKEx order type
// Deprecated: Use ConstantsConverter.ToOKExOrderType with typed constants instead
func (c *Converter) ConvertOrderType(orderType string) okexconstants.OrderType {
	// Convert string to typed constant
	var typedOrderType commontypes.OrderType
	switch orderType {
	case "market":
		typedOrderType = commontypes.OrderTypeMarket
	case "limit":
		typedOrderType = commontypes.OrderTypeLimit
	case "post_only":
		typedOrderType = commontypes.OrderTypePostOnly
	case "fok":
		typedOrderType = commontypes.OrderTypeFOK
	case "ioc":
		typedOrderType = commontypes.OrderTypeIOC
	default:
		typedOrderType = commontypes.OrderTypeLimit
	}
	return c.constantsConverter.ToOKExOrderType(typedOrderType)
}

// ConvertInstrumentType converts common instrument type to OKEx instrument type
// Use ConstantsConverter.ToOKExInstrumentType for the same functionality
func (c *Converter) ConvertInstrumentType(instrumentType commontypes.InstrumentType) okexconstants.InstrumentType {
	return c.constantsConverter.ToOKExInstrumentType(instrumentType)
}

// ConvertAccountConfig converts OKEx account config to common AccountConfig type
func (c *Converter) ConvertAccountConfig(okexConfig *account.Config) *commontypes.AccountConfig {
	if okexConfig == nil {
		return nil
	}

	// Convert position mode
	posMode := "net_mode"
	if okexConfig.PosMode == ("long_short_mode") {
		posMode = "long_short_mode"
	}

	return &commontypes.AccountConfig{
		UID:          okexConfig.UID,
		Level:        okexConfig.Level,
		AutoLoan:     okexConfig.AutoLoan,
		PositionMode: posMode,
		Extra: map[string]interface{}{
			"levelTmp":   okexConfig.LevelTmp,
			"acctLv":     okexConfig.AcctLv,
			"greeksType": okexConfig.GreeksType,
		},
	}
}

// ConvertBalanceAndPosition converts OKEx BalanceAndPosition event to common type
func (c *Converter) ConvertBalanceAndPosition(okexBNP *account.BalanceAndPosition) *commontypes.BalanceAndPositionUpdate {
	if okexBNP == nil {
		return nil
	}

	// Convert balances
	balances := make([]*commontypes.Balance, 0, len(okexBNP.BalData))
	for _, bal := range okexBNP.BalData {
		cashBal := float64(bal.CashBal)
		balances = append(balances, &commontypes.Balance{
			Currency:  bal.Ccy,
			Available: commontypes.Decimal(strconv.FormatFloat(cashBal, 'f', -1, 64)),
			Frozen:    commontypes.Decimal("0"), // BalData doesn't have frozen field
			Total:     commontypes.Decimal(strconv.FormatFloat(cashBal, 'f', -1, 64)),
			Extra: map[string]interface{}{
				"uTime": bal.UTime,
			},
		})
	}

	// Convert positions
	positions := make([]*commontypes.Position, 0, len(okexBNP.PosData))
	for _, pos := range okexBNP.PosData {
		positions = append(positions, &commontypes.Position{
			Symbol:        pos.InstId,
			Side:          string(pos.PosSide),
			Quantity:      commontypes.Decimal(strconv.FormatFloat(float64(pos.Pos), 'f', -1, 64)),
			AvgPrice:      commontypes.Decimal(strconv.FormatFloat(float64(pos.AvgPx), 'f', -1, 64)),
			MarkPrice:     commontypes.Decimal("0"), // PosData doesn't have mark price
			UnrealizedPnL: commontypes.Decimal("0"), // PosData doesn't have unrealized PnL
			RealizedPnL:   commontypes.Decimal(strconv.FormatFloat(float64(pos.SettledPnl), 'f', -1, 64)),
			Leverage:      commontypes.Decimal("0"), // PosData doesn't have leverage
			MarginMode:    string(pos.MgnMode),
			UpdatedAt:     commontypes.Timestamp(time.Time(pos.UTime)),
			Extra: map[string]interface{}{
				"posId":          pos.PosId,
				"tradeId":        pos.TradeId,
				"instType":       pos.InstType,
				"ccy":            pos.Ccy,
				"posCcy":         pos.PosCcy,
				"nonSettleAvgPx": pos.NonSettleAvgPx,
			},
		})
	}

	return &commontypes.BalanceAndPositionUpdate{
		Balances:  balances,
		Positions: positions,
		EventType: string(okexBNP.EventType),
		UpdatedAt: commontypes.Timestamp(time.Time(okexBNP.PTime)),
		Extra: map[string]interface{}{
			"trades": okexBNP.Trades,
		},
	}
}

// ConvertAccountEvent converts OKEx Account event to common AccountUpdate
func (c *Converter) ConvertAccountEvent(okexBalances []*account.Balance, eventType string) *commontypes.AccountUpdate {
	if len(okexBalances) == 0 {
		return nil
	}

	balances := make([]*commontypes.Balance, 0)
	var updateTime time.Time
	var totalEquity float64

	for _, okexBal := range okexBalances {
		for _, detail := range okexBal.Details {
			balances = append(balances, &commontypes.Balance{
				Currency:  detail.Ccy,
				Available: commontypes.Decimal(strconv.FormatFloat(float64(detail.AvailBal), 'f', -1, 64)),
				Frozen:    commontypes.Decimal(strconv.FormatFloat(float64(detail.FrozenBal), 'f', -1, 64)),
				Total:     commontypes.Decimal(strconv.FormatFloat(float64(detail.Eq), 'f', -1, 64)),
				Extra: map[string]interface{}{
					"cashBal":   detail.CashBal,
					"upl":       detail.Upl,
					"ordFrozen": detail.OrdFrozen,
				},
			})
		}
		updateTime = time.Time(okexBal.UTime)
		totalEquity = float64(okexBal.TotalEq)
	}

	return &commontypes.AccountUpdate{
		Balances:    balances,
		EventType:   eventType,
		UpdatedAt:   commontypes.Timestamp(updateTime),
		TotalEquity: commontypes.Decimal(strconv.FormatFloat(totalEquity, 'f', -1, 64)),
		Extra:       map[string]interface{}{},
	}
}

// ConvertPositionEvent converts OKEx Position event to common PositionUpdate
func (c *Converter) ConvertPositionEvent(okexPositions []*account.Position, eventType string) *commontypes.PositionUpdate {
	if len(okexPositions) == 0 {
		return nil
	}

	positions := make([]*commontypes.Position, 0, len(okexPositions))
	var updateTime time.Time

	for _, pos := range okexPositions {
		positions = append(positions, c.ConvertPosition(pos))
		updateTime = time.Time(pos.UTime)
	}

	return &commontypes.PositionUpdate{
		Positions: positions,
		EventType: eventType,
		UpdatedAt: commontypes.Timestamp(updateTime),
		Extra:     map[string]interface{}{},
	}
}

// ConvertOrderEvent converts OKEx Order event to common OrderUpdate
func (c *Converter) ConvertOrderEvent(okexOrders []*trade.Order) *commontypes.OrderUpdate {
	if len(okexOrders) == 0 {
		return nil
	}

	orders := make([]*commontypes.Order, 0, len(okexOrders))
	var updateTime time.Time

	for _, order := range okexOrders {
		orders = append(orders, c.ConvertOrder(order))
		updateTime = time.Time(order.UTime)
	}

	return &commontypes.OrderUpdate{
		Orders:    orders,
		UpdatedAt: commontypes.Timestamp(updateTime),
		Extra:     map[string]interface{}{},
	}
}

// ConvertInstrument converts OKEx instrument to common Instrument type
func (c *Converter) ConvertInstrument(okexInst *publicdata.Instrument, ticker market.Ticker) *commontypes.Instrument {
	if okexInst == nil {
		return nil
	}

	// Convert instrument type
	instType := c.constantsConverter.FromOKExInstrumentType(okexInst.InstType)
	if instType == commontypes.InstrumentSwap {
		return &commontypes.Instrument{
			Symbol:            okexInst.InstID,
			BaseCurrency:      okexInst.CtValCcy,
			QuoteCurrency:     okexInst.SettleCcy,
			CtVal:             commontypes.Decimal(strconv.FormatFloat(float64(okexInst.CtVal), 'f', -1, 64)),
			InstrumentType:    commontypes.InstrumentSwap,
			Status:            string(okexInst.State),
			MinOrderSize:      commontypes.Decimal(strconv.FormatFloat(float64(okexInst.MinSz), 'f', -1, 64)),
			MaxOrderSize:      "0",
			PricePrecision:    commontypes.Decimal(strconv.FormatFloat(float64(okexInst.TickSz), 'f', -1, 64)),
			QuantityPrecision: commontypes.Decimal(strconv.FormatFloat(float64(okexInst.LotSz), 'f', -1, 64)),
			MaxLever:          int(okexInst.Lever),
			LastPrice:         commontypes.Decimal(strconv.FormatFloat(float64(ticker.Last), 'f', -1, 64)),
		}
	}

	return &commontypes.Instrument{
		Symbol:            okexInst.InstID,
		BaseCurrency:      okexInst.BaseCcy,
		QuoteCurrency:     okexInst.QuoteCcy,
		CtVal:             commontypes.Decimal(strconv.FormatFloat(float64(okexInst.CtVal), 'f', -1, 64)),
		InstrumentType:    instType,
		Status:            string(okexInst.State),
		MinOrderSize:      commontypes.Decimal(strconv.FormatFloat(float64(okexInst.MinSz), 'f', -1, 64)),
		MaxOrderSize:      "0",
		PricePrecision:    commontypes.Decimal(strconv.FormatFloat(float64(okexInst.TickSz), 'f', -1, 64)),
		QuantityPrecision: commontypes.Decimal(strconv.FormatFloat(float64(okexInst.LotSz), 'f', -1, 64)),
		MaxLever:          int(okexInst.Lever),
		LastPrice:         commontypes.Decimal(strconv.FormatFloat(float64(ticker.Last), 'f', -1, 64)),
	}
}

// ConvertTicker converts OKEx ticker to common Ticker type
func (c *Converter) ConvertTicker(okexTicker *market.Ticker) *commontypes.Ticker {
	if okexTicker == nil {
		return nil
	}

	return &commontypes.Ticker{
		Symbol:    okexTicker.InstID,
		LastPrice: commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.Last), 'f', -1, 64)),
		BidPrice:  commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.BidPx), 'f', -1, 64)),
		AskPrice:  commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.AskPx), 'f', -1, 64)),
		High24h:   commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.High24h), 'f', -1, 64)),
		Low24h:    commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.Low24h), 'f', -1, 64)),
		Volume24h: commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.Vol24h), 'f', -1, 64)),
		Timestamp: commontypes.Timestamp(okexTicker.TS),
		Extra: map[string]interface{}{
			"open24h":   okexTicker.Open24h,
			"volCcy24h": okexTicker.VolCcy24h,
			"instType":  okexTicker.InstType,
		},
	}
}

// ConvertTickerToUpdate converts an OKEx ticker to common TickerUpdate for WebSocket events
func (c *Converter) ConvertTickerToUpdate(okexTicker *market.Ticker) *commontypes.TickerUpdate {
	if okexTicker == nil {
		return nil
	}

	// Calculate 24h price change and percent change
	last := float64(okexTicker.Last)
	open24h := float64(okexTicker.Open24h)
	priceChange := last - open24h
	percentChange := 0.0
	if open24h != 0 {
		percentChange = (priceChange / open24h) * 100
	}

	return &commontypes.TickerUpdate{
		Symbol:           okexTicker.InstID,
		LastPrice:        commontypes.Decimal(strconv.FormatFloat(last, 'f', -1, 64)),
		BidPrice:         commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.BidPx), 'f', -1, 64)),
		BidSize:          commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.BidSz), 'f', -1, 64)),
		AskPrice:         commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.AskPx), 'f', -1, 64)),
		AskSize:          commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.AskSz), 'f', -1, 64)),
		High24h:          commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.High24h), 'f', -1, 64)),
		Low24h:           commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.Low24h), 'f', -1, 64)),
		Open24h:          commontypes.Decimal(strconv.FormatFloat(open24h, 'f', -1, 64)),
		Volume24h:        commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.Vol24h), 'f', -1, 64)),
		QuoteVolume24h:   commontypes.Decimal(strconv.FormatFloat(float64(okexTicker.VolCcy24h), 'f', -1, 64)),
		PriceChange24h:   commontypes.Decimal(strconv.FormatFloat(priceChange, 'f', -1, 64)),
		PercentChange24h: commontypes.Decimal(strconv.FormatFloat(percentChange, 'f', -1, 64)),
		Timestamp:        commontypes.Timestamp(okexTicker.TS),
		Extra: map[string]interface{}{
			"lastSz":   okexTicker.LastSz,
			"sodUtc0":  okexTicker.SodUtc0,
			"sodUtc8":  okexTicker.SodUtc8,
			"instType": okexTicker.InstType,
		},
	}
}

// ConvertCandleToUpdate converts an OKEx candle to common CandleUpdate for WebSocket events
func (c *Converter) ConvertCandleToUpdate(symbol, interval string, okexCandle *market.Candle) *commontypes.CandleUpdate {
	if okexCandle == nil {
		return nil
	}

	return &commontypes.CandleUpdate{
		Symbol:      symbol,
		Interval:    interval,
		Open:        commontypes.Decimal(strconv.FormatFloat(okexCandle.O, 'f', -1, 64)),
		High:        commontypes.Decimal(strconv.FormatFloat(okexCandle.H, 'f', -1, 64)),
		Low:         commontypes.Decimal(strconv.FormatFloat(okexCandle.L, 'f', -1, 64)),
		Close:       commontypes.Decimal(strconv.FormatFloat(okexCandle.C, 'f', -1, 64)),
		Volume:      commontypes.Decimal(strconv.FormatFloat(okexCandle.Vol, 'f', -1, 64)),
		QuoteVolume: commontypes.Decimal(strconv.FormatFloat(okexCandle.VolCcy, 'f', -1, 64)),
		Timestamp:   commontypes.Timestamp(okexCandle.TS),
		Confirmed:   okexCandle.Confirm,
		Extra: map[string]interface{}{
			"volCcyQuote": okexCandle.VolCcyQuote,
		},
	}
}

// ConvertCandle converts an OKEx candle to common Candle for REST API
func (c *Converter) ConvertCandle(okexCandle *market.Candle, symbol, interval string) *commontypes.Candle {
	if okexCandle == nil {
		return nil
	}

	return &commontypes.Candle{
		Symbol:      symbol,
		Interval:    interval,
		Open:        commontypes.Decimal(strconv.FormatFloat(okexCandle.O, 'f', -1, 64)),
		High:        commontypes.Decimal(strconv.FormatFloat(okexCandle.H, 'f', -1, 64)),
		Low:         commontypes.Decimal(strconv.FormatFloat(okexCandle.L, 'f', -1, 64)),
		Close:       commontypes.Decimal(strconv.FormatFloat(okexCandle.C, 'f', -1, 64)),
		Volume:      commontypes.Decimal(strconv.FormatFloat(okexCandle.Vol, 'f', -1, 64)),
		QuoteVolume: commontypes.Decimal(strconv.FormatFloat(okexCandle.VolCcy, 'f', -1, 64)),
		Timestamp:   commontypes.Timestamp(okexCandle.TS),
		Confirmed:   okexCandle.Confirm,
		Extra: map[string]interface{}{
			"volCcyQuote": okexCandle.VolCcyQuote,
		},
	}
}
