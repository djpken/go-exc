package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/djpken/go-exc/exchanges/bitmart/types"
	"github.com/djpken/go-exc/exchanges/bitmart/ws"
)

const apiKey = "f42ce865e6123a77b507659452384a2b48165991"
const secretKey = "72cca2bc92ee262f7e5398e3344b15638d43224a6fc5404567fc884f655a4b73"
const memo = ".Vm3djpcl3gj94" // BitMart requires a memo

func main() {
	// Example selection
	examples := map[int]func(){
		1: example1_TickerSubscription,
		2: example2_DepthSubscription,
		3: example3_TradeSubscription,
		4: example4_KlineSubscription,
		5: example5_MultiplePublicChannels,
		6: example6_PrivateOrderUpdates,
		7: example7_PrivateBalanceUpdates,
		8: example8_PrivateTradeUpdates,
		9: example9_AllPrivateChannels,
	}

	fmt.Println("BitMart WebSocket API Examples")
	fmt.Println("==============================")
	fmt.Println("1. Ticker Subscription (Public)")
	fmt.Println("2. Depth/Order Book Subscription (Public)")
	fmt.Println("3. Trade Subscription (Public)")
	fmt.Println("4. Kline/Candlestick Subscription (Public)")
	fmt.Println("5. Multiple Public Channels")
	fmt.Println("6. Private Order Updates (Requires API Key)")
	fmt.Println("7. Private Balance Updates (Requires API Key)")
	fmt.Println("8. Private Trade Updates (Requires API Key)")
	fmt.Println("9. All Private Channels (Requires API Key)")
	fmt.Print("\nSelect example (1-9): ")

	var choice int
	fmt.Scan(&choice)

	if fn, ok := examples[choice]; ok {
		fn()
	} else {
		fmt.Println("Invalid choice")
	}
}

