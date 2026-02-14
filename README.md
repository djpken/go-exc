# go-exc

A unified Go library for cryptocurrency exchange APIs, supporting multiple exchanges with a consistent interface.

**Migrated from `go-okex` with multi-exchange support architecture.**

## Features

- 🔄 **Unified Interface**: Common API across multiple exchanges
- 🎯 **Type-Safe**: Strong typing with proper error handling
- 📦 **Modular Design**: Use native exchange APIs or unified interface
- 🔌 **Extensible**: Easy to add support for new exchanges
- ⚡ **High Performance**: Efficient WebSocket and REST implementations
- 🔙 **Backward Compatible**: Existing `go-okex` code works with minimal changes

## Supported Exchanges

- ✅ **OKEx** (Full native client support + adapter layer)
- 🚧 **BitMart** (In Progress - Infrastructure ready)
- 🚧 **Binance** (Planned)
- 🚧 **BingX** (Planned)

## Installation

```bash
go get github.com/djpken/go-exc
```

## Requirements

- Go 1.25 or higher
- Dependencies managed via `go.mod`

## Quick Start

### Using OKEx Native Client (Recommended)

The native OKEx client provides full access to all OKEx features:

```go
package main

import (
    "context"
    "log"

    "github.com/djpken/go-exc/exchanges/okex"
)

func main() {
    ctx := context.Background()

    // Create OKEx client
    client, err := okex.NewClient(
        ctx,
        "your-api-key",
        "your-secret-key",
        "your-passphrase",
        okex.NormalServer, // or okex.DemoServer for testing
    )
    if err != nil {
        log.Fatal(err)
    }

    // Access all native OKEx APIs
    // REST API examples
    balance, err := client.Rest.Account.GetBalance(/* params */)
    positions, err := client.Rest.Account.GetPositions(/* params */)
    orderResp, err := client.Rest.Trade.PlaceOrder(/* params */)

    // WebSocket examples
    // Subscribe to order book
    orderBookChan := make(chan *public.OrderBook)
    err = client.Ws.Public.OrderBook(req, orderBookChan)

    // Subscribe to private channels
    positionChan := make(chan *private.BalanceAndPosition)
    err = client.Ws.Private.BalanceAndPosition(positionChan)
}
```

### Using New Adapter Layer (Experimental)

The adapter layer provides a unified interface across exchanges:

```go
package main

import (
    "context"
    "log"

    "github.com/djpken/go-exc"
)

func main() {
    ctx := context.Background()

    // Create OKEx client through adapter
    client, err := exc.NewOKExClient(
        ctx,
        "api-key",
        "secret-key",
        "passphrase",
        false, // testMode
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Use adapter APIs
    balance, err := client.REST().Account().GetBalance(ctx)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Total Equity: %s\n", balance.TotalEquity)

    // Or access native client for full features
    nativeClient := client.GetNativeClient()
    // Use all native OKEx features
}
```

### WebSocket Ticker Subscriptions (Unified Interface)

Subscribe to real-time ticker updates across any exchange:

```go
package main

import (
    "context"
    "fmt"
    "log"

    exc "github.com/djpken/go-exc"
)

func main() {
    ctx := context.Background()

    // Works with any exchange - just change exc.OKX to exc.BitMart, etc.
    client, err := exc.NewExchange(ctx, exc.OKX, exc.Config{
        APIKey:     "your-api-key",
        SecretKey:  "your-secret-key",
        Passphrase: "your-passphrase",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create channel for ticker updates
    tickerCh := make(chan *exc.TickerUpdate, 100)

    // Subscribe to symbols
    symbols := []string{"BTC-USDT", "ETH-USDT"}
    if err := client.SubscribeTickers(tickerCh, symbols...); err != nil {
        log.Fatal(err)
    }

    // Process real-time ticker updates
    for update := range tickerCh {
        fmt.Printf("%s: %s (24h change: %s%%)\n",
            update.Symbol, update.LastPrice, update.PercentChange24h)
    }
}
```

