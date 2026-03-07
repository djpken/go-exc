package bingx

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/djpken/go-exc/exchanges/bingx/ws"
	commontypes "github.com/djpken/go-exc/types"
)

// WebSocketAdapter adapts the BingX WebSocket client to the common interface
type WebSocketAdapter struct {
	client        *ws.ClientWs
	privateClient *ws.PrivateClientWs
	converter     *Converter

	tickerChannels map[string]chan *commontypes.TickerUpdate
	candleChannels map[string]map[string]chan *commontypes.CandleUpdate // interval->symbol->ch

	// Private channel fan-out targets (ACCOUNT_UPDATE carries both balance and position data)
	accountChannel  chan *commontypes.AccountUpdate
	positionChannel chan *commontypes.PositionUpdate
}

func NewWebSocketAdapter(client *ws.ClientWs, privateClient *ws.PrivateClientWs) *WebSocketAdapter {
	return &WebSocketAdapter{
		client:         client,
		privateClient:  privateClient,
		converter:      NewConverter(),
		tickerChannels: make(map[string]chan *commontypes.TickerUpdate),
		candleChannels: make(map[string]map[string]chan *commontypes.CandleUpdate),
	}
}

func (a *WebSocketAdapter) Connect() error {
	return a.client.Connect()
}

func (a *WebSocketAdapter) Close() error {
	return a.client.Close()
}

// ─── Tickers ─────────────────────────────────────────────────────────────────

// tickerMsg is the expected structure of a BingX ticker WebSocket push
type tickerMsg struct {
	DataType string `json:"dataType"`
	Data     struct {
		E string `json:"e"`
		S string `json:"s"`
		P string `json:"p"` // price change
		C string `json:"c"` // last price
		H string `json:"h"` // high
		L string `json:"l"` // low
		V string `json:"v"` // volume
		O string `json:"o"` // open price
		B string `json:"B"` // best bid price
		A string `json:"A"` // best ask price
		// E_ shadows the outer "E" event time field
		E_ int64 `json:"E"` // event time
	} `json:"data"`
}

func (a *WebSocketAdapter) SubscribeTickers(userCh chan *commontypes.TickerUpdate, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("bingx: no symbols specified")
	}
	if !a.client.IsConnected() {
		if err := a.Connect(); err != nil {
			return fmt.Errorf("bingx: ws connect: %w", err)
		}
	}

	for _, symbol := range symbols {
		dataType := symbol + "@ticker"
		a.tickerChannels[symbol] = userCh

		sym := symbol // capture
		a.client.RegisterHandler(dataType, func(data []byte) {
			var msg tickerMsg
			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}
			conv := a.converter
			update := &commontypes.TickerUpdate{
				Symbol:    msg.Data.S,
				LastPrice: conv.str(msg.Data.C),
				High24h:   conv.str(msg.Data.H),
				Low24h:    conv.str(msg.Data.L),
				Volume24h: conv.str(msg.Data.V),
				BidPrice:  conv.str(msg.Data.B),
				AskPrice:  conv.str(msg.Data.A),
				Timestamp: commontypes.Timestamp(time.UnixMilli(msg.Data.E_)),
				Extra: map[string]interface{}{
					"priceChange": msg.Data.P,
					"openPrice":   msg.Data.O,
				},
			}
			select {
			case userCh <- update:
			default:
			}
		})

		if err := a.client.Subscribe(dataType); err != nil {
			return fmt.Errorf("bingx: subscribe ticker %s: %w", sym, err)
		}
	}
	return nil
}

func (a *WebSocketAdapter) UnsubscribeTickers(symbols ...string) error {
	for _, symbol := range symbols {
		dataType := symbol + "@ticker"
		a.client.UnregisterHandler(dataType)
		if err := a.client.Unsubscribe(dataType); err != nil {
			return err
		}
		delete(a.tickerChannels, symbol)
	}
	return nil
}

// ─── Candles ─────────────────────────────────────────────────────────────────

