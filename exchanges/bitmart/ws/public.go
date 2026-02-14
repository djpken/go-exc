package ws

import (
	"encoding/json"
	"fmt"

	"github.com/djpken/go-exc/exchanges/bitmart/events/public"
)

// Public provides access to BitMart public WebSocket channels
type Public struct {
	*ClientWs
	tickerCh        chan *public.TickerEvent
	futuresTickerCh chan *public.FuturesTickerEvent
	depthCh         chan *public.DepthEvent
	tradeCh         chan *public.TradeEvent
	klineCh         chan *public.KlineEvent
}

// NewPublic creates a new Public instance
func NewPublic(c *ClientWs) *Public {
	return &Public{ClientWs: c}
}

// SubscribeTicker subscribes to ticker channel (legacy spot ticker)
//
// Channel: spot/ticker:{symbol}
// Note: For futures ticker, use SubscribeFuturesTicker instead
func (p *Public) SubscribeTicker(symbol string, ch ...chan *public.TickerEvent) error {
	if len(ch) > 0 {
		p.tickerCh = ch[0]
	} else {
		p.tickerCh = make(chan *public.TickerEvent, 100)
	}

	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("spot/ticker:%s", normalizedSymbol)

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event public.TickerEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal spot ticker event: %v\n", err)
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

// SubscribeFuturesTicker subscribes to futures ticker channel
//
// Channel: futures/ticker:{symbol}
// Note: BitMart v2 API uses futures channels and symbol format without underscore (e.g., BTCUSDT)
func (p *Public) SubscribeFuturesTicker(symbol string, ch ...chan *public.FuturesTickerEvent) error {
	if len(ch) > 0 {
		p.futuresTickerCh = ch[0]
	} else {
		p.futuresTickerCh = make(chan *public.FuturesTickerEvent, 100)
	}

	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/ticker:%s", normalizedSymbol)

	// Register message handler
	p.RegisterHandler(channel, func(data []byte) {
		var event public.FuturesTickerEvent
		if err := json.Unmarshal(data, &event); err != nil {
			fmt.Printf("Failed to unmarshal futures ticker event: %v\n", err)
			fmt.Printf("Raw data: %s\n", string(data))
			return
		}
		select {
		case p.futuresTickerCh <- &event:
		default:
			// Channel full, drop message
		}
	})

	return p.Subscribe(channel)
}

// UnsubscribeTicker unsubscribes from ticker channel (legacy spot ticker)
func (p *Public) UnsubscribeTicker(symbol string) error {
	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("spot/ticker:%s", normalizedSymbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// UnsubscribeFuturesTicker unsubscribes from futures ticker channel
func (p *Public) UnsubscribeFuturesTicker(symbol string) error {
	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/ticker:%s", normalizedSymbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeDepth subscribes to order book depth channel
//
// Channel: futures/depth{depth}:{symbol}
// Depth: 5, 20, 50
// Note: BitMart v2 API uses futures channels and symbol format without underscore (e.g., BTCUSDT)
func (p *Public) SubscribeDepth(symbol string, depth int, ch ...chan *public.DepthEvent) error {
	if len(ch) > 0 {
		p.depthCh = ch[0]
	} else {
		p.depthCh = make(chan *public.DepthEvent, 100)
	}

	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/depth%d:%s", depth, normalizedSymbol)

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
	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/depth%d:%s", depth, normalizedSymbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeTrade subscribes to trade channel
//
// Channel: futures/trade:{symbol}
// Note: BitMart v2 API uses futures channels and symbol format without underscore (e.g., BTCUSDT)
func (p *Public) SubscribeTrade(symbol string, ch ...chan *public.TradeEvent) error {
	if len(ch) > 0 {
		p.tradeCh = ch[0]
	} else {
		p.tradeCh = make(chan *public.TradeEvent, 100)
	}

	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/trade:%s", normalizedSymbol)

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
	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/trade:%s", normalizedSymbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// SubscribeKline subscribes to kline/candlestick channel
//
// Channel: futures/kline{step}:{symbol}
// Step: 1m, 3m, 5m, 15m, 30m, 45m, 1H, 2H, 3H, 4H, 1D, 1W, 1M
// Note: BitMart v2 API uses futures channels and symbol format without underscore (e.g., BTCUSDT)
func (p *Public) SubscribeKline(symbol string, step string, ch ...chan *public.KlineEvent) error {
	if len(ch) > 0 {
		p.klineCh = ch[0]
	} else {
		p.klineCh = make(chan *public.KlineEvent, 100)
	}

	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/kline%s:%s", step, normalizedSymbol)

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
	// Convert symbol format: BTC_USDT -> BTCUSDT (remove underscore)
	normalizedSymbol := normalizeSymbol(symbol)
	channel := fmt.Sprintf("futures/kline%s:%s", step, normalizedSymbol)
	p.UnregisterHandler(channel)

	return p.Unsubscribe(channel)
}

// GetTickerChan returns the ticker channel (legacy spot ticker)
func (p *Public) GetTickerChan() chan *public.TickerEvent {
	return p.tickerCh
}

// GetFuturesTickerChan returns the futures ticker channel
func (p *Public) GetFuturesTickerChan() chan *public.FuturesTickerEvent {
	return p.futuresTickerCh
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

// normalizeSymbol converts symbol format from BTC_USDT to BTCUSDT (removes underscore)
// BitMart v2 API uses symbol format without underscore
func normalizeSymbol(symbol string) string {
	// Remove underscore from symbol
	result := ""
	for _, char := range symbol {
		if char != '_' {
			result += string(char)
		}
	}
	return result
}
