---
name: go-exc-new-exchange
description: >
  Add a new cryptocurrency exchange to the go-exc unified library
  (github.com/djpken/go-exc). Use this skill whenever the user wants to
  integrate a new exchange, says something like "add Binance support",
  "implement a new exchange", "新增交易所", or mentions adding support for
  any exchange name (Bybit, Gate.io, Kraken, Huobi, etc.). This skill
  covers the full implementation: REST client, WebSocket client, converters,
  adapters, and factory registration — all following the project's established
  patterns.
---

# Adding a New Exchange to go-exc

## Overview

Each exchange lives in `exchanges/<name>/` and must implement the `Exchange` interface defined in `exc.go`. Follow the BingX implementation as the canonical reference — it's the most recently added and cleanest.

Module path: `github.com/djpken/go-exc`
Common types import alias: `commontypes "github.com/djpken/go-exc/types"`

---

## Step 1 — Gather Exchange Info

Before writing any code, determine:
- **Exchange name** (e.g., `bybit`) — used as package name and directory
- **Base REST URL** (prod + test/demo if different)
- **WebSocket URL** (pub + private, prod + test)
- **Auth method** (HMAC-SHA256 is common; check the exchange docs)
- **Auth headers** (e.g., `X-BX-APIKEY`, `X-MBX-APIKEY`, `OK-ACCESS-KEY`)
- **Response envelope** (common: `{"code":0,"msg":"","data":{...}}`)
- **Symbol format** (hyphen `BTC-USDT`, underscore `BTC_USDT`, or none `BTCUSDT`)
- **Private WS auth method** (listen-key, JWT, HMAC challenge)
- **Any extra config fields** needed (e.g., BitMart needs `memo`)

Ask the user for the above if not already provided. Check the `docs/` directory inside the exchange folder if it exists.

---

## Step 2 — Directory Structure

Create this exact layout:

```
exchanges/<name>/
├── <name>.go           # Exchange struct + interface implementation
├── converter.go        # Native ↔ common type conversion
├── rest_adapter.go     # RESTAdapter + Market/Account/TradeAPIAdapter
├── ws_adapter.go       # WebSocketAdapter + Subscribe*/Unsubscribe* methods
├── rest/
│   ├── client.go       # ClientRest, auth signing, GET/POST/DELETE helpers
│   ├── market.go       # Market data endpoints + response structs
│   ├── account.go      # Account/position/leverage endpoints + structs
│   └── trade.go        # Order endpoints + structs
└── ws/
    ├── client.go       # Public WebSocket client (GZIP decode if needed)
    └── private.go      # Private WebSocket client (listen-key or auth)
```

---

## Step 3 — REST Client (`rest/client.go`)

Pattern to follow (from BingX):

```go
package rest

type ClientRest struct {
    ctx        context.Context
    httpClient *http.Client
    apiKey     string
    secretKey  string
    baseURL    string

    Market  *Market
    Account *Account
    Trade   *Trade
}

// Response envelope — adapt to what this exchange actually returns
type Response[T any] struct {
    Code      int    `json:"code"`
    Msg       string `json:"msg"`
    Data      T      `json:"data"`
    Timestamp int64  `json:"timestamp"`
}

func NewClientRest(ctx context.Context, apiKey, secretKey, baseURL string) *ClientRest {
    c := &ClientRest{...}
    c.Market = NewMarket(c)
    c.Account = NewAccount(c)
    c.Trade = NewTrade(c)
    return c
}
```

Key points:
- Implement `sign()` using the exchange's auth algorithm (usually HMAC-SHA256 hex or base64)
- Provide `GET`, `POST`, `DELETE`, `PUT` helpers for authenticated requests
- Provide `GETPublic` for unauthenticated market data
- Check the business error code in `executeAndDecode` (non-zero code = API error)
- Set `Timeout: 30 * time.Second` on the HTTP client

---

## Step 4 — REST Endpoint Files

Each file follows this pattern:

