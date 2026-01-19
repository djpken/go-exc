package ws

import (
	"encoding/json"
	"fmt"

	"github.com/djpken/go-exc/exchanges/bitmart/events/public"
)

// Public provides access to BitMart public WebSocket channels
type Public struct {
	*ClientWs
	tickerCh chan *public.TickerEvent
	depthCh  chan *public.DepthEvent
	tradeCh  chan *public.TradeEvent
	klineCh  chan *public.KlineEvent
}

// NewPublic creates a new Public instance
func NewPublic(c *ClientWs) *Public {
	return &Public{ClientWs: c}
}

// SubscribeTicker subscribes to ticker channel
//
// Channel: spot/ticker:{symbol}
func (p *Public) SubscribeTicker(symbol string, ch ...chan *public.TickerEvent) error {
	if len(ch) > 0 {
		p.tickerCh = ch[0]
	} else {
		p.tickerCh = make(chan *public.TickerEvent, 100)
	}

	channel := fmt.Sprintf("spot/ticker:%s", symbol)

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event public.TickerEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal ticker event: %v\n", err)
			return
		}
		select {
		case p.tickerCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeTicker unsubscribes from ticker channel
func (p *Public) UnsubscribeTicker(symbol string) error {
	channel := fmt.Sprintf("spot/ticker:%s", symbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeDepth subscribes to order book depth channel
//
// Channel: spot/depth{depth}:{symbol}
// Depth: 5, 20, 50
func (p *Public) SubscribeDepth(symbol string, depth int, ch ...chan *public.DepthEvent) error {
	if len(ch) > 0 {
		p.depthCh = ch[0]
	} else {
		p.depthCh = make(chan *public.DepthEvent, 100)
	}

	channel := fmt.Sprintf("spot/depth%d:%s", depth, symbol)

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event public.DepthEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal depth event: %v\n", err)
			return
		}
		select {
		case p.depthCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeDepth unsubscribes from depth channel
func (p *Public) UnsubscribeDepth(symbol string, depth int) error {
	channel := fmt.Sprintf("spot/depth%d:%s", depth, symbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeTrade subscribes to trade channel
//
// Channel: spot/trade:{symbol}
func (p *Public) SubscribeTrade(symbol string, ch ...chan *public.TradeEvent) error {
	if len(ch) > 0 {
		p.tradeCh = ch[0]
	} else {
		p.tradeCh = make(chan *public.TradeEvent, 100)
	}

	channel := fmt.Sprintf("spot/trade:%s", symbol)

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event public.TradeEvent
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
func (p *Public) UnsubscribeTrade(symbol string) error {
	channel := fmt.Sprintf("spot/trade:%s", symbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeKline subscribes to kline/candlestick channel
//
// Channel: spot/kline{step}:{symbol}
// Step: 1m, 3m, 5m, 15m, 30m, 45m, 1H, 2H, 3H, 4H, 1D, 1W, 1M
func (p *Public) SubscribeKline(symbol string, step string, ch ...chan *public.KlineEvent) error {
	if len(ch) > 0 {
		p.klineCh = ch[0]
	} else {
		p.klineCh = make(chan *public.KlineEvent, 100)
	}

	channel := fmt.Sprintf("spot/kline%s:%s", step, symbol)

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event public.KlineEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal kline event: %v\n", err)
			return
		}
		select {
		case p.klineCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeKline unsubscribes from kline channel
func (p *Public) UnsubscribeKline(symbol string, step string) error {
	channel := fmt.Sprintf("spot/kline%s:%s", step, symbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// GetTickerChan returns the ticker channel
func (p *Public) GetTickerChan() chan *public.TickerEvent {
	return p.tickerCh
}

// GetDepthChan returns the depth channel
func (p *Public) GetDepthChan() chan *public.DepthEvent {
	return p.depthCh
}

// GetTradeChan returns the trade channel
func (p *Public) GetTradeChan() chan *public.TradeEvent {
	return p.tradeCh
}

// GetKlineChan returns the kline channel
func (p *Public) GetKlineChan() chan *public.KlineEvent {
	return p.klineCh
}