// Example 1: Subscribe to ticker updates for a symbol
func example1_TickerSubscription() {
	fmt.Println("\n=== Example 1: Ticker Subscription ===")
	fmt.Println("Subscribing to BTC_USDT ticker updates...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create WebSocket client
	cfg := &ws.BitMartConfig{
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	// Connect to WebSocket
	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	fmt.Println("Connected to WebSocket server")

	// Subscribe to ticker
	if err := client.Public.SubscribeTicker("BTC_USDT"); err != nil {
		log.Fatal("Failed to subscribe to ticker:", err)
	}

	fmt.Println("Subscribed to BTC_USDT ticker")
	fmt.Println("Waiting for ticker updates (Press Ctrl+C to exit)...")

	// Handle ticker events
	tickerCh := client.Public.GetTickerChan()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case ticker := <-tickerCh:
			fmt.Printf("\n[Ticker Update]\n")
			fmt.Printf("  Symbol: %s\n", ticker.Symbol)
			fmt.Printf("  Last Price: %s\n", ticker.LastPrice)
			fmt.Printf("  24h Volume: %s\n", ticker.QuoteVolume)
			fmt.Printf("  24h High: %s\n", ticker.HighPrice)
			fmt.Printf("  24h Low: %s\n", ticker.LowPrice)
			fmt.Printf("  Timestamp: %d\n", ticker.Timestamp)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 2: Subscribe to order book depth
func example2_DepthSubscription() {
	fmt.Println("\n=== Example 2: Depth Subscription ===")
	fmt.Println("Subscribing to BTC_USDT order book (depth: 5)...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	// Subscribe to depth
	if err := client.Public.SubscribeDepth("BTC_USDT", 5); err != nil {
		log.Fatal("Failed to subscribe to depth:", err)
	}

	fmt.Println("Subscribed to BTC_USDT depth")
	fmt.Println("Waiting for depth updates (Press Ctrl+C to exit)...")

	depthCh := client.Public.GetDepthChan()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case depth := <-depthCh:
			fmt.Printf("\n[Depth Update]\n")
			fmt.Printf("  Symbol: %s\n", depth.Symbol)
			fmt.Printf("  Bids (Top 3):\n")
			for i := 0; i < 3 && i < len(depth.Bids); i++ {
				fmt.Printf("    [%d] Price: %s, Size: %s\n", i+1, depth.Bids[i][0], depth.Bids[i][1])
			}
			fmt.Printf("  Asks (Top 3):\n")
			for i := 0; i < 3 && i < len(depth.Asks); i++ {
				fmt.Printf("    [%d] Price: %s, Size: %s\n", i+1, depth.Asks[i][0], depth.Asks[i][1])
			}
			fmt.Printf("  Timestamp: %d\n", depth.Timestamp)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 3: Subscribe to trade updates
func example3_TradeSubscription() {
	fmt.Println("\n=== Example 3: Trade Subscription ===")
	fmt.Println("Subscribing to BTC_USDT trades...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	// Subscribe to trades
	if err := client.Public.SubscribeTrade("BTC_USDT"); err != nil {
		log.Fatal("Failed to subscribe to trades:", err)
	}

	fmt.Println("Subscribed to BTC_USDT trades")
	fmt.Println("Waiting for trade updates (Press Ctrl+C to exit)...")

	tradeCh := client.Public.GetTradeChan()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case trade := <-tradeCh:
			fmt.Printf("\n[Trade]\n")
			fmt.Printf("  Symbol: %s\n", trade.Symbol)
			fmt.Printf("  Price: %s\n", trade.Price)
			fmt.Printf("  Size: %s\n", trade.Size)
			fmt.Printf("  Side: %s\n", trade.Side)
			fmt.Printf("  Timestamp: %d\n", trade.Timestamp)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 4: Subscribe to kline/candlestick data
func example4_KlineSubscription() {
	fmt.Println("\n=== Example 4: Kline Subscription ===")
	fmt.Println("Subscribing to BTC_USDT 1m kline...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	// Subscribe to kline
	if err := client.Public.SubscribeKline("BTC_USDT", "1m"); err != nil {
		log.Fatal("Failed to subscribe to kline:", err)
	}

	fmt.Println("Subscribed to BTC_USDT 1m kline")
	fmt.Println("Waiting for kline updates (Press Ctrl+C to exit)...")

	klineCh := client.Public.GetKlineChan()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case kline := <-klineCh:
			fmt.Printf("\n[Kline]\n")
			fmt.Printf("  Symbol: %s\n", kline.Symbol)
			fmt.Printf("  Open: %s\n", kline.Open)
			fmt.Printf("  High: %s\n", kline.High)
			fmt.Printf("  Low: %s\n", kline.Low)
			fmt.Printf("  Close: %s\n", kline.Close)
			fmt.Printf("  Volume: %s\n", kline.Volume)
			fmt.Printf("  Timestamp: %d\n", kline.Timestamp)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 5: Subscribe to multiple public channels
func example5_MultiplePublicChannels() {
	fmt.Println("\n=== Example 5: Multiple Public Channels ===")
	fmt.Println("Subscribing to ticker, depth, and trades...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	// Subscribe to multiple channels
	if err := client.Public.SubscribeTicker("BTC_USDT"); err != nil {
		log.Fatal("Failed to subscribe to ticker:", err)
	}

	if err := client.Public.SubscribeDepth("BTC_USDT", 5); err != nil {
		log.Fatal("Failed to subscribe to depth:", err)
	}

	if err := client.Public.SubscribeTrade("BTC_USDT"); err != nil {
		log.Fatal("Failed to subscribe to trades:", err)
	}

	fmt.Println("Subscribed to multiple channels")
	fmt.Println("Waiting for updates (Press Ctrl+C to exit)...")

	tickerCh := client.Public.GetTickerChan()
	depthCh := client.Public.GetDepthChan()
	tradeCh := client.Public.GetTradeChan()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case ticker := <-tickerCh:
			fmt.Printf("\n[Ticker] %s: %s\n", ticker.Symbol, ticker.LastPrice)

		case depth := <-depthCh:
			fmt.Printf("\n[Depth] %s: Best Bid: %s, Best Ask: %s\n",
				depth.Symbol, depth.Bids[0][0], depth.Asks[0][0])

		case trade := <-tradeCh:
			fmt.Printf("\n[Trade] %s: %s @ %s (%s)\n",
				trade.Symbol, trade.Size, trade.Price, trade.Side)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 6: Subscribe to private order updates
func example6_PrivateOrderUpdates() {
	fmt.Println("\n=== Example 6: Private Order Updates ===")

	// Get credentials from environment
	apiKey := apiKey
	secretKey := secretKey
	memo := memo

	if apiKey == "" || secretKey == "" || memo == "" {
		log.Fatal("Please set BITMART_API_KEY, BITMART_SECRET_KEY, and BITMART_MEMO environment variables")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	fmt.Println("Connected, authenticating...")

	// Authenticate
	if err := client.Login(); err != nil {
		log.Fatal("Failed to login:", err)
	}

	// Wait a bit for authentication
	time.Sleep(1 * time.Second)

	// Subscribe to order updates
	if err := client.Private.SubscribeOrder(); err != nil {
		log.Fatal("Failed to subscribe to orders:", err)
	}

	fmt.Println("Subscribed to order updates")
	fmt.Println("Waiting for order updates (Press Ctrl+C to exit)...")

	orderCh := client.Private.GetOrderChan()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case order := <-orderCh:
			fmt.Printf("\n[Order Update]\n")
			fmt.Printf("  Order ID: %s\n", order.OrderID)
			fmt.Printf("  Symbol: %s\n", order.Symbol)
			fmt.Printf("  Side: %s\n", order.Side)
			fmt.Printf("  Type: %s\n", order.Type)
			fmt.Printf("  Status: %s\n", order.Status)
			fmt.Printf("  Price: %s\n", order.Price)
			fmt.Printf("  Size: %s\n", order.Size)
			fmt.Printf("  Filled: %s\n", order.FilledSize)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 7: Subscribe to private balance updates
func example7_PrivateBalanceUpdates() {
	fmt.Println("\n=== Example 7: Private Balance Updates ===")

	apiKey := apiKey
	secretKey := secretKey
	memo := memo

	if apiKey == "" || secretKey == "" || memo == "" {
		log.Fatal("Please set BITMART_API_KEY, BITMART_SECRET_KEY, and BITMART_MEMO environment variables")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	if err := client.Login(); err != nil {
		log.Fatal("Failed to login:", err)
	}

	time.Sleep(1 * time.Second)

	if err := client.Private.SubscribeBalance(); err != nil {
		log.Fatal("Failed to subscribe to balance:", err)
	}

	fmt.Println("Subscribed to balance updates")
	fmt.Println("Waiting for balance updates (Press Ctrl+C to exit)...")

	balanceCh := client.Private.GetBalanceChan()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case balance := <-balanceCh:
			fmt.Printf("\n[Balance Update]\n")
			fmt.Printf("  Currency: %s\n", balance.Currency)
			fmt.Printf("  Available: %s\n", balance.Available)
			fmt.Printf("  Frozen: %s\n", balance.Frozen)
			fmt.Printf("  Total: %s\n", balance.Total)
			fmt.Printf("  Timestamp: %d\n", balance.Timestamp)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 8: Subscribe to private trade updates
func example8_PrivateTradeUpdates() {
	fmt.Println("\n=== Example 8: Private Trade Updates ===")

	apiKey := apiKey
	secretKey := secretKey
	memo := memo

	if apiKey == "" || secretKey == "" || memo == "" {
		log.Fatal("Please set BITMART_API_KEY, BITMART_SECRET_KEY, and BITMART_MEMO environment variables")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	if err := client.Login(); err != nil {
		log.Fatal("Failed to login:", err)
	}

	time.Sleep(1 * time.Second)

	if err := client.Private.SubscribeTrade(); err != nil {
		log.Fatal("Failed to subscribe to trades:", err)
	}

	fmt.Println("Subscribed to trade updates")
	fmt.Println("Waiting for trade updates (Press Ctrl+C to exit)...")

	tradeCh := client.Private.GetTradeChan()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case trade := <-tradeCh:
			fmt.Printf("\n[Trade Execution]\n")
			fmt.Printf("  Trade ID: %s\n", trade.TradeID)
			fmt.Printf("  Order ID: %s\n", trade.OrderID)
			fmt.Printf("  Symbol: %s\n", trade.Symbol)
			fmt.Printf("  Side: %s\n", trade.Side)
			fmt.Printf("  Price: %s\n", trade.Price)
			fmt.Printf("  Size: %s\n", trade.Size)
			fmt.Printf("  Fee: %s %s\n", trade.Fee, trade.FeeCurrency)
			fmt.Printf("  Exec Time: %d\n", trade.ExecTime)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

// Example 9: Subscribe to all private channels
func example9_AllPrivateChannels() {
	fmt.Println("\n=== Example 9: All Private Channels ===")

	apiKey := apiKey
	secretKey := secretKey
	memo := memo

	if apiKey == "" || secretKey == "" || memo == "" {
		log.Fatal("Please set BITMART_API_KEY, BITMART_SECRET_KEY, and BITMART_MEMO environment variables")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &ws.BitMartConfig{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Memo:      memo,
		WSBaseURL: string(types.ProductionWSServer),
	}

	client, err := ws.NewClientWs(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create WebSocket client:", err)
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Close()

	fmt.Println("Connected, authenticating...")

	if err := client.Login(); err != nil {
		log.Fatal("Failed to login:", err)
	}

	time.Sleep(1 * time.Second)

	// Subscribe to all private channels
	if err := client.Private.SubscribeOrder(); err != nil {
		log.Fatal("Failed to subscribe to orders:", err)
	}

	if err := client.Private.SubscribeBalance(); err != nil {
		log.Fatal("Failed to subscribe to balance:", err)
	}

	if err := client.Private.SubscribeTrade(); err != nil {
		log.Fatal("Failed to subscribe to trades:", err)
	}

	fmt.Println("Subscribed to all private channels")
	fmt.Println("Waiting for updates (Press Ctrl+C to exit)...")

	orderCh := client.Private.GetOrderChan()
	balanceCh := client.Private.GetBalanceChan()
	tradeCh := client.Private.GetTradeChan()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case order := <-orderCh:
			fmt.Printf("\n[Order] %s: %s %s @ %s (%s)\n",
				order.OrderID, order.Side, order.Symbol, order.Price, order.Status)

		case balance := <-balanceCh:
			fmt.Printf("\n[Balance] %s: %s (Available: %s)\n",
				balance.Currency, balance.Total, balance.Available)

		case trade := <-tradeCh:
			fmt.Printf("\n[Trade] %s: %s @ %s (Fee: %s %s)\n",
				trade.Symbol, trade.Size, trade.Price, trade.Fee, trade.FeeCurrency)

		case <-sigCh:
			fmt.Println("\nShutting down...")
			return
		}
	}
}