```go
// market.go
package rest

type Market struct { client *ClientRest }
func NewMarket(c *ClientRest) *Market { return &Market{client: c} }

// Response structs — one per endpoint
type TickerData struct { ... }
type TickerResponse struct {
    Code int        `json:"code"`
    Data TickerData `json:"data"`
}

// One method per endpoint with doc comment showing HTTP method + path
// GET /openApi/swap/v2/quote/ticker
func (m *Market) GetTicker(symbol string) (*TickerResponse, error) { ... }
```

Required methods per file:
- **market.go**: `GetTicker`, `GetTickers`, `GetOrderBook`, `GetKlines`, `GetContracts/GetInstruments`
- **account.go**: `GetBalance`, `GetPositions`, `GetLeverage`, `SetLeverage`
- **trade.go**: `PlaceOrder`, `CancelOrder`, `GetOrder`

---

## Step 5 — WebSocket Client (`ws/client.go`)

Public WebSocket client must support:
- `Connect()` / `Close()` / `IsConnected()`
- `Subscribe(channel)` / `Unsubscribe(channel)`
- `RegisterHandler(channel, func([]byte))` / `UnregisterHandler(channel)`
- Auto-reconnect loop
- GZIP decompression if the exchange compresses messages

Private WebSocket (`ws/private.go`) usually needs:
- Listen-key management (create/extend/delete via REST) OR HMAC auth challenge
- `RegisterHandler(eventType, func([]byte))`
- Auto-reconnect with re-auth

---

## Step 6 — Converter (`converter.go`)

```go
package <name>

type Converter struct{}
func NewConverter() *Converter { return &Converter{} }

// Helper: safe string → Decimal
func (c *Converter) str(s string) commontypes.Decimal {
    if s == "" { return commontypes.ZeroDecimal }
    d, err := commontypes.NewDecimal(s)
    if err != nil { return commontypes.ZeroDecimal }
    return d
}
```