### WebSocket Candlestick Subscriptions (Unified Interface)

Subscribe to real-time candlestick/kline updates across any exchange:

```go
package main

import (
    "context"
    "fmt"
    "log"

    exc "github.com/djpken/go-exc"
)

func main() {
    ctx := context.Background()

    // Works with any exchange - just change exc.OKX to exc.BitMart, etc.
    client, err := exc.NewExchange(ctx, exc.OKX, exc.Config{
        APIKey:     "your-api-key",
        SecretKey:  "your-secret-key",
        Passphrase: "your-passphrase",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create channel for candle updates
    candleCh := make(chan *exc.CandleUpdate, 100)

    // Subscribe to symbols with specific interval
    // Intervals: "1m", "5m", "15m", "30m", "1H", "4H", "1D", etc.
    symbols := []string{"BTC-USDT", "ETH-USDT"}
    if err := client.SubscribeCandles(candleCh, "1m", symbols...); err != nil {
        log.Fatal(err)
    }

    // Process real-time candle updates
    for update := range candleCh {
        if update.Confirmed {
            fmt.Printf("%s candle closed: O=%s H=%s L=%s C=%s Vol=%s\n",
                update.Symbol, update.Open, update.High, update.Low,
                update.Close, update.Volume)
        } else {
            fmt.Printf("%s candle forming: Current=%s\n",
                update.Symbol, update.Close)
        }
    }
}
```

## Migration from go-okex

### No Code Changes Required!

The library maintains full backward compatibility with `go-okex`. Simply update your import paths:

```go
// Old
import "github.com/djpken/go-okex"
import "github.com/djpken/go-okex/api"

// New
import "github.com/djpken/go-exc/exchanges/okex"
// All APIs remain the same!
```

### Module Path Update

Update your `go.mod`:

```go
// Old
module github.com/djpken/go-okex

// New
module github.com/djpken/go-exc
```

Then update import paths throughout your codebase:

```bash
find . -type f -name "*.go" -exec sed -i '' 's|github.com/djpken/go-okex|github.com/djpken/go-exc/exchanges/okex|g' {} +
```

## Project Structure

```
go-exc/
├── exc.go                    # Core unified interfaces
├── config.go                 # Configuration types
├── errors.go                 # Error definitions
├── factory.go                # Exchange factory
├── types/                    # Common types across exchanges
│   ├── common.go            # Decimal, Timestamp
│   ├── order.go             # Order types
│   ├── balance.go           # Balance types
│   ├── position.go          # Position types
│   └── market.go            # Market data types
└── exchanges/
    └── okex/                # OKEx implementation
        ├── okex.go          # Exchange adapter
        ├── converter.go     # Type converters
        ├── rest_adapter.go  # REST API adapter
        ├── ws_adapter.go    # WebSocket adapter
        ├── client_legacy.go # Original client
        ├── types/           # OKEx-specific types
        │   └── definitions.go
        ├── rest/            # REST API implementation
        ├── ws/              # WebSocket implementation
        ├── models/          # Data models
        ├── requests/        # Request types
        ├── responses/       # Response types
        └── events/          # Event types
```

## Decimal Type for Precise Financial Calculations

The `Decimal` type is used throughout go-exc for representing monetary values, prices, and quantities with precision. Unlike `float64`, `Decimal` avoids floating-point precision issues that are critical in financial applications.

### Basic Usage

```go
import exc "github.com/djpken/go-exc"

price := exc.Decimal("45678.50")
quantity := exc.Decimal("1.5")
```

### Arithmetic Operations

```go
// Basic math
sum, err := price.Add(exc.Decimal("100"))      // price + 100
diff, err := price.Sub(exc.Decimal("100"))     // price - 100
product, err := price.Mul(quantity)            // price * quantity
quotient, err := price.Div(quantity)           // price / quantity

// Unary operations
abs, err := price.Abs()    // |price|
neg, err := price.Neg()    // -price
```

