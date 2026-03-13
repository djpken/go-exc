package bingx

import (
	"context"
	"fmt"
	"strconv"

	"github.com/djpken/go-exc/exchanges/bingx/rest"
	commontypes "github.com/djpken/go-exc/types"
)

// RESTAdapter exposes grouped REST operations via the common interface
type RESTAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

func NewRESTAdapter(client *rest.ClientRest) *RESTAdapter {
	return &RESTAdapter{client: client, converter: NewConverter()}
}

func (a *RESTAdapter) Market() *MarketAPIAdapter {
	return &MarketAPIAdapter{client: a.client, converter: a.converter}
}

func (a *RESTAdapter) Account() *AccountAPIAdapter {
	return &AccountAPIAdapter{client: a.client, converter: a.converter}
}

func (a *RESTAdapter) Trade() *TradeAPIAdapter {
	return &TradeAPIAdapter{client: a.client, converter: a.converter}
}

// ─── Market ──────────────────────────────────────────────────────────────────

type MarketAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

func (a *MarketAPIAdapter) GetTicker(_ context.Context, symbol string) (*commontypes.Ticker, error) {
	resp, err := a.client.Market.GetTicker(symbol)
	if err != nil {
		return nil, err
	}
	return a.converter.ConvertTicker(&resp.Data), nil
}

func (a *MarketAPIAdapter) GetTickers(_ context.Context) ([]*commontypes.Ticker, error) {
	resp, err := a.client.Market.GetTickers()
	if err != nil {
		return nil, err
	}
	tickers := make([]*commontypes.Ticker, 0, len(resp.Data))
	for i := range resp.Data {
		tickers = append(tickers, a.converter.ConvertTicker(&resp.Data[i]))
	}
	return tickers, nil
}

func (a *MarketAPIAdapter) GetOrderBook(_ context.Context, symbol string, depth int) (*commontypes.OrderBook, error) {
	resp, err := a.client.Market.GetOrderBook(symbol, depth)
	if err != nil {
		return nil, err
	}
	return a.converter.ConvertOrderBook(&resp.Data, symbol), nil
}

func (a *MarketAPIAdapter) GetCandles(_ context.Context, req commontypes.GetCandlesRequest) ([]*commontypes.Candle, error) {
	interval, err := a.converter.ConvertIntervalToREST(req.Interval)
	if err != nil {
		return nil, err
	}

	var startMs, endMs int64
	if req.StartTime != nil {
		startMs = req.StartTime.UnixMilli()
	}
	if req.EndTime != nil {
		endMs = req.EndTime.UnixMilli()
	}

	resp, err := a.client.Market.GetKlines(req.Symbol, interval, startMs, endMs, req.Limit, 0)
	if err != nil {
		return nil, err
	}

	candles := make([]*commontypes.Candle, 0, len(resp.Data))
	for i := range resp.Data {
		candles = append(candles, a.converter.ConvertKline(&resp.Data[i], req.Symbol, req.Interval))
	}
	return candles, nil
}

func (a *MarketAPIAdapter) GetInstruments(_ context.Context) ([]*commontypes.Instrument, error) {
	resp, err := a.client.Market.GetContracts()
	if err != nil {
		return nil, err
	}
	instruments := make([]*commontypes.Instrument, 0, len(resp.Data))
	for i := range resp.Data {
		instruments = append(instruments, a.converter.ConvertInstrument(&resp.Data[i]))
	}
	return instruments, nil
}

// ─── Account ─────────────────────────────────────────────────────────────────

type AccountAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

func (a *AccountAPIAdapter) GetBalance(_ context.Context, currencies ...string) (*commontypes.AccountBalance, error) {
	currency := ""
	if len(currencies) > 0 {
		currency = currencies[0]
	}
	resp, err := a.client.Account.GetBalance(currency)
	if err != nil {
		return nil, err
	}
	return a.converter.ConvertBalance(&resp.Data), nil
}

func (a *AccountAPIAdapter) GetPositions(_ context.Context, symbols ...string) ([]*commontypes.Position, error) {
	symbol := ""
	if len(symbols) == 1 {
		symbol = symbols[0]
	}
	resp, err := a.client.Account.GetPositions(symbol)
	if err != nil {
		return nil, err
	}
	positions := make([]*commontypes.Position, 0, len(resp.Data))
	for i := range resp.Data {
		positions = append(positions, a.converter.ConvertPosition(&resp.Data[i]))
	}
	return positions, nil
}