// klineMsg is the expected structure of a BingX kline WebSocket push
type klineMsg struct {
	DataType string `json:"dataType"`
	Data     struct {
		E string `json:"e"`
		S string `json:"s"`
		K struct {
			T int64  `json:"t"` // kline start time
			O string `json:"o"`
			H string `json:"h"`
			L string `json:"l"`
			C string `json:"c"`
			Q string `json:"q"` // volume
			I string `json:"i"` // interval
		} `json:"K"`
	} `json:"data"`
}

func (a *WebSocketAdapter) SubscribeCandles(userCh chan *commontypes.CandleUpdate, interval string, symbols ...string) error {
	if len(symbols) == 0 {
		return fmt.Errorf("bingx: no symbols specified")
	}
	wsInterval, err := a.converter.ConvertIntervalToWS(interval)
	if err != nil {
		return err
	}
	if !a.client.IsConnected() {
		if err := a.Connect(); err != nil {
			return fmt.Errorf("bingx: ws connect: %w", err)
		}
	}

	if a.candleChannels[interval] == nil {
		a.candleChannels[interval] = make(map[string]chan *commontypes.CandleUpdate)
	}

	for _, symbol := range symbols {
		dataType := fmt.Sprintf("%s@kline_%s", symbol, wsInterval)
		a.candleChannels[interval][symbol] = userCh

		sym := symbol
		iv := interval
		conv := a.converter
		a.client.RegisterHandler(dataType, func(data []byte) {
			var msg klineMsg
			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}
			k := msg.Data.K
			update := &commontypes.CandleUpdate{
				Symbol:    sym,
				Interval:  iv,
				Open:      conv.str(k.O),
				High:      conv.str(k.H),
				Low:       conv.str(k.L),
				Close:     conv.str(k.C),
				Volume:    conv.str(k.Q),
				Timestamp: commontypes.Timestamp(time.UnixMilli(k.T)),
				Confirmed: false, // BingX pushes forming candles; treat as unconfirmed
			}
			select {
			case userCh <- update:
			default:
			}
		})

		if err := a.client.Subscribe(dataType); err != nil {
			return fmt.Errorf("bingx: subscribe candle %s: %w", sym, err)
		}
	}
	return nil
}

func (a *WebSocketAdapter) UnsubscribeCandles(interval string, symbols ...string) error {
	wsInterval, err := a.converter.ConvertIntervalToWS(interval)
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		dataType := fmt.Sprintf("%s@kline_%s", symbol, wsInterval)
		a.client.UnregisterHandler(dataType)
		if err := a.client.Unsubscribe(dataType); err != nil {
			return err
		}
		if ch, ok := a.candleChannels[interval]; ok {
			delete(ch, symbol)
		}
	}
	return nil
}

// ─── Private channels ────────────────────────────────────────────────────────

// accountUpdateMsg is the structure of a BingX ACCOUNT_UPDATE private push.
// BingX sends both balance changes (B) and position changes (P) in the same event.
type accountUpdateMsg struct {
	EventType string `json:"e"` // "ACCOUNT_UPDATE"
	EventTime int64  `json:"E"`
	Account   struct {
		Reason   string `json:"m"`
		Balances []struct {
			Asset  string `json:"a"`  // asset name
			Wallet string `json:"wb"` // wallet balance
			Cross  string `json:"cw"` // cross wallet balance
			Delta  string `json:"bc"` // balance change amount
		} `json:"B"`
		Positions []struct {
			Symbol        string `json:"s"`  // symbol
			Amount        string `json:"pa"` // position amount
			EntryPrice    string `json:"ep"` // entry price
			RealizedPnL   string `json:"cr"` // accumulated realized PNL
			UnrealizedPnL string `json:"up"` // unrealized PNL
			MarginType    string `json:"mt"` // margin type (isolated/crossed)
			IsoWallet     string `json:"iw"` // isolated wallet (isolated mode only)
			PosSide       string `json:"ps"` // LONG / SHORT / BOTH
		} `json:"P"`
	} `json:"a"`
}

