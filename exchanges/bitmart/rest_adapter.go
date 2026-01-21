package bitmart

import (
	"context"
	"fmt"

	tradereq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/trade"
	accountreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/account"
	marketreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/market"
	fundingreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/funding"
	"github.com/djpken/go-exc/exchanges/bitmart/rest"
	commontypes "github.com/djpken/go-exc/types"
)

// RESTAdapter adapts BitMart REST client to common interface
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

// Trade returns the trade API adapter
func (a *RESTAdapter) Trade() *TradeAPIAdapter {
	return &TradeAPIAdapter{
		client:    a.client,
		converter: a.converter,
	}
}

// Account returns the account API adapter
func (a *RESTAdapter) Account() *AccountAPIAdapter {
	return &AccountAPIAdapter{
		client:    a.client,
		converter: a.converter,
	}
}

// Market returns the market API adapter
func (a *RESTAdapter) Market() *MarketAPIAdapter {
	return &MarketAPIAdapter{
		client:    a.client,
		converter: a.converter,
	}
}

// Funding returns the funding API adapter
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
	req := tradereq.PlaceOrderRequest{
		Symbol: symbol,
		Side:   a.converter.ConvertOrderSide(side),
		Type:   a.converter.ConvertOrderType(orderType),
		Size:   a.converter.formatFloat(quantity),
	}

	if price > 0 {
		req.Price = a.converter.formatFloat(price)
	}

	if extra != nil {
		if clientOrderID, ok := extra["clientOrderID"].(string); ok {
			req.ClientOrderID = clientOrderID
		}
	}

	resp, err := a.client.Trade.PlaceOrder(req)
	if err != nil {
		return nil, err
	}

	// Query order details
	orderReq := tradereq.GetOrderRequest{OrderID: resp.Data.OrderID}
	orderResp, err := a.client.Trade.GetOrder(orderReq)
	if err != nil {
		return nil, err
	}

	return a.converter.ConvertOrderDetail(&orderResp.Data), nil
}

// CancelOrder cancels an existing order
func (a *TradeAPIAdapter) CancelOrder(ctx context.Context, symbol, orderID string, extra map[string]interface{}) error {
	req := tradereq.CancelOrderRequest{
		Symbol:  symbol,
		OrderID: orderID,
	}

	_, err := a.client.Trade.CancelOrder(req)
	return err
}

// GetOrder gets order details
func (a *TradeAPIAdapter) GetOrder(ctx context.Context, symbol, orderID string, extra map[string]interface{}) (*commontypes.Order, error) {
	req := tradereq.GetOrderRequest{OrderID: orderID}
	resp, err := a.client.Trade.GetOrder(req)
	if err != nil {
		return nil, err
	}

	return a.converter.ConvertOrderDetail(&resp.Data), nil
}

// AccountAPIAdapter implements account operations
type AccountAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetBalance gets account balance
func (a *AccountAPIAdapter) GetBalance(ctx context.Context) (*commontypes.AccountBalance, error) {
	req := accountreq.GetWalletBalanceRequest{}
	balances, err := a.client.Account.GetWalletBalance(req)
	if err != nil {
		return nil, err
	}

	return a.converter.ConvertAccountBalance(balances.Data.Balances), nil
}

// GetPositions gets account positions (not supported for spot trading)
func (a *AccountAPIAdapter) GetPositions(ctx context.Context) ([]*commontypes.Position, error) {
	return []*commontypes.Position{}, nil
}

// MarketAPIAdapter implements market data operations
type MarketAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetTicker gets ticker information
func (a *MarketAPIAdapter) GetTicker(ctx context.Context, symbol string) (*commontypes.Ticker, error) {
	req := marketreq.GetTickerRequest{Symbol: symbol}
	ticker, err := a.client.Market.GetTicker(req)
	if err != nil {
		return nil, err
	}

	return a.converter.ConvertTicker(&ticker.Data), nil
}

// GetOrderBook gets order book
func (a *MarketAPIAdapter) GetOrderBook(ctx context.Context, symbol string, depth int) (*commontypes.OrderBook, error) {
	validDepth := 20 // default
	switch depth {
	case 5, 20, 50:
		validDepth = depth
	}

	req := marketreq.GetOrderBookRequest{
		Symbol: symbol,
		Depth:  validDepth,
	}

	orderBook, err := a.client.Market.GetOrderBook(req)
	if err != nil {
		return nil, err
	}

	return a.converter.ConvertOrderBook(&orderBook.Data, symbol), nil
}

// FundingAPIAdapter implements funding operations
type FundingAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetDepositAddress gets deposit address
func (a *FundingAPIAdapter) GetDepositAddress(ctx context.Context, currency string) (string, error) {
	req := fundingreq.GetDepositAddressRequest{Currency: currency}
	address, err := a.client.Funding.GetDepositAddress(req)
	if err != nil {
		return "", err
	}

	if address.Data.Address != "" {
		if address.Data.Tag != "" {
			return fmt.Sprintf("%s:%s", address.Data.Address, address.Data.Tag), nil
		}
		return address.Data.Address, nil
	}

	return "", fmt.Errorf("no deposit address found for currency %s", currency)
}

// Withdraw initiates a withdrawal
func (a *FundingAPIAdapter) Withdraw(ctx context.Context, currency string, amount float64, address, tag string, extra map[string]interface{}) (string, error) {
	req := fundingreq.WithdrawRequest{
		Currency: currency,
		Amount:   a.converter.formatFloat(amount),
		Address:  address,
	}

	if tag != "" {
		req.AddressTag = tag
	}

	resp, err := a.client.Funding.Withdraw(req)
	if err != nil {
		return "", err
	}

	return resp.Data.WithdrawID, nil
}
