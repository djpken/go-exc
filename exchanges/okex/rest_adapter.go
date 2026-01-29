package okex

import (
	"context"
	"fmt"
	"strconv"

	okexconstants "github.com/djpken/go-exc/exchanges/okex/constants"
	"github.com/djpken/go-exc/exchanges/okex/models/market"
	accountreq "github.com/djpken/go-exc/exchanges/okex/requests/rest/account"
	marketreq "github.com/djpken/go-exc/exchanges/okex/requests/rest/market"
	publicreq "github.com/djpken/go-exc/exchanges/okex/requests/rest/public"
	tradereq "github.com/djpken/go-exc/exchanges/okex/requests/rest/trade"
	"github.com/djpken/go-exc/exchanges/okex/responses"
	"github.com/djpken/go-exc/exchanges/okex/rest"
	commontypes "github.com/djpken/go-exc/types"
)

// checkAPIError checks if the API response contains an error
// OKEx returns Code > 0 for errors
func checkAPIError(basic responses.Basic) error {
	if basic.Code > 0 {
		return fmt.Errorf("OKEx API error: code=%d, msg=%s", basic.Code, basic.Msg)
	}
	return nil
}

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
		TdMode:  okexconstants.TradeCashMode, // Default to cash mode
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
			req.TdMode = okexconstants.TradeMode(tdMode)
		}
		if posSide, ok := extra["posSide"].(string); ok {
			req.PosSide = okexconstants.PositionSide(posSide)
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

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	if len(resp.PlaceOrders) == 0 {
		return nil, fmt.Errorf("no order data returned")
	}

	// Check if the order placement was successful
	if resp.PlaceOrders[0].SCode != 0 {
		return nil, fmt.Errorf("order placement failed: code=%d, msg=%s", resp.PlaceOrders[0].SCode, resp.PlaceOrders[0].SMsg)
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

	// Note: CandleOrder is a typo in the original SDK, should be CancelOrder
	// It returns PlaceOrder response type which is reused for cancel responses
	resp, err := a.client.Trade.CandleOrder([]tradereq.CancelOrder{req})
	if err != nil {
		return err
	}

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return err
	}

	// Check if the cancellation was successful
	// The response reuses PlaceOrders field name
	if len(resp.PlaceOrders) > 0 && resp.PlaceOrders[0].SCode != 0 {
		return fmt.Errorf("order cancellation failed: code=%d, msg=%s", resp.PlaceOrders[0].SCode, resp.PlaceOrders[0].SMsg)
	}

	return nil
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

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	if len(resp.Orders) == 0 {
		return nil, fmt.Errorf("order not found: symbol=%s, orderId=%s", symbol, orderId)
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

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	if len(resp.Balances) == 0 {
		// Return empty balance instead of nil
		return &commontypes.AccountBalance{
			Balances: []*commontypes.Balance{},
		}, nil
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

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	positions := make([]*commontypes.Position, 0, len(resp.Positions))
	for _, pos := range resp.Positions {
		positions = append(positions, a.converter.ConvertPosition(pos))
	}

	return positions, nil
}

// GetLeverage gets leverage configuration
func (a *AccountAPIAdapter) GetLeverage(ctx context.Context, req commontypes.GetLeverageRequest) ([]*commontypes.Leverage, error) {
	okexReq := accountreq.GetLeverage{
		InstID:  req.Symbols,
		MgnMode: okexconstants.MarginMode(req.MarginMode),
	}

	resp, err := a.client.Account.GetLeverage(okexReq)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	leverages := make([]*commontypes.Leverage, 0, len(resp.Leverages))
	for _, lev := range resp.Leverages {
		leverages = append(leverages, &commontypes.Leverage{
			Symbol:       lev.InstID,
			Leverage:     int(lev.Lever),
			MarginMode:   string(lev.MgnMode),
			PositionSide: string(lev.PosSide),
			Extra:        map[string]interface{}{},
		})
	}

	return leverages, nil
}

// SetLeverage sets leverage for a trading pair
func (a *AccountAPIAdapter) SetLeverage(ctx context.Context, req commontypes.SetLeverageRequest) (*commontypes.Leverage, error) {
	okexReq := accountreq.SetLeverage{
		Lever:   int64(req.Leverage),
		MgnMode: okexconstants.MarginMode(req.MarginMode),
	}

	// OKEx supports either InstID (symbol) or Ccy (currency)
	if req.Symbol != "" {
		okexReq.InstID = req.Symbol
	} else if req.Currency != "" {
		okexReq.Ccy = req.Currency
	}

	if req.PositionSide != "" {
		okexReq.PosSide = okexconstants.PositionSide(req.PositionSide)
	}

	resp, err := a.client.Account.SetLeverage(okexReq)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	if len(resp.Leverages) == 0 {
		return nil, fmt.Errorf("no leverage data returned")
	}

	lev := resp.Leverages[0]
	return &commontypes.Leverage{
		Symbol:       lev.InstID,
		Leverage:     int(lev.Lever),
		MarginMode:   string(lev.MgnMode),
		PositionSide: string(lev.PosSide),
		Extra:        map[string]interface{}{},
	}, nil
}

// MarketAPIAdapter implements market data operations
type MarketAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetTicker gets ticker information
func (a *MarketAPIAdapter) GetTicker(ctx context.Context, symbol string) (*commontypes.Ticker, error) {
	// Build request
	req := marketreq.GetTickers{
		InstType: okexconstants.SpotInstrument, // Default to spot
	}

	// Call OKEx API
	resp, err := a.client.Market.GetTicker(req)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	// Check if we have tickers
	if len(resp.Tickers) == 0 {
		return nil, fmt.Errorf("no ticker data returned for symbol %s", symbol)
	}

	// Use converter to convert ticker
	return a.converter.ConvertTicker(resp.Tickers[0]), nil
}

// GetTickers gets ticker information for multiple instruments
func (a *MarketAPIAdapter) GetTickers(ctx context.Context, req commontypes.GetTickersRequest) ([]*commontypes.Ticker, error) {
	// Build OKEx request
	okexReq := marketreq.GetTickers{
		InstType: a.converter.ConvertInstrumentType(req.InstrumentType),
	}

	// Call OKEx API
	resp, err := a.client.Market.GetTickers(okexReq)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	// Convert all tickers
	tickers := make([]*commontypes.Ticker, 0, len(resp.Tickers))
	for _, okexTicker := range resp.Tickers {
		tickers = append(tickers, a.converter.ConvertTicker(okexTicker))
	}

	return tickers, nil
}

// GetInstruments gets trading instrument information
func (a *MarketAPIAdapter) GetInstruments(ctx context.Context, req commontypes.GetInstrumentsRequest) ([]*commontypes.Instrument, error) {
	// Build OKEx request
	okexReq := publicreq.GetInstruments{
		InstType: a.converter.ConvertInstrumentType(req.InstrumentType),
	}

	// Call OKEx API
	resp, err := a.client.PublicData.GetInstruments(okexReq)
	if err != nil {
		return nil, err
	}
	tickerResp, err := a.client.Market.GetTickers(marketreq.GetTickers{
		InstType: okexReq.InstType,
	})
	if err != nil {
		return nil, err
	}

	m := make(map[string]market.Ticker, len(tickerResp.Tickers))
	for _, ticker := range tickerResp.Tickers {
		m[ticker.InstID] = *ticker
	}

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	// Convert all instruments
	instruments := make([]*commontypes.Instrument, 0, len(resp.Instruments))
	for _, okexInst := range resp.Instruments {
		instruments = append(instruments, a.converter.ConvertInstrument(okexInst, m[okexInst.InstID]))
	}

	return instruments, nil
}

// GetCandles gets historical candlestick/kline data
func (a *MarketAPIAdapter) GetCandles(ctx context.Context, req commontypes.GetCandlesRequest) ([]*commontypes.Candle, error) {
	// Build OKEx request
	okexReq := marketreq.GetCandlesticks{
		InstID: req.Symbol,
		Bar:    okexconstants.BarSize(req.Interval), // OKEx uses the same format as our common interval
	}

	// Set limit if specified
	if req.Limit > 0 {
		okexReq.Limit = int64(req.Limit)
	}

	// Set time range if specified
	// OKEx uses Before/After for pagination
	// Before: timestamp to request data before (exclusive)
	// After: timestamp to request data after (exclusive)
	if req.StartTime != nil {
		// After parameter: get data after this timestamp
		okexReq.After = req.StartTime.UnixMilli()
	}
	if req.EndTime != nil {
		// Before parameter: get data before this timestamp
		okexReq.Before = req.EndTime.UnixMilli()
	}

	// Call OKEx API
	resp, err := a.client.Market.GetCandlesticks(okexReq)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if err := checkAPIError(resp.Basic); err != nil {
		return nil, err
	}

	// Convert response to common types
	candles := make([]*commontypes.Candle, 0, len(resp.Candles))
	for _, okexCandle := range resp.Candles {
		candle := a.converter.ConvertCandle(okexCandle, req.Symbol, req.Interval)
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

// Note: GetDepositAddress and Withdraw methods have been removed from the Exchange interface
// If you need these features, use the native REST client directly via exchange.REST().Funding()
