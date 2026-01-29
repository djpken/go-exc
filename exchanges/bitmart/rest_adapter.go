package bitmart

import (
	"context"
	"fmt"

	accountreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/account"
	"github.com/djpken/go-exc/exchanges/bitmart/requests/rest/contract"
	fundingreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/funding"
	marketreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/market"
	tradereq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/trade"
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

	return a.converter.ConvertAccountBalance(balances.Data.Wallet), nil
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

// GetTickers gets ticker information for all symbols
func (a *MarketAPIAdapter) GetTickers(ctx context.Context) ([]*commontypes.Ticker, error) {
	return nil, commontypes.ErrNotSupported
}

// GetInstruments gets trading instrument information
func (a *MarketAPIAdapter) GetInstruments(ctx context.Context) ([]*commontypes.Instrument, error) {
	// Get all symbols
	resp, err := a.client.Contract.GetContractDetails(
		contract.GetContractDetailsRequest{})
	if err != nil {
		return nil, err
	}

	// For each symbol, get detailed information
	instruments := make([]*commontypes.Instrument, 0, len(resp.Data.Symbols))
	for _, symbol := range resp.Data.Symbols {

		instrument := a.converter.ConvertInstrument(&symbol)
		if instrument != nil {
			instruments = append(instruments, instrument)
		}
	}

	return instruments, nil
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

// GetCandles gets historical candlestick/kline data
func (a *MarketAPIAdapter) GetCandles(ctx context.Context, req commontypes.GetCandlesRequest) ([]*commontypes.Candle, error) {
	// Convert interval to BitMart step (in minutes)
	step, err := a.converter.ConvertIntervalToStep(req.Interval)
	if err != nil {
		return nil, err
	}

	// Build BitMart request
	bitmartReq := marketreq.GetKlineRequest{
		Symbol: req.Symbol,
		Step:   step,
		Limit:  req.Limit,
	}

	// Set time range if specified
	if req.StartTime != nil {
		bitmartReq.FromTime = req.StartTime.Unix()
	}
	if req.EndTime != nil {
		bitmartReq.ToTime = req.EndTime.Unix()
	}

	// Call BitMart API
	resp, err := a.client.Market.GetKlines(bitmartReq)
	if err != nil {
		return nil, err
	}

	// Convert response to common types
	// BitMart returns: [timestamp, open, high, low, close, volume, quote_volume]
	candles := make([]*commontypes.Candle, 0, len(resp.Data))
	for _, klineData := range resp.Data {
		candle := a.converter.ConvertKlineArrayToCandle(klineData, req.Symbol, req.Interval)
		if candle != nil {
			candles = append(candles, candle)
		}
	}

	return candles, nil
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
