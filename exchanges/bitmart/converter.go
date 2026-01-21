package bitmart

import (
	"fmt"
	"time"

	accountmodels "github.com/djpken/go-exc/exchanges/bitmart/models/account"
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
		ID:               order.OrderID,
		Symbol:           order.Symbol,
		Side:             order.Side,
		Type:             order.Type,
		Quantity:         commontypes.Decimal(order.Size),
		Price:            commontypes.Decimal(order.Price),
		FilledQuantity:   commontypes.Decimal(order.FilledSize),
		Status:           c.ConvertOrderStatus(order.Status),
		CreatedAt:        commontypes.Timestamp(time.Unix(0, order.CreateTime*int64(time.Millisecond))),
		UpdatedAt:        commontypes.Timestamp(time.Unix(0, order.UpdateTime*int64(time.Millisecond))),
		Extra:            map[string]interface{}{"notional": order.Notional},
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

	return &commontypes.Balance{
		Currency:  balance.Currency,
		Available: commontypes.Decimal(balance.Available),
		Frozen:    commontypes.Decimal(balance.Frozen),
		Total:     commontypes.Decimal(balance.Total),
	}
}

// ConvertAccountBalance converts BitMart balances to common account balance
func (c *Converter) ConvertAccountBalance(balances []accountmodels.Balance) *commontypes.AccountBalance {
	if balances == nil {
		return nil
	}

	commonBalances := make([]*commontypes.Balance, 0, len(balances))
	for _, bal := range balances {
		commonBalances = append(commonBalances, &commontypes.Balance{
			Currency:  bal.Currency,
			Available: commontypes.Decimal(bal.Available),
			Frozen:    commontypes.Decimal(bal.Frozen),
			Total:     commontypes.Decimal(bal.Total),
		})
	}

	return &commontypes.AccountBalance{
		Balances: commonBalances,
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
		Timestamp int64                       `json:"timestamp"`
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
