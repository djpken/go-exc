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

- `examples/okex/` - OKEx specific examples
- `examples/books.go` - WebSocket order book subscription

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