// registerAccountUpdateHandler (re-)registers a single ACCOUNT_UPDATE handler that fans out
// to accountChannel and positionChannel. It must be called whenever either channel changes.
func (a *WebSocketAdapter) registerAccountUpdateHandler() error {
	conv := a.converter
	return a.privateClient.RegisterHandler("ACCOUNT_UPDATE", func(data []byte) {
		var msg accountUpdateMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}

		// --- account balance fan-out ---
		if ch := a.accountChannel; ch != nil {
			balances := make([]*commontypes.Balance, 0, len(msg.Account.Balances))
			for _, b := range msg.Account.Balances {
				balances = append(balances, &commontypes.Balance{
					Currency:  b.Asset,
					Total:     conv.str(b.Wallet),
					Available: conv.str(b.Cross),
					Extra: map[string]interface{}{
						"balanceChange": b.Delta,
					},
				})
			}
			update := &commontypes.AccountUpdate{
				Balances:  balances,
				EventType: "update",
				UpdatedAt: commontypes.Timestamp(time.UnixMilli(msg.EventTime)),
				Extra: map[string]interface{}{
					"reason": msg.Account.Reason,
				},
			}
			select {
			case ch <- update:
			default:
			}
		}

		// --- position fan-out ---
		if ch := a.positionChannel; ch != nil && len(msg.Account.Positions) > 0 {
			positions := make([]*commontypes.Position, 0, len(msg.Account.Positions))
			for _, p := range msg.Account.Positions {
				var posSide commontypes.PositionSide
				switch p.PosSide {
				case "LONG":
					posSide = commontypes.PositionSideLong
				case "SHORT":
					posSide = commontypes.PositionSideShort
				default:
					posSide = commontypes.PositionSideNet
				}
				var marginMode commontypes.MarginMode
				switch p.MarginType {
				case "isolated":
					marginMode = commontypes.MarginModeIsolated
				default:
					marginMode = commontypes.MarginModeCross
				}
				positions = append(positions, &commontypes.Position{
					Symbol:        p.Symbol,
					PosSide:       posSide,
					Quantity:      conv.str(p.Amount),
					AvgPrice:      conv.str(p.EntryPrice),
					UnrealizedPnL: conv.str(p.UnrealizedPnL),
					RealizedPnL:   conv.str(p.RealizedPnL),
					MarginMode:    marginMode,
					Extra: map[string]interface{}{
						"isoWallet": p.IsoWallet,
					},
				})
			}
			update := &commontypes.PositionUpdate{
				Positions: positions,
				EventType: "update",
				UpdatedAt: commontypes.Timestamp(time.UnixMilli(msg.EventTime)),
			}
			select {
			case ch <- update:
			default:
			}
		}
	})
}

// orderTradeUpdateMsg is the structure of a BingX ORDER_TRADE_UPDATE private push.
type orderTradeUpdateMsg struct {
	EventType string `json:"e"` // "ORDER_TRADE_UPDATE"
	EventTime int64  `json:"E"`
	Order     struct {
		Symbol        string `json:"s"`
		ClientOrderID string `json:"c"`
		OrderID       int64  `json:"i"`
		Side          string `json:"S"` // BUY / SELL
		Type          string `json:"o"` // MARKET / LIMIT / ...
		Quantity      string `json:"q"`
		Price         string `json:"p"`
		AvgPrice      string `json:"ap"`
		ExecType      string `json:"x"` // execution type (NEW, TRADE, CANCELED, ...)
		Status        string `json:"X"` // order status
		FilledQty     string `json:"l"` // last filled qty
		TotalFilled   string `json:"z"` // cumulative filled qty
		Fee           string `json:"n"`
		FeeAsset      string `json:"N"`
		TradeTime     int64  `json:"T"`
		PosSide       string `json:"ps"` // LONG / SHORT / BOTH
		RealizedPnL   string `json:"rp"`
	} `json:"o"`
}

// SubscribeAccount registers for ACCOUNT_UPDATE balance events.
// The currencies parameter is ignored — BingX pushes all assets in a single event.
func (a *WebSocketAdapter) SubscribeAccount(userCh chan *commontypes.AccountUpdate, _ ...string) error {
	if a.privateClient == nil {
		return commontypes.ErrNotSupported
	}
	a.accountChannel = userCh
	return a.registerAccountUpdateHandler()
}