### Comparison Methods

```go
// Compare
cmp, err := price1.Cmp(price2)    // Returns -1, 0, or 1

// Boolean comparisons
isLess, err := price1.LessThan(price2)               // <
isGreater, err := price1.GreaterThan(price2)         // >
isEqual, err := price1.Equal(price2)                 // ==
isLessOrEq, err := price1.LessThanOrEqual(price2)    // <=
isGreaterOrEq, err := price1.GreaterThanOrEqual(price2) // >=
```

### Max and Min Operations (New!)

```go
// Find maximum
highest, err := price1.Max(price2)

// Find minimum
lowest, err := price1.Min(price2)

// Chain operations for multiple values
highest, _ := price1.Max(price2)
highest, _ = highest.Max(price3)
```

### Sign Checks (New!)

```go
// Check if positive (> 0)
if price.IsPositive() {
    fmt.Println("Price is positive")
}

// Check if negative (< 0)
if pnl.IsNegative() {
    fmt.Println("Position in loss")
}

// Check if zero
if balance.IsZero() {
    fmt.Println("No balance")
}
```

### Practical Examples

**Find Best Execution Price Across Exchanges:**
```go
binanceBid := exc.Decimal("45678.50")
okexBid := exc.Decimal("45679.20")

// Find highest bid for selling
bestBid, _ := binanceBid.Max(okexBid)
```

**Risk Management with Stop Loss:**
```go
entryPrice := exc.Decimal("50000")
stopLossPercent := exc.Decimal("0.95")  // 95% (5% loss)

stopLossPrice, _ := entryPrice.Mul(stopLossPercent)

if currentPrice.LessThan(stopLossPrice) {
    // Trigger stop loss
}
```

**Portfolio Analysis:**
```go
position1Value, _ := position1Price.Mul(position1Qty)
position2Value, _ := position2Price.Mul(position2Qty)
totalValue, _ := position1Value.Add(position2Value)

// Find largest position
largest, _ := position1Value.Max(position2Value)
```

For more examples, see [examples/decimal_math/](examples/decimal_math/)

## Type-Safe Constants for Positions and Trading

go-exc uses type-safe constants instead of strings for position sides, margin modes, and instrument types to catch errors at compile time.

### Position Side Constants

```go
// Use typed constants instead of strings
position := &exc.Position{
    Symbol:  "BTC-USDT",
    PosSide: exc.PositionSideLong,  // Not "long" string!
}

// Available constants
exc.PositionSideLong   // Long position
exc.PositionSideShort  // Short position
exc.PositionSideNet    // Net position (one-way mode)
```

### Margin Mode Constants

```go
// Type-safe margin mode
position.MarginMode = exc.MarginModeCross     // Not "cross" string!

// Available constants
exc.MarginModeCross     // Cross margin
exc.MarginModeIsolated  // Isolated margin
```

### Instrument Type Constants

```go
// Query instruments with type-safe constants
instruments, _ := client.GetInstruments(ctx, exc.GetInstrumentsRequest{
    InstrumentType: exc.InstrumentSwap,  // Not "swap" string!
})

// Available constants
exc.InstrumentSpot      // Spot trading
exc.InstrumentMargin    // Margin trading
exc.InstrumentFutures   // Futures contracts
exc.InstrumentSwap      // Perpetual swaps
exc.InstrumentOption    // Options
```

### Setting Leverage (Type-Safe)

```go
leverage, err := client.SetLeverage(ctx, exc.SetLeverageRequest{
    Symbol:     "BTC-USDT",
    Leverage:   10,
    MarginMode: exc.MarginModeCross,    // Type-safe!
    PosSide:    exc.PositionSideLong,   // Type-safe!
})
```

