package ws

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/djpken/go-exc/exchanges/bitmart/events/private"
)

// Private provides access to BitMart private WebSocket channels
type Private struct {
	*ClientWs
	orderCh          chan *private.OrderEvent
	balanceCh        chan *private.BalanceEvent
	tradeCh          chan *private.TradeEvent
	futuresAssetCh   chan *private.FuturesAssetEvent
	futuresPositionCh chan *private.FuturesPositionEvent
}

// NewPrivate creates a new Private instance
func NewPrivate(c *ClientWs) *Private {
	return &Private{ClientWs: c}
}

// SubscribeOrder subscribes to order update channel
//
// Channel: spot/user/order
// Requires authentication
func (p *Private) SubscribeOrder(ch ...chan *private.OrderEvent) error {
	if !p.IsAuthenticated() {
		return errors.New("not authenticated, please call Login() first")
	}

	if len(ch) > 0 {
		p.orderCh = ch[0]
	} else {
		p.orderCh = make(chan *private.OrderEvent, 100)
	}

	channel := "spot/user/order"

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event private.OrderEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal order event: %v\n", err)
			return
		}
		select {
		case p.orderCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeOrder unsubscribes from order channel
func (p *Private) UnsubscribeOrder() error {
	channel := "spot/user/order"
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeBalance subscribes to balance update channel
//
// Channel: spot/user/balance
// Requires authentication
func (p *Private) SubscribeBalance(ch ...chan *private.BalanceEvent) error {
	if !p.IsAuthenticated() {
		return errors.New("not authenticated, please call Login() first")
	}

	if len(ch) > 0 {
		p.balanceCh = ch[0]
	} else {
		p.balanceCh = make(chan *private.BalanceEvent, 100)
	}

	channel := "spot/user/balance"

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event private.BalanceEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal balance event: %v\n", err)
			return
		}
		select {
		case p.balanceCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeBalance unsubscribes from balance channel
func (p *Private) UnsubscribeBalance() error {
	channel := "spot/user/balance"
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeTrade subscribes to trade execution channel
//
// Channel: spot/user/trade
// Requires authentication
func (p *Private) SubscribeTrade(ch ...chan *private.TradeEvent) error {
	if !p.IsAuthenticated() {
		return errors.New("not authenticated, please call Login() first")
	}

	if len(ch) > 0 {
		p.tradeCh = ch[0]
	} else {
		p.tradeCh = make(chan *private.TradeEvent, 100)
	}

	channel := "spot/user/trade"

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event private.TradeEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal trade event: %v\n", err)
			return
		}
		select {
		case p.tradeCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeTrade unsubscribes from trade channel
func (p *Private) UnsubscribeTrade() error {
	channel := "spot/user/trade"
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// GetOrderChan returns the order channel
func (p *Private) GetOrderChan() chan *private.OrderEvent {
	return p.orderCh
}

// GetBalanceChan returns the balance channel
func (p *Private) GetBalanceChan() chan *private.BalanceEvent {
	return p.balanceCh
}

// GetTradeChan returns the trade channel
func (p *Private) GetTradeChan() chan *private.TradeEvent {
	return p.tradeCh
}

// SubscribeFuturesAsset subscribes to futures asset balance update channel
//
// Channel: futures/asset:CURRENCY (e.g., futures/asset:USDT, futures/asset:BTC, futures/asset:ETH)
// Requires authentication
//
// Example:
//   err := client.Private.SubscribeFuturesAsset(assetCh, "USDT", "BTC")
func (p *Private) SubscribeFuturesAsset(ch chan *private.FuturesAssetEvent, currencies ...string) error {
	if !p.IsAuthenticated() {
		return errors.New("not authenticated, please call Login() first")
	}

	if len(currencies) == 0 {
		return errors.New("at least one currency must be specified")
	}

	p.futuresAssetCh = ch

	// Build channel list: futures/asset:CURRENCY
	channels := make([]string, len(currencies))
	for i, currency := range currencies {
		channels[i] = fmt.Sprintf("futures/asset:%s", currency)
	}

	// Register message handler for each channel
	for _, channel := range channels {
		p.RegisterHandler(channel, func(data []byte) {
			var event private.FuturesAssetEvent
			if err := json.Unmarshal(data, &event); err != nil {
				fmt.Printf("Failed to unmarshal futures asset event: %v\n", err)
				return
			}
			select {
			case p.futuresAssetCh <- &event:
			default:
				// Channel full, drop message
			}
		})
	}

	// Subscribe to all channels at once
	subscribeMsg := map[string]interface{}{
		"action": "subscribe",
		"args":   channels,
	}

	// fmt.Printf("[WS] Subscribing to futures asset channels: %v\n", channels)

	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conn.WriteJSON(subscribeMsg)
}

// UnsubscribeFuturesAsset unsubscribes from futures asset channels
func (p *Private) UnsubscribeFuturesAsset(currencies ...string) error {
	if len(currencies) == 0 {
		return errors.New("at least one currency must be specified")
	}

	// Build channel list
	channels := make([]string, len(currencies))
	for i, currency := range currencies {
		channels[i] = fmt.Sprintf("futures/asset:%s", currency)
	}

	// Unregister handlers
	for _, channel := range channels {
		p.UnregisterHandler(channel)
	}

	// Unsubscribe from all channels at once
	unsubscribeMsg := map[string]interface{}{
		"action": "unsubscribe",
		"args":   channels,
	}

	// fmt.Printf("[WS] Unsubscribing from futures asset channels: %v\n", channels)

	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conn.WriteJSON(unsubscribeMsg)
}

// SubscribeFuturesPosition subscribes to futures position update channel
//
// Channel: futures/position
// Requires authentication
//
// Example:
//   err := client.Private.SubscribeFuturesPosition(positionCh)
func (p *Private) SubscribeFuturesPosition(ch chan *private.FuturesPositionEvent) error {
	if !p.IsAuthenticated() {
		return errors.New("not authenticated, please call Login() first")
	}

	p.futuresPositionCh = ch

	channel := "futures/position"

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event private.FuturesPositionEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal futures position event: %v\n", err)
			return
		}
		select {
		case p.futuresPositionCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeFuturesPosition unsubscribes from futures position channel
func (p *Private) UnsubscribeFuturesPosition() error {
	channel := "futures/position"
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// GetFuturesAssetChan returns the futures asset channel
func (p *Private) GetFuturesAssetChan() chan *private.FuturesAssetEvent {
	return p.futuresAssetCh
}

// GetFuturesPositionChan returns the futures position channel
func (p *Private) GetFuturesPositionChan() chan *private.FuturesPositionEvent {
	return p.futuresPositionCh
}