// UnsubscribeAccount removes the account balance listener.
// If a position listener is still active, the underlying handler is re-registered without account fan-out.
func (a *WebSocketAdapter) UnsubscribeAccount(_ ...string) error {
	if a.privateClient == nil {
		return nil
	}
	a.accountChannel = nil
	if a.positionChannel == nil {
		a.privateClient.UnregisterHandler("ACCOUNT_UPDATE")
		return nil
	}
	return a.registerAccountUpdateHandler()
}

// SubscribePosition registers for position updates carried inside ACCOUNT_UPDATE events.
// BingX sends position changes in the same event as balance changes (a.P array).
func (a *WebSocketAdapter) SubscribePosition(userCh chan *commontypes.PositionUpdate, _ commontypes.WebSocketSubscribeRequest) error {
	if a.privateClient == nil {
		return commontypes.ErrNotSupported
	}
	a.positionChannel = userCh
	return a.registerAccountUpdateHandler()
}

// UnsubscribePosition removes the position listener.
// If an account listener is still active, the underlying handler is re-registered without position fan-out.
func (a *WebSocketAdapter) UnsubscribePosition(_ commontypes.WebSocketSubscribeRequest) error {
	if a.privateClient == nil {
		return nil
	}
	a.positionChannel = nil
	if a.accountChannel == nil {
		a.privateClient.UnregisterHandler("ACCOUNT_UPDATE")
		return nil
	}
	return a.registerAccountUpdateHandler()
}

// SubscribeOrders registers a handler for ORDER_TRADE_UPDATE events.
func (a *WebSocketAdapter) SubscribeOrders(userCh chan *commontypes.OrderUpdate, _ commontypes.WebSocketSubscribeRequest) error {
	if a.privateClient == nil {
		return commontypes.ErrNotSupported
	}
	conv := a.converter
	return a.privateClient.RegisterHandler("ORDER_TRADE_UPDATE", func(data []byte) {
		var msg orderTradeUpdateMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}
		o := msg.Order
		order := &commontypes.Order{
			ID:             fmt.Sprintf("%d", o.OrderID),
			Symbol:         o.Symbol,
			ClientOrderID:  o.ClientOrderID,
			Side:           conv.ConvertOrderSide(o.Side),
			Type:           conv.ConvertOrderType(o.Type),
			Status:         conv.ConvertOrderStatus(o.Status),
			Price:          conv.str(o.Price),
			Quantity:       conv.str(o.Quantity),
			FilledQuantity: conv.str(o.TotalFilled),
			Fee:            conv.str(o.Fee),
			FeeCurrency:    o.FeeAsset,
			UpdatedAt:      commontypes.Timestamp(time.UnixMilli(o.TradeTime)),
			Extra: map[string]interface{}{
				"avgPrice":    o.AvgPrice,
				"execType":    o.ExecType,
				"posSide":     o.PosSide,
				"realizedPnL": o.RealizedPnL,
			},
		}
		update := &commontypes.OrderUpdate{
			Orders:    []*commontypes.Order{order},
			UpdatedAt: commontypes.Timestamp(time.UnixMilli(msg.EventTime)),
		}
		select {
		case userCh <- update:
		default:
		}
	})
}

// UnsubscribeOrders removes the ORDER_TRADE_UPDATE handler.
func (a *WebSocketAdapter) UnsubscribeOrders(_ commontypes.WebSocketSubscribeRequest) error {
	if a.privateClient == nil {
		return nil
	}
	a.privateClient.UnregisterHandler("ORDER_TRADE_UPDATE")
	return nil
}

// ─── Private channels (not supported via unified interface for BingX) ─────────

func (a *WebSocketAdapter) SetChannels(
	errCh chan *commontypes.WebSocketError,
	subCh chan *commontypes.WebSocketSubscribe,
	unsubCh chan *commontypes.WebSocketUnsubscribe,
	loginCh chan *commontypes.WebSocketLogin,
	successCh chan *commontypes.WebSocketSuccess,
	systemMsgCh chan *commontypes.WebSocketSystemMessage,
	systemErrCh chan *commontypes.WebSocketSystemError,
) error {
	return commontypes.ErrNotSupported
}
