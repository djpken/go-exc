# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`go-exc` is a unified Go library for cryptocurrency exchange APIs, supporting multiple exchanges with a consistent interface.

**Key Design Principle**: Exchange interface with unified methods
- All exchanges implement the same `Exchange` interface
- Write code once, works with any exchange
- No need for exchange-specific logic in your application code

## Quick Start

### Build & Run

```bash
# Build the entire project
go build ./...

# Run simple example (single exchange usage)
go run examples/simple_example/main.go

# Run GetConfig example (handling unsupported features)
go run examples/getconfig_example/main.go

# Run native API examples
go run examples/bitmart_rest/main.go
go run examples/okex_rest/main.go
```

### Simple Usage

```go
// 1. Create exchange client
client, _ := exc.NewExchange(ctx, exc.BitMart, config)
defer client.Close()

// 2. Use unified methods - works for ALL exchanges
ticker, _ := client.GetTicker(ctx, "BTC_USDT")
balance, _ := client.GetBalance(ctx)
order, _ := client.PlaceOrder(ctx, exc.PlaceOrderRequest{...})
```

## Core Concepts

### Exchange Interface

All exchanges implement the `Exchange` interface with unified methods:

```go
type Exchange interface {
    // Basic methods
    Name() string
    REST() interface{}      // Exchange-specific adapter
    WebSocket() interface{} // Exchange-specific adapter
    Close() error

    // Unified API - works the same across all exchanges
    GetTicker(ctx, symbol) (*Ticker, error)
    GetOrderBook(ctx, symbol, depth) (*OrderBook, error)
    GetConfig(ctx) (*AccountConfig, error)    // ⚠️ BitMart: not supported
    GetBalance(ctx, currencies...) (*Balance, error)
    GetPositions(ctx, symbols...) ([]*Position, error)
    PlaceOrder(ctx, req) (*Order, error)
    CancelOrder(ctx, req) error
    GetOrder(ctx, req) (*Order, error)
    GetDepositAddress(ctx, currency) (string, error)
    Withdraw(ctx, req) (string, error)
}
```

**Key Benefits:**
- ✅ Same interface for all exchanges
- ✅ Type-safe operations
- ✅ Easy to switch exchanges (just change config)
- ✅ Write generic functions that work with any exchange

### Creating Exchange Clients

**Pattern 1: Factory Pattern (Recommended)**
```go
client, err := exc.NewExchange(ctx, exc.BitMart, exc.Config{
    APIKey:    "your-api-key",
    SecretKey: "your-secret-key",
    Extra:     map[string]interface{}{"memo": "your-memo"}, // BitMart only
})
```

**Pattern 2: Direct Client Creation**
```go
// BitMart
client, _ := exc.NewBitMartClient(ctx, apiKey, secretKey, memo)

// OKEx
client, _ := exc.NewOKExClient(ctx, apiKey, secretKey, passphrase, testMode)
```

## Directory Structure

```
go-exc/
├── exc.go                  # Exchange interface definition
├── factory.go              # Factory for creating exchanges
├── config.go               # Configuration types
├── errors.go               # Error definitions (including ErrNotSupported)
├── types/                  # Common types shared across exchanges
│   ├── common.go           # Decimal, Timestamp, Request types
│   ├── order.go            # Order types and constants
│   ├── balance.go          # Balance types
│   ├── position.go         # Position types
│   └── market.go           # Market data types
├── exchanges/
│   ├── bitmart/            # BitMart implementation
│   │   ├── bitmart.go      # Exchange interface implementation
│   │   ├── converter.go    # Type conversion
│   │   ├── rest_adapter.go # REST API adapter
│   │   ├── ws_adapter.go   # WebSocket adapter
│   │   └── rest/           # Native REST API
│   └── okex/               # OKEx implementation
│       ├── okex.go         # Exchange interface implementation
│       ├── converter.go    # Type conversion
│       ├── rest_adapter.go # REST API adapter
│       ├── ws_adapter.go   # WebSocket adapter
│       └── rest/           # Native REST API
└── examples/
    ├── simple_example/     # Single exchange usage (unified interface)
    ├── getconfig_example/  # Handling unsupported features
    ├── bitmart_rest/       # BitMart REST API (native client)
    ├── bitmart_ws/         # BitMart WebSocket (native client)
    ├── okex_rest/          # OKEx REST API (native client)
    └── okex_ws/            # OKEx WebSocket (native client)
```

