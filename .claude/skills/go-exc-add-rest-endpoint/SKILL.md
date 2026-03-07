---
name: go-exc-add-rest-endpoint
description: >
  Add a new REST API endpoint to an existing exchange in go-exc
  (github.com/djpken/go-exc). Use this skill when the user wants to add a
  specific API endpoint to an already-implemented exchange, says things like
  "add the funding rate endpoint to BingX", "implement GetOpenOrders for
  BitMart", "新增 REST endpoint", "expose X API for exchange Y", or wants to
  call a specific exchange API endpoint that's not currently in the library.
  This covers: defining request/response structs in rest/, adding the HTTP
  call, adding a converter, and wiring it up through the adapter. Use
  go-exc-add-method instead if the goal is to add the method to the unified
  Exchange interface across all exchanges.
---

# Adding a New REST Endpoint to an Existing Exchange

## Overview

Adding an endpoint involves three layers:
1. **`rest/<file>.go`** — raw HTTP call + response struct
2. **`converter.go`** — native struct → common type
3. **`rest_adapter.go`** — adapter method using rest + converter

Module: `github.com/djpken/go-exc`
Current exchanges: **bitmart**, **okex**, **bingx**

---

## Step 1 — Identify the Endpoint

Determine:
- Which **exchange** (bitmart/okex/bingx)
- The **HTTP method** (GET/POST/PUT/DELETE)
- The **path** (e.g., `/openApi/swap/v2/quote/fundingRate`)
- Whether it's **public** (no auth) or **private** (needs signing)
- The **request parameters** (query string or body)
- The **response JSON structure** (get from exchange docs or existing `docs/` folder)
- Which **rest file** it belongs to:
  - `market.go` — public data (prices, depth, klines, instruments)
  - `account.go` — account state (balance, positions, leverage)
  - `trade.go` — order operations (place, cancel, query)

---

## Step 2 — Add Response Structs + HTTP Call (`rest/<file>.go`)

Follow the existing pattern exactly:

```go
// FundingRateData holds funding rate information for a single symbol
type FundingRateData struct {
    Symbol      string `json:"symbol"`
    FundingRate string `json:"fundingRate"`
    FundingTime int64  `json:"fundingTime"`
}

// FundingRateResponse is the full API response for funding rate
type FundingRateResponse struct {
    Code int             `json:"code"`
    Data FundingRateData `json:"data"`
}

// GetFundingRate retrieves the current funding rate for a symbol
// GET /openApi/swap/v2/quote/fundingRate
func (m *Market) GetFundingRate(symbol string) (*FundingRateResponse, error) {
    var result FundingRateResponse
    params := map[string]string{"symbol": symbol}
    if err := m.client.GETPublic("/openApi/swap/v2/quote/fundingRate", params, &result); err != nil {
        return nil, err
    }
    return &result, nil
}
```

Key rules:
- **One struct per response** — name it `<DataType>Data` for the inner data, `<DataType>Response` for the envelope
- **Doc comment** must include HTTP method and full path
- **Use `GETPublic`** for unauthenticated endpoints, `GET`/`POST`/`DELETE` for authenticated
- **All numeric fields as `string`** in JSON structs — converted to `Decimal` later in the converter
- **All timestamp fields as `int64`** (Unix milliseconds)

For endpoints returning a list:
```go
type FundingRateResponse struct {
    Code int               `json:"code"`
    Data []FundingRateData `json:"data"`
}
```

For POST endpoints, params go in the body:
```go
func (t *Trade) PlaceAlgoOrder(symbol, side string, ...) (*AlgoOrderResponse, error) {
    var result AlgoOrderResponse
    params := map[string]string{
        "symbol": symbol,
        "side":   side,
    }
    if err := t.client.POST("/openApi/swap/v2/trade/order/algo", params, &result); err != nil {
        return nil, err
    }
    return &result, nil
}
```

---

## Step 3 — Add Converter Method (`converter.go`)

Add a method to the `Converter` struct:

```go
// ConvertFundingRate converts exchange FundingRateData to the common FundingRate type
func (c *Converter) ConvertFundingRate(d *rest.FundingRateData) *commontypes.FundingRate {
    if d == nil {
        return nil
    }
    return &commontypes.FundingRate{
        Symbol:      d.Symbol,
        Rate:        c.str(d.FundingRate),
        NextFunding: commontypes.Timestamp(time.UnixMilli(d.FundingTime)),
    }
}
```

If there's no matching common type yet (i.e., this is exchange-specific), either:
- Return a `map[string]interface{}` for simple cases
- Add a new type to `types/market.go` (see `go-exc-add-method` skill)
- Use the `Extra` field on an existing type to carry exchange-specific data

---

## Step 4 — Add Adapter Method (`rest_adapter.go`)

Add to the appropriate adapter struct (`MarketAPIAdapter`, `AccountAPIAdapter`, or `TradeAPIAdapter`):

```go
func (a *MarketAPIAdapter) GetFundingRate(_ context.Context, symbol string) (*commontypes.FundingRate, error) {
    resp, err := a.client.Market.GetFundingRate(symbol)
    if err != nil {
        return nil, err
    }
    return a.converter.ConvertFundingRate(&resp.Data), nil
}
```

For list responses:
```go
func (a *MarketAPIAdapter) GetFundingRates(_ context.Context) ([]*commontypes.FundingRate, error) {
    resp, err := a.client.Market.GetFundingRates()
    if err != nil {
        return nil, err
    }
    result := make([]*commontypes.FundingRate, 0, len(resp.Data))
    for i := range resp.Data {
        result = append(result, a.converter.ConvertFundingRate(&resp.Data[i]))
    }
    return result, nil
}
```

---

## Step 5 — Build Check

```bash
go build ./...
```

Must pass cleanly. Typical mistakes:
- Missing import in `rest_adapter.go` (`"time"`, `"fmt"`, `"strconv"`)
- Struct field name mismatch between response JSON and converter

---

## Exposing via the Unified Interface (Optional)

If the new endpoint should be part of the `Exchange` interface (callable from any exchange), also follow the **go-exc-add-method** skill after this one.

If it's exchange-specific only, users access it via the adapter:
```go
client.REST().(*bingx.RESTAdapter).Market().GetFundingRate(ctx, "BTC-USDT")
```

---

## Quick Reference: REST Client Methods

| Method | Auth | Body |
|--------|------|------|
| `GETPublic(path, params, result)` | None | Query string |
| `GET(path, params, result)` | Signed | Query string |
| `POST(path, params, result)` | Signed | URL-encoded body |
| `PUT(path, params, result)` | Signed | Query string |
| `DELETE(path, params, result)` | Signed | Query string |

All methods decode the exchange response envelope and return a business error if the status code is non-zero.