**Benefits:**
- ✅ Compile-time type checking (catch typos before runtime)
- ✅ IDE autocomplete support
- ✅ Self-documenting code
- ✅ Safe refactoring

For more examples, see [examples/position_types/](examples/position_types/)

## Supported OKEx APIs

### REST API
- ✅ Trade (place/cancel/amend orders, etc.)
- ✅ Account (balance, positions, configuration)
- ✅ Funding (deposits, withdrawals, transfers)
- ✅ Market Data (tickers, order books, candles)
- ✅ Public Data (instruments, system status)
- ✅ Trading Data (volume, ratios, etc.)
- ✅ Sub-Account management

### WebSocket API
- ✅ Private Channels (account, positions, orders)
- ✅ Public Channels (tickers, order books, trades)
- ✅ Trade Operations (place/cancel orders via WS)
- ✅ Unified Ticker Subscriptions (works across all exchanges)
- ✅ Unified Candlestick Subscriptions (works across all exchanges)

## Development Status

### ✅ Phase 1: Complete
- [x] Module rename and restructuring
- [x] Dependency updates (WebSocket v1.5.3, Go 1.23)
- [x] Fixed circular dependencies
- [x] Core architecture design

### ✅ Phase 2: Complete
- [x] Core interface definitions
- [x] Common type system
- [x] Error handling framework
- [x] Factory pattern implementation

### ✅ Phase 3: Complete
- [x] OKEx adapter structure
- [x] Type converters
- [x] REST API adapters (basic operations)
- [x] WebSocket adapter structure

### 🚧 Phase 4: In Progress
- [x] README documentation
- [ ] API usage examples
- [ ] Migration guide
- [ ] Architecture documentation

### 📋 Phase 5: Planned
- [ ] Complete REST API coverage
- [ ] Full WebSocket event handling
- [ ] Binance integration
- [ ] Comprehensive test suite
- [ ] Performance benchmarks

## Examples

See the [examples/](examples/) directory for complete usage examples:

**Unified Interface Examples:**
- `examples/simple_example/` - Single exchange usage with unified interface
- `examples/getconfig_example/` - Handling unsupported features gracefully
- `examples/get_candles/` - Fetching historical candlestick/kline data
- `examples/websocket_tickers/` - Real-time ticker subscriptions via WebSocket
- `examples/websocket_candles/` - Real-time candlestick/kline subscriptions via WebSocket
- `examples/websocket_tickers_dynamic/` - Dynamic ticker subscription example
- `examples/decimal_math/` - Decimal type math operations and trading scenarios
- `examples/position_types/` - Type-safe constants for positions and trading

**Native Client Examples:**
- `examples/bitmart_rest/` - BitMart REST API usage
- `examples/bitmart_ws/` - BitMart WebSocket subscriptions
- `examples/okex_rest/` - OKEx REST API usage
- `examples/okex_ws/` - OKEx WebSocket subscriptions

## Dependencies

- **Go**: 1.23 or higher
- **gorilla/websocket**: v1.5.3

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the same terms as the original `go-okex` project. See [LICENSE](LICENSE) for details.

## Disclaimer

This package is provided as-is, without any express or implied warranties. The user assumes all risks associated with the use of this package. Use at your own risk.

## Documentation

- [CHANGELOG.md](CHANGELOG.md) - Version history and changes
- [MIGRATION.md](MIGRATION.md) - Migration guide from go-okex
- [Architecture Documentation](.code/ARCHITECTURE.md) - Detailed architecture design
- [Refactoring Plan](.code/REFACTORING_PLAN.md) - Project refactoring implementation plan

## Resources

- [OKEx API Documentation](https://www.okx.com/docs-v5/en/)
- [Issues](https://github.com/djpken/go-exc/issues)
- [Discussions](https://github.com/djpken/go-exc/discussions)

## Acknowledgments

Originally based on [go-okex](https://github.com/amir-the-h/okex) by amir-the-h, refactored for multi-exchange support.