Implement these converters (adapt field names to exchange's actual response):
- `ConvertTicker(*rest.TickerData) *commontypes.Ticker`
- `ConvertOrderBook(*rest.OrderBookData, symbol) *commontypes.OrderBook`
- `ConvertKline(*rest.KlineEntry, symbol, interval) *commontypes.Candle`
- `ConvertBalance(*rest.BalanceData) *commontypes.AccountBalance`
- `ConvertPosition(*rest.PositionData) *commontypes.Position`
- `ConvertLeverage(symbol, *rest.LeverageData) []*commontypes.Leverage`
- `ConvertOrder(*rest.OrderData) *commontypes.Order`
- `ConvertInstrument(*rest.InstrumentData) *commontypes.Instrument`
- `ConvertOrderStatus(string) string` — map exchange statuses to common ones
- `ConvertIntervalToWS(string) (string, error)` — map common intervals to WS format
- `ConvertIntervalToREST(string) (string, error)` — map common intervals to REST format

Common interval map (adjust for exchange):
```go
// Common → REST
"1m"→"1m", "3m"→"3m", "5m"→"5m", "15m"→"15m", "30m"→"30m",
"1H"→"1h", "2H"→"2h", "4H"→"4h", "6H"→"6h", "12H"→"12h",
"1D"→"1d", "1W"→"1w", "1M"→"1M"
```

Position side mapping:
- `"LONG"` → `commontypes.PositionSideLong`
- `"SHORT"` → `commontypes.PositionSideShort`
- anything else → `commontypes.PositionSideNet`

---

## Step 7 — REST Adapter (`rest_adapter.go`)

```go
package <name>

type RESTAdapter struct {
    client    *rest.ClientRest
    converter *Converter
}
func NewRESTAdapter(c *rest.ClientRest) *RESTAdapter { ... }
func (a *RESTAdapter) Market() *MarketAPIAdapter  { ... }
func (a *RESTAdapter) Account() *AccountAPIAdapter { ... }
func (a *RESTAdapter) Trade() *TradeAPIAdapter    { ... }

type MarketAPIAdapter struct { client *rest.ClientRest; converter *Converter }
// Implement: GetTicker, GetTickers, GetOrderBook, GetCandles, GetInstruments

type AccountAPIAdapter struct { ... }
// Implement: GetBalance, GetPositions, GetLeverage, SetLeverage

type TradeAPIAdapter struct { ... }
// Implement: PlaceOrder, CancelOrder, GetOrderDetail
```

Each adapter method:
1. Calls the corresponding `rest/` method
2. Converts the response using the `Converter`
3. Returns the common type

---

## Step 8 — WebSocket Adapter (`ws_adapter.go`)

```go
type WebSocketAdapter struct {
    client         *ws.ClientWs
    privateClient  *ws.PrivateClientWs
    converter      *Converter
    tickerChannels map[string]chan *commontypes.TickerUpdate
    candleChannels map[string]map[string]chan *commontypes.CandleUpdate
}
```

Implement all Subscribe/Unsubscribe methods. Pattern for each:
1. Check `IsConnected()`, auto-connect if not
2. Build the exchange-specific channel/dataType string
3. `RegisterHandler(dataType, func([]byte) { parse → convert → send to userCh })`
4. Call `Subscribe(dataType)` on the WS client
5. Store channel reference in the map

For unsupported methods, return `commontypes.ErrNotSupported`.

---

## Step 9 — Main Exchange File (`<name>.go`)

```go
package <name>

const (
    defaultRESTURL = "https://..."
    testRESTURL    = "https://..."  // or same as prod
    defaultWSURL   = "wss://..."
    testWSURL      = "wss://..."
)

type <Name>Exchange struct {
    restClient *rest.ClientRest
    wsClient   *ws.ClientWs
    privateWS  *ws.PrivateClientWs
    restAPI    *RESTAdapter
    wsAPI      *WebSocketAdapter
    ctx        context.Context
    testMode   bool
}

func New<Name>Exchange(ctx context.Context, apiKey, secretKey string, testMode bool) (*<Name>Exchange, error) { ... }
```

Implement all Exchange interface methods by delegating to `e.restAPI` or `e.wsAPI`.

For unsupported features:
```go
func (e *<Name>Exchange) GetConfig(_ context.Context) (*commontypes.AccountConfig, error) {
    return nil, commontypes.ErrNotSupported
}
```

---

## Step 10 — Register in Factory

**`config.go`** — add constants:
```go
const (
    // existing ...
    NewExchange     ExchangeType = "NEW-EXCHANGE"
    NewExchangeTest ExchangeType = "NEW-EXCHANGE-TEST"
)
```

**`factory.go`** — add import and cases:
```go
import "github.com/djpken/go-exc/exchanges/<name>"

case NewExchange:
    return new<Name>Exchange(ctx, cfg, false)
case NewExchangeTest:
    return new<Name>Exchange(ctx, cfg, true)
```

Add the private constructor:
```go
func new<Name>Exchange(ctx context.Context, cfg Config, testMode bool) (Exchange, error) {
    client, err := <name>.New<Name>Exchange(ctx, cfg.APIKey, cfg.SecretKey, testMode)
    if err != nil { return nil, err }
    return client, nil
}
```

If the exchange needs extra config (like BitMart's `memo`):
```go
func new<Name>Exchange(ctx context.Context, cfg Config, testMode bool) (Exchange, error) {
    extra := ""
    if cfg.Extra != nil {
        if v, ok := cfg.Extra["extra_field"].(string); ok {
            extra = v
        }
    }
    client, err := <name>.New<Name>Exchange(ctx, cfg.APIKey, cfg.SecretKey, extra, testMode)
    ...
}
```

---

## Step 11 — Build Check

Run `go build ./...` and fix any compile errors. The build must pass cleanly.

---

## Common Pitfalls

- **Decimal creation**: Use `commontypes.NewDecimal(s)` or `commontypes.MustDecimal(s)` — never call `Decimal("string")` directly (compile error)
- **Zero decimal**: Use `commontypes.ZeroDecimal` for empty/zero values
- **Symbol format**: Document the expected format in the package-level comment
- **Interval format**: WS and REST may use different formats (e.g., BingX WS uses `1min` vs REST `1m`)
- **Channel direction**: User channels receive updates via non-blocking send: `select { case ch <- update: default: }`
- **Private WS**: Must handle auth before any private subscriptions
- **go.mod**: After adding the exchange, run `go mod tidy` if any new dependencies were added
