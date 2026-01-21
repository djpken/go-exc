package okex

import (
	"strconv"
	"time"

	"github.com/djpken/go-exc/exchanges/okex/models/account"
	"github.com/djpken/go-exc/exchanges/okex/models/trade"
	okexconstants "github.com/djpken/go-exc/exchanges/okex/constants"
	commontypes "github.com/djpken/go-exc/types"
)

// Converter handles conversion between OKEx types and common types
type Converter struct{}

// NewConverter creates a new Converter
func NewConverter() *Converter {
	return &Converter{}
}

// ConvertOrder converts OKEx order to common Order type
func (c *Converter) ConvertOrder(okexOrder *trade.Order) *commontypes.Order {
	if okexOrder == nil {
		return nil
	}

	price := float64(okexOrder.Px)
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
		CreatedAt:         commontypes.Timestamp(time.Time(okexOrder.CTime)),
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
		UpdatedAt:   commontypes.Timestamp(time.Time(okexBalance.UTime)),
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

// ConvertOrderSide converts common order side to OKEx order side
func (c *Converter) ConvertOrderSide(side string) okexconstants.OrderSide {
	switch side {
	case "buy":
		return okexconstants.OrderBuy
	case "sell":
		return okexconstants.OrderSell
	default:
		return okexconstants.OrderBuy
	}
}

// ConvertOrderType converts common order type to OKEx order type
func (c *Converter) ConvertOrderType(orderType string) okexconstants.OrderType {
	switch orderType {
	case "market":
		return okexconstants.OrderMarket
	case "limit":
		return okexconstants.OrderLimit
	case "post_only":
		return okexconstants.OrderPostOnly
	case "fok":
		return okexconstants.OrderFOK
	case "ioc":
		return okexconstants.OrderIOC
	default:
		return okexconstants.OrderLimit
	}
}

// ConvertAccountConfig converts OKEx account config to common AccountConfig type
func (c *Converter) ConvertAccountConfig(okexConfig *account.Config) *commontypes.AccountConfig {
	if okexConfig == nil {
		return nil
	}

	// Convert position mode
	posMode := "net_mode"
	if okexConfig.PosMode == okexconstants.PosModeType("long_short_mode") {
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
