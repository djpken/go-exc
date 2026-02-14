package bitmart

import (
	"context"
	"fmt"
	"time"

	accountreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/account"
	contractreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/contract"
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
// Routes to either spot or contract API based on parameters
func (a *TradeAPIAdapter) PlaceOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	// Determine if this is a contract or spot order
	isContract := false
	posSide := req.PosSide

	// Check if position side is specified (contract trading)
	if posSide != "" && posSide != commontypes.PositionSideNet {
		isContract = true
	}

	extra := req.Extra
	if extra != nil {
		if _, hasLeverage := extra["leverage"]; hasLeverage {
			isContract = true
		}
		if _, hasOpenType := extra["open_type"]; hasOpenType {
			isContract = true
		}
		if tdMode, hasTdMode := extra["tdMode"]; hasTdMode && tdMode != "" {
			isContract = true
		}
	}

	// Route to appropriate API
	if isContract {
		return a.placeContractOrder(ctx, req)
	}
	return a.placeSpotOrder(ctx, req)
}

// placeSpotOrder places a spot trading order
func (a *TradeAPIAdapter) placeSpotOrder(ctx context.Context, commonReq commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	req := tradereq.PlaceOrderRequest{
		Symbol:        commonReq.Symbol,
		ClientOrderID: commonReq.ClientOrderID,
		Side:          a.converter.ConvertOrderSide(string(commonReq.Side)),
		Type:          a.converter.ConvertOrderType(commonReq.Type),
		Size:          a.converter.formatFloat(commonReq.Quantity),
	}

	price := commonReq.Price
	if price > 0 {
		req.Price = a.converter.formatFloat(price)
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

// placeContractOrder places a contract/futures order
func (a *TradeAPIAdapter) placeContractOrder(ctx context.Context, commonReq commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	// Convert side and position side to BitMart contract side value
	contractSide := a.converter.ConvertToContractSide(commonReq.Side, commonReq.PosSide)

	req := contractreq.SubmitOrderRequest{
		Symbol:        commonReq.Symbol,
		ClientOrderID: commonReq.ClientOrderID,
		Side:          contractSide,
		Type:          a.converter.ConvertOrderType(commonReq.Type),
		Size:          int(commonReq.Quantity), // Contract orders use integer size
	}

	price := commonReq.Price
	if price > 0 {
		req.Price = a.converter.formatFloat(price)
	}

	extra := commonReq.Extra
	// Handle contract-specific parameters from extra
	if extra != nil {
		if leverage, ok := extra["leverage"].(string); ok {
			req.Leverage = leverage
		}
		if openType, ok := extra["open_type"].(string); ok {
			req.OpenType = openType
		} else if tdMode, ok := extra["tdMode"].(commontypes.MarginMode); ok {
			// Convert MarginMode to BitMart open_type
			switch tdMode {
			case commontypes.MarginModeCross:
				req.OpenType = "cross"
			case commontypes.MarginModeIsolated:
				req.OpenType = "isolated"
			}
		}
		if mode, ok := extra["mode"].(int); ok {
			req.Mode = mode
		}
	}

	resp, err := a.client.Contract.SubmitOrder(req)
	if err != nil {
		return nil, err
	}

	// Convert contract order response to common Order type
	return a.converter.ConvertContractOrder(resp), nil
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

// GetOrderDetail gets order details
// Supports both spot and contract orders
//
// Usage:
//   - Spot order: GetOrderDetail(ctx, symbol, orderID, nil)
//   - Contract order: GetOrderDetail(ctx, symbol, orderID, map[string]interface{}{"account_type": types.AccountTypeFutures})
//
// For contract orders, returns aggregated trade executions as an Order object
func (a *TradeAPIAdapter) GetOrderDetail(ctx context.Context, commonReq commontypes.GetOrderRequest) (*commontypes.Order, error) {
	// Check if this is a contract order
	accountType := commontypes.AccountTypeSpot // default to spot
	extra := commonReq.Extra
	if extra != nil {
		if accType, ok := extra["account_type"].(string); ok {
			accountType = accType
		}
	}
	switch accountType {
	case commontypes.AccountTypeFutures:
		// Query contract trades for this order
		req := contractreq.GetContractTradesRequest{}
		if commonReq.Symbol != "" {
			req.Symbol = commonReq.Symbol
		}
		if commonReq.OrderID != "" {
			req.OrderID = commonReq.OrderID
		}
		if commonReq.ClientOrderID != "" {
			req.ClientOrderID = commonReq.ClientOrderID
		}
		resp, err := a.client.Contract.GetContractTrades(req)
		if err != nil {
			return nil, err
		}

		if len(resp.Data) == 0 {
			return nil, fmt.Errorf("no trades found for order %s", commonReq.OrderID)
		}

		// Convert contract trades to Order
		return a.converter.ConvertContractTrades(resp.Data), nil

	case commontypes.AccountTypeSpot:
		// Query spot order detail (existing implementation)
		req := tradereq.GetOrderRequest{OrderID: commonReq.OrderID, ClientOrderID: commonReq.ClientOrderID}
		resp, err := a.client.Trade.GetOrder(req)
		if err != nil {
			return nil, err
		}

		return a.converter.ConvertOrderDetail(&resp.Data), nil

	default:
		return nil, commontypes.ErrNotSupported
	}
}

// AccountAPIAdapter implements account operations
type AccountAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

// GetBalance gets account balance
// Supports both spot and futures accounts via account type indicator
//
// Usage:
//   - Spot account (default): GetBalance(ctx)
//   - Futures account: GetBalance(ctx, types.AccountTypeFutures, "USDT", "BTC", ...)
//   - Filter currencies: GetBalance(ctx, "USDT", "BTC") for spot account
//
// The first parameter can be a special account type constant (types.AccountTypeFutures, types.AccountTypeSpot)
// to specify which account to query. If no account type is specified, spot account is queried by default.
func (a *AccountAPIAdapter) GetBalance(ctx context.Context, typee string, currencies ...string) (*commontypes.AccountBalance, error) {
	currencyList := currencies

	// Query based on account type
	switch typee {
	case commontypes.AccountTypeFutures:
		// Query futures/contract account
		resp, err := a.client.Contract.GetContractAssets()
		if err != nil {
			return nil, err
		}

		// Convert contract assets to account balance
		totalEquity := commontypes.ZeroDecimal
		balances := make([]*commontypes.Balance, 0)
		for _, asset := range resp.Data {
			// Filter by currency if specified
			if len(currencyList) > 0 {
				found := false
				for _, currency := range currencyList {
					if asset.Currency == currency {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			available := a.converter.stringToDecimal(asset.AvailableBalance)
			frozen := a.converter.stringToDecimal(asset.FrozenBalance)
			total, _ := available.Add(frozen)

			balance := &commontypes.Balance{
				Currency:        asset.Currency,
				Available:       available,
				Frozen:          frozen,
				Total:           total,
				PositionDeposit: a.converter.stringToDecimal(asset.PositionDeposit),
				Extra: map[string]interface{}{
					"position_deposit": asset.PositionDeposit,
					"equity":           asset.Equity,
					"unrealized":       asset.Unrealized,
				},
			}
			if _, err := totalEquity.Add(available); err != nil {
				return nil, err
			}
			balances = append(balances, balance)
		}

		return &commontypes.AccountBalance{
			Balances:    balances,
			TotalEquity: totalEquity,
			UpdatedAt:   commontypes.Timestamp{}, // BitMart doesn't provide timestamp in this API
			Extra: map[string]interface{}{
				"account_type": "futures",
			},
		}, nil

	case commontypes.AccountTypeSpot:
		// Query spot account (existing implementation)
		req := accountreq.GetWalletBalanceRequest{}
		balances, err := a.client.Account.GetWalletBalance(req)
		if err != nil {
			return nil, err
		}

		result := a.converter.ConvertAccountBalance(balances.Data.Wallet)

		// Filter by currency if specified
		if len(currencyList) > 0 {
			filteredBalances := make([]*commontypes.Balance, 0)
			for _, balance := range result.Balances {
				for _, currency := range currencyList {
					if balance.Currency == currency {
						filteredBalances = append(filteredBalances, balance)
						break
					}
				}
			}
			result.Balances = filteredBalances
		}

		if result.Extra == nil {
			result.Extra = make(map[string]interface{})
		}
		result.Extra["account_type"] = "spot"

		return result, nil

	default:
		return nil, commontypes.ErrNotSupported
	}
}

// GetPositions gets account positions (contract/futures positions)
// Uses BitMart's position-v2 API to query contract positions
//
// Notes:
// - Returns only positions with non-zero quantity
// - Supports both one-way and hedge position modes
// - For one-way mode: position_side is "both", direction determined by quantity sign
// - For hedge mode: separate long/short positions
func (a *AccountAPIAdapter) GetPositions(ctx context.Context, symbols ...string) ([]*commontypes.Position, error) {
	var allPositions []*commontypes.Position

	if len(symbols) > 0 {
		// Query specific symbols
		for _, symbol := range symbols {
			req := contractreq.GetPositionV2Request{
				Symbol: symbol,
			}

			resp, err := a.client.Contract.GetPositionV2(req)
			if err != nil {
				return nil, fmt.Errorf("failed to get positions for %s: %w", symbol, err)
			}

			// Convert positions
			for _, pos := range resp.Data {
				convertedPos := a.converter.ConvertPositionV2ToPosition(&pos)
				if convertedPos != nil {
					allPositions = append(allPositions, convertedPos)
				}
			}
		}
	} else {
		// Query all positions (only returns positions with non-zero quantity)
		req := contractreq.GetPositionV2Request{}
		resp, err := a.client.Contract.GetPositionV2(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get positions: %w", err)
		}

		// Convert positions
		for _, pos := range resp.Data {
			convertedPos := a.converter.ConvertPositionV2ToPosition(&pos)
			if convertedPos != nil {
				allPositions = append(allPositions, convertedPos)
			}
		}
	}

	return allPositions, nil
}

// GetLeverage gets leverage configuration for positions
func (a *AccountAPIAdapter) GetLeverage(ctx context.Context, symbols []string) ([]*commontypes.Leverage, error) {
	// If specific symbols are requested, query each one
	if len(symbols) > 0 {
		leverages := make([]*commontypes.Leverage, 0, len(symbols))
		for _, symbol := range symbols {
			req := contractreq.GetPositionV2Request{
				Symbol: symbol,
			}
			resp, err := a.client.Contract.GetPositionV2(req)
			if err != nil {
				return nil, err
			}

			// Convert each position to leverage
			for i := range resp.Data {
				leverage := a.converter.ConvertPositionV2ToLeverage(&resp.Data[i])
				if leverage != nil {
					leverages = append(leverages, leverage)
				}
			}
		}
		return leverages, nil
	}

	// If no specific symbols, get all positions
	req := contractreq.GetPositionV2Request{}
	resp, err := a.client.Contract.GetPositionV2(req)
	if err != nil {
		return nil, err
	}

	// Convert positions to leverages
	leverages := make([]*commontypes.Leverage, 0, len(resp.Data))
	for i := range resp.Data {
		leverage := a.converter.ConvertPositionV2ToLeverage(&resp.Data[i])
		if leverage != nil {
			leverages = append(leverages, leverage)
		}
	}

	return leverages, nil
}

// SetLeverage sets leverage for a contract trading pair
func (a *AccountAPIAdapter) SetLeverage(ctx context.Context, req commontypes.SetLeverageRequest) (*commontypes.Leverage, error) {
	// Convert common types to BitMart types
	var openType string
	switch req.MarginMode {
	case commontypes.MarginModeCross:
		openType = "cross"
	case commontypes.MarginModeIsolated:
		openType = "isolated"
	default:
		openType = "isolated" // Default to isolated
	}

	// Build BitMart request
	leverageReq := contractreq.SubmitLeverageRequest{
		Symbol:   req.Symbol,
		Leverage: fmt.Sprintf("%d", req.Leverage),
		OpenType: openType,
	}

	// Call BitMart API
	resp, err := a.client.Contract.SubmitLeverage(leverageReq)
	if err != nil {
		return nil, err
	}

	// Convert response to common Leverage type
	return a.converter.ConvertSubmitLeverageResponse(resp), nil
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
		contractreq.GetContractDetailsRequest{})
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
// Supports both spot and contract markets via account_type parameter in Extra
//
// Usage:
//   - Spot klines (default): GetCandles(ctx, GetCandlesRequest{Symbol: "BTC_USDT", Interval: "1m"})
//   - Contract klines: GetCandles(ctx, GetCandlesRequest{Symbol: "BTCUSDT", Interval: "1m", Extra: map[string]interface{}{"account_type": types.AccountTypeFutures}})
//
// Note: Contract symbols use no separator (BTCUSDT), spot symbols use underscore (BTC_USDT)
func (a *MarketAPIAdapter) GetCandles(ctx context.Context, req commontypes.GetCandlesRequest) ([]*commontypes.Candle, error) {
	// Convert interval to BitMart step (in minutes)
	step, err := a.converter.ConvertIntervalToStep(req.Interval)
	if err != nil {
		return nil, err
	}

	// Check if this is a contract kline request
	accountType := commontypes.AccountTypeSpot // default to spot
	if req.Extra != nil {
		if accType, ok := req.Extra["account_type"].(string); ok {
			accountType = accType
		}
	}

	switch accountType {
	case commontypes.AccountTypeFutures:
		// Query contract klines
		var startTime, endTime int64
		var limit int

		// Get limit (default to 100 if not specified)
		if req.Limit > 0 {
			limit = req.Limit
		} else {
			limit = 100 // default limit
		}

		// Ensure limit doesn't exceed API maximum
		if limit > 500 {
			limit = 500
		}

		// Calculate time range based on provided parameters
		// Priority: 1. Both provided 2. EndTime + Limit 3. StartTime + Limit 4. Current time + Limit
		if req.StartTime != nil && req.EndTime != nil {
			// Both times provided - use directly
			startTime = req.StartTime.Unix()
			endTime = req.EndTime.Unix()
		} else if req.EndTime != nil {
			// Only end time provided - calculate start time from end time and limit
			endTime = req.EndTime.Unix()
			startTime = endTime - int64(limit*step*60) // step is in minutes
		} else if req.StartTime != nil {
			// Only start time provided - calculate end time from start time and limit
			startTime = req.StartTime.Unix()
			endTime = startTime + int64(limit*step*60)
		} else {
			// No time provided - use current time as end, calculate start from limit
			endTime = time.Now().Unix()
			startTime = endTime - int64(limit*step*60)
		}

		contractReq := contractreq.GetContractKlineRequest{
			Symbol:    req.Symbol,
			Step:      int64(step),
			StartTime: startTime,
			EndTime:   endTime,
		}

		resp, err := a.client.Contract.GetKline(contractReq)
		if err != nil {
			return nil, err
		}

		// Convert contract kline data to common Candle type
		// Note: BitMart returns candles in ascending order (oldest first), no need to reverse
		candles := make([]*commontypes.Candle, 0, len(resp.Data))
		for i := range resp.Data {
			candle := a.converter.ConvertContractKlineToCandle(&resp.Data[i], req.Symbol, req.Interval)
			if candle != nil {
				candles = append(candles, candle)
			}
		}

		return candles, nil

	case commontypes.AccountTypeSpot:
		// Query spot klines (existing implementation)
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
		// Note: BitMart returns candles in ascending order (oldest first), no need to reverse
		// BitMart format: [timestamp, open, high, low, close, volume, quote_volume]
		candles := make([]*commontypes.Candle, 0, len(resp.Data))
		for _, klineData := range resp.Data {
			candle := a.converter.ConvertKlineArrayToCandle(klineData, req.Symbol, req.Interval)
			if candle != nil {
				candles = append(candles, candle)
			}
		}

		return candles, nil

	default:
		return nil, commontypes.ErrNotSupported
	}
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