func (a *AccountAPIAdapter) GetLeverage(_ context.Context, symbols []string) ([]*commontypes.Leverage, error) {
	var all []*commontypes.Leverage
	for _, sym := range symbols {
		resp, err := a.client.Account.GetLeverage(sym)
		if err != nil {
			return nil, fmt.Errorf("bingx: get leverage for %s: %w", sym, err)
		}
		all = append(all, a.converter.ConvertLeverage(sym, &resp.Data)...)
	}
	return all, nil
}

func (a *AccountAPIAdapter) SetLeverage(_ context.Context, req commontypes.SetLeverageRequest) (*commontypes.Leverage, error) {
	side := "LONG"
	switch req.PosSide {
	case commontypes.PositionSideShort:
		side = "SHORT"
	case commontypes.PositionSideNet:
		// one-way mode: set both sides
		side = "LONG"
	}

	resp, err := a.client.Account.SetLeverage(req.Symbol, side, req.Leverage)
	if err != nil {
		return nil, err
	}

	var posSide commontypes.PositionSide
	if resp.Data.Side == "SHORT" {
		posSide = commontypes.PositionSideShort
	} else {
		posSide = commontypes.PositionSideLong
	}

	return &commontypes.Leverage{
		Symbol:   resp.Data.Symbol,
		Leverage: resp.Data.Leverage,
		PosSide:  posSide,
	}, nil
}

// ─── Trade ───────────────────────────────────────────────────────────────────

type TradeAPIAdapter struct {
	client    *rest.ClientRest
	converter *Converter
}

func (a *TradeAPIAdapter) PlaceOrder(_ context.Context, req commontypes.PlaceOrderRequest) (*commontypes.Order, error) {
	side := string(req.Side)
	positionSide := ""
	switch req.PosSide {
	case commontypes.PositionSideLong:
		positionSide = "LONG"
	case commontypes.PositionSideShort:
		positionSide = "SHORT"
	case commontypes.PositionSideNet:
		positionSide = "BOTH"
	}

	resp, err := a.client.Trade.PlaceOrder(
		req.Symbol, side, positionSide, req.Type,
		req.Price, req.Quantity,
		req.ClientOrderID,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return a.converter.ConvertOrder(&resp.Data.Order), nil
}

// PlaceSingleOrder places exactly one order and wraps the result in PlaceOrderResult.
func (a *TradeAPIAdapter) PlaceSingleOrder(ctx context.Context, req commontypes.PlaceOrderRequest) (*commontypes.PlaceOrderResult, error) {
	order, err := a.PlaceOrder(ctx, req)
	if err != nil {
		return &commontypes.PlaceOrderResult{Error: err}, nil
	}
	return &commontypes.PlaceOrderResult{Order: order}, nil
}

// PlaceMultiOrder places multiple orders sequentially and returns per-order results.
// BingX has no native batch-order API, so orders are sent one by one.
func (a *TradeAPIAdapter) PlaceMultiOrder(ctx context.Context, reqs []commontypes.PlaceOrderRequest) ([]*commontypes.PlaceOrderResult, error) {
	results := make([]*commontypes.PlaceOrderResult, len(reqs))
	for i, req := range reqs {
		order, err := a.PlaceOrder(ctx, req)
		results[i] = &commontypes.PlaceOrderResult{Order: order, Error: err}
	}
	return results, nil
}

func (a *TradeAPIAdapter) CancelOrder(_ context.Context, symbol, orderID string, _ map[string]interface{}) error {
	var oid int64
	if orderID != "" {
		var err error
		oid, err = strconv.ParseInt(orderID, 10, 64)
		if err != nil {
			return fmt.Errorf("bingx: invalid orderId %q: %w", orderID, err)
		}
	}
	_, err := a.client.Trade.CancelOrder(symbol, oid, "")
	return err
}

func (a *TradeAPIAdapter) GetOrderDetail(_ context.Context, req commontypes.GetOrderRequest) (*commontypes.Order, error) {
	var oid int64
	if req.OrderID != "" {
		var err error
		oid, err = strconv.ParseInt(req.OrderID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bingx: invalid orderId %q: %w", req.OrderID, err)
		}
	}
	resp, err := a.client.Trade.GetOrder(req.Symbol, oid, req.ClientOrderID)
	if err != nil {
		return nil, err
	}
	return a.converter.ConvertOrder(&resp.Data), nil
}
