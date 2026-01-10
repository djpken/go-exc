package main

import (
	"context"
	"fmt"
	"log"

	"github.com/djpken/go-exc"
	"github.com/djpken/go-exc/exchanges/okex"
)

// This example demonstrates the multi-exchange architecture
// Currently only OKEx is fully implemented, but the structure
// shows how to support multiple exchanges in the future

func main() {
	ctx := context.Background()

	log.Println("=== Multi-Exchange Example ===")

	// Example 1: Using the factory (currently returns error for most exchanges)
	log.Println("--- Example 1: Using Factory Pattern ---")
	demonstrateFactory(ctx)

	// Example 2: Direct exchange instantiation (recommended for now)
	log.Println("\n--- Example 2: Direct Instantiation ---")
	demonstrateDirectInstantiation(ctx)

	// Example 3: Future multi-exchange usage pattern
	log.Println("\n--- Example 3: Future Multi-Exchange Pattern ---")
	demonstrateFuturePattern(ctx)

	log.Println("\n✅ Multi-exchange example completed!")
}

// demonstrateFactory shows how the factory pattern will work
func demonstrateFactory(ctx context.Context) {
	// Create configuration
	cfg := exc.Config{
		APIKey:     "your-api-key",
		SecretKey:  "your-secret-key",
		Passphrase: "your-passphrase",
		TestMode:   true,
	}

	// Try to create OKEx exchange through factory
	// Note: This currently returns ErrNotImplemented
	exchange, err := exc.NewExchange(ctx, exc.OKEx, cfg)
	if err != nil {
		log.Printf("Factory creation (expected to fail): %v", err)
	} else {
		log.Printf("Exchange created: %s", exchange.Name())
		defer exchange.Close()
	}

	// Try other exchanges (will fail until implemented)
	_, err = exc.NewExchange(ctx, exc.Binance, cfg)
	log.Printf("Binance: %v (not yet implemented)", err)

	_, err = exc.NewExchange(ctx, exc.BitMart, cfg)
	log.Printf("BitMart: %v (not yet implemented)", err)

	_, err = exc.NewExchange(ctx, exc.BingX, cfg)
	log.Printf("BingX: %v (not yet implemented)", err)
}

// demonstrateDirectInstantiation shows the recommended current approach
func demonstrateDirectInstantiation(ctx context.Context) {
	// For now, use direct instantiation for OKEx
	okexClient, err := exc.NewOKExClient(
		ctx,
		"your-api-key",
		"your-secret-key",
		"your-passphrase",
		true, // testMode
	)
	if err != nil {
		log.Printf("Failed to create OKEx client: %v", err)
		return
	}
	defer okexClient.Close()

	log.Printf("✅ OKEx client created: %s", okexClient.Name())

	// Access REST API through adapter
	restAPI := okexClient.REST()
	log.Printf("REST API available: %v", restAPI != nil)

	// Access WebSocket through adapter
	wsAPI := okexClient.WebSocket()
	log.Printf("WebSocket API available: %v", wsAPI != nil)

	// Access native OKEx features
	nativeClient := okexClient.GetNativeClient()
	log.Printf("Native client available: %v", nativeClient != nil)

	// Note: The unified adapter layer is under development
	// For now, use the native client directly for full functionality
	// Example: Get balance through native client
	// balance, err := nativeClient.Rest.Account.GetBalance(accountreq.GetBalance{})
	// if err != nil {
	//     log.Printf("Get balance error (expected with test credentials): %v", err)
	// } else {
	//     log.Printf("Balance retrieved: %v", balance != nil)
	// }
}

// demonstrateFuturePattern shows how the library will be used once fully implemented
func demonstrateFuturePattern(ctx context.Context) {
	log.Println("This demonstrates the future unified interface pattern:")
	log.Println()

	// Pseudo-code for future implementation
	fmt.Println(`// Future code example:

// Configure multiple exchanges
exchanges := []struct {
    name exc.ExchangeType
    cfg  exc.Config
}{
    {exc.OKEx, okexConfig},
    {exc.Binance, binanceConfig},
    {exc.BitMart, bitmartConfig},
}

// Create clients for all exchanges
for _, ex := range exchanges {
    client, err := exc.NewExchange(ctx, ex.name, ex.cfg)
    if err != nil {
        log.Printf("Failed to create %s: %v", ex.name, err)
        continue
    }
    defer client.Close()

    // Use unified interface across all exchanges
    balance, err := client.REST().Account().GetBalance(ctx)
    ticker, err := client.REST().Market().GetTicker(ctx, "BTC-USDT")
    order, err := client.REST().Trade().PlaceOrder(ctx, orderParams)

    // Subscribe to WebSocket across all exchanges
    client.WebSocket().Subscribe("ticker", params)
}`)

	log.Println()
	log.Println("Benefits of this approach:")
	log.Println("  ✓ Unified API across exchanges")
	log.Println("  ✓ Easy to switch between exchanges")
	log.Println("  ✓ Simplified portfolio management")
	log.Println("  ✓ Consistent error handling")
	log.Println("  ✓ Type-safe operations")
}

// Example helper function for managing multiple exchanges
type ExchangeManager struct {
	exchanges map[string]interface{}
}

func NewExchangeManager() *ExchangeManager {
	return &ExchangeManager{
		exchanges: make(map[string]interface{}),
	}
}

func (m *ExchangeManager) AddOKEx(ctx context.Context, apiKey, secretKey, passphrase string, testMode bool) error {
	client, err := okex.NewOKExExchange(ctx, apiKey, secretKey, passphrase, testMode)
	if err != nil {
		return err
	}
	m.exchanges["okex"] = client
	return nil
}

func (m *ExchangeManager) Get(name string) (interface{}, bool) {
	ex, ok := m.exchanges[name]
	return ex, ok
}

func (m *ExchangeManager) Close() {
	for name, ex := range m.exchanges {
		if closer, ok := ex.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				log.Printf("Error closing %s: %v", name, err)
			}
		}
	}
}
