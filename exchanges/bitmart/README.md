# BitMart Exchange Integration

BitMart integration for the go-exc multi-exchange library.

## Current Status: 🚧 In Progress

### ✅ Completed

#### Phase 1: Infrastructure Preparation
- [x] Directory structure created (21 directories)
- [x] Configuration system implemented
- [x] Native client skeleton created
- [x] REST client infrastructure (HTTP client with authentication)
- [x] WebSocket client infrastructure (connection management)

#### Phase 2: REST API Native Implementation ✅
- [x] **Market Data API** - Full implementation
  - GetTicker, GetTickers (all symbols)
  - GetOrderBook (with depth options)
  - GetTrades (recent trades)
  - GetKlines (candlestick data with configurable timeframes)
  - GetSymbols, GetSymbolDetail (trading pair information)
- [x] **Account API** - Full implementation
  - GetBalance (spot wallet balance)
  - GetWalletBalance (multi-wallet support)
- [x] **Trading API** - Full implementation
  - PlaceOrder (limit, market, limit_maker, IOC orders)
  - CancelOrder, CancelAllOrders
  - GetOrder (order details)
  - GetOrders (order list with filters)
  - GetTrades (trade history)
- [x] **Funding API** - Full implementation
  - GetDepositAddress (with chain support)
  - Withdraw (with chain and tag support)
  - GetDepositHistory
  - GetWithdrawHistory

#### Core Files Created (22 files)
**Configuration & Utils:**
- `config.go` - BitMart configuration with API Key, Secret, and Memo
- `client_legacy.go` - Main native client structure
- `types/definitions.go` - BitMart-specific types and constants
- `utils/utils.go` - Utility functions (HMAC-SHA256 signature, etc.)

**REST Client:**
- `rest/client.go` - REST HTTP client with authentication
- `rest/market.go` - Market Data API (7 endpoints)
- `rest/account.go` - Account API (2 endpoints)
- `rest/trade.go` - Trading API (6 endpoints)
- `rest/funding.go` - Funding API (4 endpoints)

**Data Models** (models/):
- `market/market_models.go` - Ticker, OrderBook, Trade, Kline, Symbol
- `account/account_models.go` - Balance, WalletBalance
- `trade/trade_models.go` - Order, OrderDetail, Trade
- `funding/funding_models.go` - DepositAddress, DepositRecord, WithdrawRecord

**Request Types** (requests/rest/):
- `market/market_requests.go` - 7 request types
- `account/account_requests.go` - 2 request types
- `trade/trade_requests.go` - 6 request types
- `funding/funding_requests.go` - 4 request types

**Response Types** (responses/):
- `market/market_responses.go` - 7 response types
- `account/account_responses.go` - 2 response types
- `trade/trade_responses.go` - 6 response types
- `funding/funding_responses.go` - 4 response types

**WebSocket Client:**
- `ws/client.go` - WebSocket client infrastructure

### 🚧 Next Steps (Phase 3-6)

#### Phase 3: WebSocket API Native Implementation
- [ ] Public channels (tickers, depth, trades)
- [ ] Private channels (orders, positions, balance)
- [ ] Event handling and distribution

#### Phase 4: Adapter Layer Implementation
- [ ] BitMart Exchange interface implementation
- [ ] Type converters (BitMart ↔ Common types)
- [ ] REST adapter
- [ ] WebSocket adapter
- [ ] Factory support

#### Phase 5: Examples and Documentation
- [ ] REST API usage examples
- [ ] WebSocket usage examples
- [ ] Multi-exchange examples (OKEx + BitMart)
- [ ] API documentation

#### Phase 6: Testing and Validation
- [ ] Unit tests (type conversion, config validation)
- [ ] Integration tests (actual API calls)
- [ ] Compilation verification
- [ ] Unified interface compatibility testing

## Directory Structure

```
exchanges/bitmart/
├── client_legacy.go        # Native client main structure
├── config.go               # BitMart configuration
├── rest/
│   └── client.go          # REST HTTP client with auth
├── ws/
│   └── client.go          # WebSocket client
├── types/
│   └── definitions.go     # BitMart-specific types
├── utils/
│   └── utils.go           # Utility functions
├── models/                # Data models (to be implemented)
│   ├── account/
│   ├── funding/
│   ├── market/
│   └── trade/
├── requests/              # Request types (to be implemented)
│   ├── rest/
│   └── ws/
├── responses/             # Response types (to be implemented)
│   ├── account/
│   ├── funding/
│   ├── market/
│   └── trade/
└── events/                # WebSocket events (to be implemented)
    ├── public/
    └── private/
```

## Configuration

```go
import "github.com/djpken/go-exc/exchanges/bitmart"

// Create BitMart client
client, err := bitmart.NewClient(ctx, apiKey, secretKey, memo)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// REST and WebSocket clients are available
// client.Rest (when REST APIs are implemented)
// client.Ws (when WebSocket APIs are implemented)
```

## API Reference

BitMart API Documentation: https://developer-pro.bitmart.com/

## Progress Tracking

See `.code/Requirements.md` (REQ-05) and `.code/Tasks.md` (Tasks 16-40) for detailed requirements and task breakdown.

## Notes

- **Authentication**: BitMart uses HMAC-SHA256 signature with API Key, Secret, and Memo
- **API Endpoints**: Production server at `https://api-cloud.bitmart.com`
- **WebSocket**: Production WebSocket at `wss://ws-manager-compress.bitmart.com/api?protocol=1.1`
- **Rate Limits**: Configured for 10 requests per second (adjustable)

## Contributing

Follow the same patterns as OKEx implementation in `exchanges/okex/` for consistency.
