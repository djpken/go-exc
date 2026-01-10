package okex

import (
	"context"
	"strconv"

	accountreq "github.com/djpken/go-exc/exchanges/okex/requests/rest/account"
	tradereq "github.com/djpken/go-exc/exchanges/okex/requests/rest/trade"
	"github.com/djpken/go-exc/exchanges/okex/rest"
	okextypes "github.com/djpken/go-exc/exchanges/okex/types"
	commontypes "github.com/djpken/go-exc/types"
)

// RESTAdapter adapts OKEx REST client to common interface
type RESTAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// NewRESTAdapter creates a new REST adapter
func NewRESTAdapter(client *rest.ClientRest) *RESTAdapter {
	return &RESTAdapter{
		client:    client,
		converter: NewConverter(),
	}
}

// TradeAPI returns the trade API adapter
func (a *RESTAdapter) Trade() *TradeAPIAdapter {
	return &TradeAPIAdapter{
		client:    a.client,
		converter: a.converter,
	}
}

// AccountAPI returns the account API adapter
func (a *RESTAdapter) Account() *AccountAPIAdapter {
	return &AccountAPIAdapter{
		client:    a.client,
		converter: a.converter,
	}
}

// MarketAPI returns the market API adapter
func (a *RESTAdapter) Market() *MarketAPIAdapter {
	return &MarketAPIAdapter{
		client:    a.client,
		converter: a.converter,
	}
}

// FundingAPI returns the funding API adapter
func (a *RESTAdapter) Funding() *FundingAPIAdapter {
	return &FundingAPIAdapter{
		client:    a.client,
		converter: a.converter,
	}
}

// TradeAPIAdapter implements trading operations
type TradeAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// PlaceOrder places a new order
func (a *TradeAPIAdapter) PlaceOrder(ctx context.Context, symbol, side, orderType string, quantity, price float64, extra map[string]interface{}) (*commontypes.Order, error) {
	// Build OKEx order request
	req := tradereq.PlaceOrder{
		InstID:  symbol,
		TdMode:  okextypes.TradeCashMode, // Default to cash mode
		Side:    a.converter.ConvertOrderSide(side),
		OrdType: a.converter.ConvertOrderType(orderType),
		Sz:      quantity,
	}

	if price > 0 {
		req.Px = price
	}

	// Apply extra parameters
	if extra != nil {
		if tdMode, ok := extra["tdMode"].(string); ok {
			req.TdMode = okextypes.TradeMode(tdMode)
		}
		if clOrdID, ok := extra["clOrdID"].(string); ok {
			req.ClOrdID = clOrdID
		}
		if tag, ok := extra["tag"].(string); ok {
			req.Tag = tag
		}
	}

	// Place order (requires slice)
	resp, err := a.client.Trade.PlaceOrder([]tradereq.PlaceOrder{req})
	if err != nil {
		return nil, err
	}

	if len(resp.PlaceOrders) == 0 {
		return nil, nil
	}

	// Convert to common order type
	// Note: PlaceOrder response doesn't include full order details
	// You may need to call GetOrder to get complete information
	ordID := strconv.FormatFloat(float64(resp.PlaceOrders[0].OrdID), 'f', 0, 64)
	return &commontypes.Order{
		ID:     ordID,
		Symbol: symbol,
		Side:   side,
		Type:   orderType,
		Extra: map[string]interface{}{
			"clOrdID": resp.PlaceOrders[0].ClOrdID,
			"sCode":   resp.PlaceOrders[0].SCode,
			"sMsg":    resp.PlaceOrders[0].SMsg,
		},
	}, nil
}

// CancelOrder cancels an existing order
func (a *TradeAPIAdapter) CancelOrder(ctx context.Context, symbol, orderId string, extra map[string]interface{}) error {
	req := tradereq.CancelOrder{
		InstID: symbol,
		OrdID:  orderId,
	}

	if extra != nil {
		if clOrdID, ok := extra["clOrdID"].(string); ok {
			req.ClOrdID = clOrdID
		}
	}

	_, err := a.client.Trade.CandleOrder([]tradereq.CancelOrder{req})
	return err
}

// GetOrder gets order details
func (a *TradeAPIAdapter) GetOrder(ctx context.Context, symbol, orderId string, extra map[string]interface{}) (*commontypes.Order, error) {
	req := tradereq.OrderDetails{
		InstID: symbol,
		OrdID:  orderId,
	}

	if extra != nil {
		if clOrdID, ok := extra["clOrdID"].(string); ok {
			req.ClOrdID = clOrdID
		}
	}

	resp, err := a.client.Trade.GetOrderDetail(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Orders) == 0 {
		return nil, nil
	}

	return a.converter.ConvertOrder(resp.Orders[0]), nil
}

// AccountAPIAdapter implements account operations
type AccountAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetBalance gets account balance
func (a *AccountAPIAdapter) GetBalance(ctx context.Context, currencies ...string) (*commontypes.AccountBalance, error) {
	req := accountreq.GetBalance{
		Ccy: currencies,
	}

	resp, err := a.client.Account.GetBalance(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Balances) == 0 {
		return nil, nil
	}

	return a.converter.ConvertBalance(resp.Balances[0]), nil
}

// GetPositions gets account positions
func (a *AccountAPIAdapter) GetPositions(ctx context.Context, symbols ...string) ([]*commontypes.Position, error) {
	req := accountreq.GetPositions{
		InstID: symbols,
	}

	resp, err := a.client.Account.GetPositions(req)
	if err != nil {
		return nil, err
	}

	positions := make([]*commontypes.Position, 0, len(resp.Positions))
	for _, pos := range resp.Positions {
		positions = append(positions, a.converter.ConvertPosition(pos))
	}

	return positions, nil
}

// MarketAPIAdapter implements market data operations
type MarketAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetTicker gets ticker information
func (a *MarketAPIAdapter) GetTicker(ctx context.Context, symbol string) (*commontypes.Ticker, error) {
	// TODO: Implement ticker conversion
	// This requires reading OKEx market API response structure
	return nil, nil
}

// GetOrderBook gets order book
func (a *MarketAPIAdapter) GetOrderBook(ctx context.Context, symbol string, depth int) (*commontypes.OrderBook, error) {
	// TODO: Implement order book conversion
	// This requires reading OKEx market API response structure
	return nil, nil
}

// FundingAPIAdapter implements funding operations
type FundingAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetDepositAddress gets deposit address
func (a *FundingAPIAdapter) GetDepositAddress(ctx context.Context, currency string) (string, error) {
	// TODO: Implement deposit address retrieval
	// This requires OKEx funding API
	return "", nil
}

// Withdraw initiates a withdrawal
func (a *FundingAPIAdapter) Withdraw(ctx context.Context, currency string, amount float64, address, tag string, extra map[string]interface{}) (string, error) {
	// TODO: Implement withdrawal
	// This requires OKEx funding API
	return "", nil
}
