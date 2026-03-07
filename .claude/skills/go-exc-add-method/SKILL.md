---
name: go-exc-add-method
description: >
  Add a new method to the go-exc Exchange interface and implement it across
  all exchanges (github.com/djpken/go-exc). Use this skill whenever the user
  wants to extend the unified Exchange interface with a new operation, says
  things like "add GetFundingRate to all exchanges", "新增統一方法", "add a
  new interface method", "implement X across all exchanges", or wants to
  expose a new capability through the common Exchange interface. This covers
  adding the method signature to exc.go, implementing it (or returning
  ErrNotSupported) in every exchange, and updating exc.go type re-exports
  if new request/response types are needed.
---

# Adding a New Method to the Exchange Interface

## Overview

The Exchange interface lives in `exc.go`. Every exchange in `exchanges/` must implement every method. If an exchange doesn't support a feature, it returns `commontypes.ErrNotSupported`.

Module: `github.com/djpken/go-exc`
Current exchanges: **bitmart**, **okex**, **bingx**

---

## Step 1 — Design the Method Signature

Before writing code, determine:
- **Method name** (follow existing naming: `GetXxx`, `SetXxx`, `SubscribeXxx`, `UnsubscribeXxx`)
- **Parameters**: use common types from `types/` package; avoid exchange-specific types
- **Return type**: return a pointer to a new or existing common type, plus `error`
- **Which exchanges support it**: identify who will implement vs return `ErrNotSupported`

Look at existing similar methods in `exc.go` for naming and signature conventions.

---

## Step 2 — Add Request/Response Types (if needed)

If the new method needs new types, add them to `types/`:

**Determine which file** — look at existing types:
- `types/market.go` — market data (Ticker, OrderBook, Candle, Instrument)
- `types/order.go` — orders and trading
- `types/balance.go` — account balances
- `types/position.go` — positions and leverage
- `types/websocket.go` — WebSocket update types

**Type pattern**:
```go
// Request type (if input has multiple fields)
type GetXxxRequest struct {
    Symbol   string
    Limit    int
    // ... exchange-agnostic fields only
    Extra    map[string]interface{} // for exchange-specific params
}

// Response type (if returning structured data)
type Xxx struct {
    Symbol    string
    Value     Decimal
    Timestamp Timestamp
    Extra     map[string]interface{} // for exchange-specific fields
}
```

---

## Step 3 — Add to Exchange Interface (`exc.go`)

Add the method signature to the `Exchange` interface with a doc comment:

```go
// GetXxx does [description].
// [parameter explanation]
// Returns: [what it returns]
// Note: Not all exchanges support this ([ExchangeName] returns ErrNotSupported)
GetXxx(ctx context.Context, req GetXxxRequest) (*Xxx, error)
```

If the new type needs to be accessible via the `exc` package, add type aliases in `exc.go`:
```go
type (
    // ... existing aliases
    GetXxxRequest = types.GetXxxRequest
    Xxx           = types.Xxx
)
```

---

## Step 4 — Implement in Each Exchange

For each exchange (`bitmart`, `okex`, `bingx`), do one of:

### A) Full implementation

1. Add the actual API call in the appropriate `rest/` or `ws/` file
2. Add a converter method in `converter.go`
3. Add the adapter method in `rest_adapter.go` or `ws_adapter.go`
4. Add the interface method in `<exchange>.go` delegating to the adapter:

```go
func (e *<Name>Exchange) GetXxx(ctx context.Context, req commontypes.GetXxxRequest) (*commontypes.Xxx, error) {
    return e.restAPI.<Group>().GetXxx(ctx, req)
}
```

### B) Not supported

Add directly in `<exchange>.go`:
```go
func (e *<Name>Exchange) GetXxx(_ context.Context, _ commontypes.GetXxxRequest) (*commontypes.Xxx, error) {
    return nil, commontypes.ErrNotSupported
}
```

---

## Step 5 — Update the Support Table in CLAUDE.md

After implementing, update the support matrix table in `CLAUDE.md` if it's tracked there:

```markdown
| `GetXxx()` | ✅ BitMart | ❌ OKEx | ✅ BingX |
```

---

## Step 6 — Build Check

```bash
go build ./...
```

The build must pass. Since Go requires all interface methods to be implemented, a compile error here means at least one exchange is missing the method.

---

## Decision Guide: Implement vs ErrNotSupported

Ask these questions:
1. Does the exchange have a corresponding REST/WS endpoint?
2. Is the data mappable to the common type without losing essential info?
3. Is it worth implementing now, or should it be a stub?

If stubbing, add a TODO comment:
```go
// GetXxx is not yet implemented for <Exchange>. The exchange supports this
// via [endpoint], but conversion is pending.
func (e *<Name>Exchange) GetXxx(_ context.Context, _ commontypes.GetXxxRequest) (*commontypes.Xxx, error) {
    return nil, commontypes.ErrNotSupported
}
```

---

## WebSocket Method Pattern

For new Subscribe/Unsubscribe pairs:

**In `exc.go`**:
```go
SubscribeXxx(ch chan *XxxUpdate, req WebSocketSubscribeRequest) error
UnsubscribeXxx(req WebSocketSubscribeRequest) error
```

**In `ws_adapter.go`** (for supported exchanges):
```go
func (a *WebSocketAdapter) SubscribeXxx(userCh chan *commontypes.XxxUpdate, req commontypes.WebSocketSubscribeRequest) error {
    // 1. Auto-connect if needed
    // 2. Build channel/dataType string
    // 3. RegisterHandler with parse + convert + non-blocking send
    // 4. Subscribe on WS client
    return nil
}
```

**Channel storage pattern** (add to `WebSocketAdapter` struct):
```go
xxxChannels map[string]chan *commontypes.XxxUpdate  // key = symbol or req identifier
```