## Configuration

### BitMart

```go
config := exc.Config{
    APIKey:    "your-api-key",
    SecretKey: "your-secret-key",
    Extra: map[string]interface{}{
        "memo": "your-memo", // Required: BitMart UID
    },
}

client, _ := exc.NewExchange(ctx, exc.BitMart, config)
```

**Exchange Types:**
- `exc.BitMart` - Production
- `exc.BitMartTEST` - Test (no separate test server, same as production)

**Symbol Format:** `BTC_USDT` (underscore)

### OKEx

```go
config := exc.Config{
    APIKey:     "your-api-key",
    SecretKey:  "your-secret-key",
    Passphrase: "your-passphrase", // Required
}

client, _ := exc.NewExchange(ctx, exc.OKX, config)
```

**Exchange Types:**
- `exc.OKX` - Production
- `exc.OKXTest` - Demo server (testnet)

**Symbol Format:** `BTC-USDT` (hyphen)

## Unified Types System

All exchanges return the same types, defined in `types/` package:

```go
// Market Data
type Ticker struct {
    Symbol    string
    LastPrice Decimal
    BidPrice  Decimal
    AskPrice  Decimal
    High24h   Decimal
    Low24h    Decimal
    Volume24h Decimal
}

type OrderBook struct {
    Symbol string
    Bids   []OrderBookLevel
    Asks   []OrderBookLevel
}

// Trading
type Order struct {
    ID                string
    Symbol            string
    Side              string   // "buy" / "sell"
    Type              string   // "limit" / "market"
    Status            string
    Price             Decimal
    Quantity          Decimal
    FilledQuantity    Decimal
    RemainingQuantity Decimal
}

// Account
type AccountBalance struct {
    Balances []Balance
}

type Balance struct {
    Currency  string
    Available Decimal
    Frozen    Decimal
    Total     Decimal
}
```

## Handling Unsupported Features

Some exchanges don't support certain features. They return `exc.ErrNotSupported`:

```go
config, err := client.GetConfig(ctx)

if errors.Is(err, exc.ErrNotSupported) {
    log.Println("This exchange doesn't support GetConfig")
    // Continue with other operations...
    return
}
```

**Unsupported Features:**
| Method | BitMart | OKEx |
|--------|---------|------|
| `GetConfig()` | ❌ | ✅ |

## Writing Exchange-Agnostic Code

The power of the `Exchange` interface is that the same code works with **any exchange**:

```go
// Create different exchanges
bitmart, _ := exc.NewExchange(ctx, exc.BitMart, bitmartConfig)
okex, _    := exc.NewExchange(ctx, exc.OKX, okexConfig)

// Use SAME code for both
ticker1, _ := bitmart.GetTicker(ctx, "BTC_USDT")
ticker2, _ := okex.GetTicker(ctx, "BTC-USDT")

// Or write generic helper functions
func getPrice(exchange exc.Exchange, symbol string) (string, error) {
    ticker, err := exchange.GetTicker(context.Background(), symbol)
    if err != nil {
        return "", err
    }
    return ticker.LastPrice.String(), nil
}
```

## Examples Guide

### `examples/simple_example/` - Single Exchange Usage

**Purpose:** Shows basic usage of a single exchange with the unified interface

**Key Features:**
- Configuration constants at the top for easy modification
- Single exchange setup and usage
- Market data queries (ticker, order book)
- Account queries (commented out, needs valid keys)
- Trading examples (commented out, needs valid keys)

```go
const (
    EXCHANGE_TYPE = exc.BitMart  // Easy to switch: exc.OKX, etc.
    API_KEY    = "your-key"
    SECRET_KEY = "your-secret"
    SYMBOL     = "BTC_USDT"
)

func main() {
    client, _ := createExchangeClient(ctx)
    ticker, _ := client.GetTicker(ctx, SYMBOL)
    orderBook, _ := client.GetOrderBook(ctx, SYMBOL, 5)
}
```

### `examples/getconfig_example/` - Handling Unsupported Features

**Purpose:** Demonstrates how to gracefully handle features not supported by all exchanges

**Key Features:**
- Shows both BitMart (not supported) and OKEx (supported)
- Error handling with `errors.Is(err, exc.ErrNotSupported)`
- Generic function that works with any exchange

```go
config, err := client.GetConfig(ctx)
if errors.Is(err, exc.ErrNotSupported) {
    log.Println("This exchange doesn't support GetConfig")
    // Continue with other operations...
}
```

### Native API Examples

**Purpose:** Show how to use exchange-specific native clients for advanced features

- `examples/bitmart_rest/` - BitMart REST API direct usage
- `examples/bitmart_ws/` - BitMart WebSocket subscriptions
- `examples/okex_rest/` - OKEx REST API direct usage
- `examples/okex_ws/` - OKEx WebSocket subscriptions

Use native clients when you need:
- Exchange-specific features not in the unified API
- Full control over API parameters
- Direct access to all exchange endpoints

## Type Conversion Architecture

Each exchange has two converters:

### 1. Type Converter (`converter.go`)
Converts data structures between native and common types:
- `ConvertOrder()` - Order conversion
- `ConvertBalance()` - Balance conversion
- `ConvertTicker()` - Ticker conversion
- `ConvertAccountConfig()` - Config conversion

### 2. Constants Converter (`constants_converter.go`)
Converts enums/constants:
- Order side: `"buy"/"sell"`
- Order type: `"limit"/"market"`
- Order status: `"open"/"filled"/"canceled"`

## Development Status

- ✅ **OKEx**: Fully implemented (REST + WebSocket)
- ✅ **BitMart**: REST API complete, WebSocket in progress
- 📋 **Binance**: Planned
- 📋 **BingX**: Planned

## Adding New Exchanges

1. Create directory: `exchanges/<exchange>/`
2. Implement `Exchange` interface in `<exchange>.go`
3. Create type converters: `converter.go` and `constants_converter.go`
4. Implement adapters: `rest_adapter.go` and `ws_adapter.go`
5. Add to factory in `factory.go`
6. Add configuration in `config.go`

## Common Patterns

### Switching Exchanges

Only need to change one line:

```go
// From BitMart to OKEx - just change this:
client, _ := exc.NewExchange(ctx, exc.OKX, config) // Was exc.BitMart

// Rest of code stays the same!
ticker, _ := client.GetTicker(ctx, symbol)
balance, _ := client.GetBalance(ctx)
```

### Batch Processing Multiple Exchanges

```go
exchanges := []exc.Exchange{bitmartClient, okexClient}

for _, ex := range exchanges {
    ticker, _ := ex.GetTicker(ctx, symbol)
    log.Printf("%s price: %s", ex.Name(), ticker.LastPrice)
}
```

### Error Handling

```go
ticker, err := client.GetTicker(ctx, "BTC_USDT")
if err != nil {
    if errors.Is(err, exc.ErrNotSupported) {
        // Feature not supported
    } else if errors.Is(err, exc.ErrRateLimitExceeded) {
        // Rate limit hit
    } else {
        // Other error
    }
}
```

## Best Practices

1. **Always use constants for configuration** - See `examples/simple_example/main.go`
2. **Write generic functions with Exchange interface** - See `examples/interface_example/main.go`
3. **Handle unsupported features gracefully** - Check for `ErrNotSupported`
4. **Use defer for Close()** - Always clean up resources
5. **Check exchange-specific symbol formats** - BitMart uses `_`, OKEx uses `-`

## API Authentication

### BitMart
```
signature = HmacSHA256(timestamp + "#" + memo + "#" + queryString, secretKey)
Headers: X-BM-KEY, X-BM-SIGN, X-BM-TIMESTAMP
```

### OKEx
```
signature = Base64(HmacSHA256(timestamp + method + path + body, secretKey))
Headers: OK-ACCESS-KEY, OK-ACCESS-SIGN, OK-ACCESS-TIMESTAMP, OK-ACCESS-PASSPHRASE
```

## Dependencies

- **Go Version**: 1.25.1
- **WebSocket**: `github.com/gorilla/websocket v1.5.3`

No other external dependencies to keep the library lightweight.
