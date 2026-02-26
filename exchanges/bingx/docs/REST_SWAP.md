# BingX Perpetual Swap (Futures) REST API

## Base URL

```
https://open-api.bingx.com
```

## Symbol Format

Perpetual swap trading pairs use a hyphen separator: `BTC-USDT`, `ETH-USDT`, etc. Always use uppercase letters.

## Authentication

All private endpoints require `X-BX-APIKEY` header and `signature` parameter. See `API_OVERVIEW.md` for authentication details.

## Position Modes

BingX supports two position modes:
- **One-Way Mode (Single):** `positionSide = BOTH`
- **Hedge Mode (Dual):** `positionSide = LONG` or `SHORT`

---

## Market Data Endpoints (Public)

### Get Contract Information

```
GET /openApi/swap/v2/quote/contracts
```

Get all perpetual swap contract specifications.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| contractId | string | Contract ID |
| symbol | string | Trading pair |
| size | string | Contract size |
| quantityPrecision | int | Quantity decimal places |
| pricePrecision | int | Price decimal places |
| feeRate | string | Taker fee rate |
| tradeMinLimit | int | Minimum trade quantity |
| maxLongLeverage | int | Maximum long position leverage |
| maxShortLeverage | int | Maximum short position leverage |
| currency | string | Settlement currency |
| asset | string | Base asset |

---

### Get Order Book Depth

```
GET /openApi/swap/v2/quote/depth
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| limit | int | No | Depth levels: 5, 10, 20, 50, 100, 500, 1000 |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| bids | array | Bid orders: `[[price, quantity], ...]` sorted by price descending |
| asks | array | Ask orders: `[[price, quantity], ...]` sorted by price ascending |
| T | int64 | Last update timestamp |

---

### Get Recent Trades

```
GET /openApi/swap/v2/quote/trades
```

Get the latest trades for a trading pair.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| limit | int | No | Number of results (default: 100, max: 500) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | Trade ID |
| price | string | Trade price |
| qty | string | Trade quantity |
| quoteQty | string | Quote asset quantity |
| time | int64 | Trade timestamp in milliseconds |
| isBuyerMaker | bool | True if buyer is maker |

---

### Get Historical Trades

```
GET /openApi/swap/v1/market/historicalTrades
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| fromId | int64 | No | Trade ID to fetch from |
| limit | int | No | Number of results (default: 500, max: 1000) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Kline/Candlestick Data (v3)

```
GET /openApi/swap/v3/quote/klines
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| interval | string | Yes | Kline interval: `1m`, `3m`, `5m`, `15m`, `30m`, `1h`, `2h`, `4h`, `6h`, `8h`, `12h`, `1d`, `3d`, `1w`, `1M` |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results (default: 500, max: 1440) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| open | string | Opening price |
| close | string | Closing price |
| high | string | Highest price |
| low | string | Lowest price |
| volume | string | Trading volume |
| time | int64 | Kline open time (milliseconds) |

---

### Get Kline Data (v2)

```
GET /openApi/swap/v2/quote/klines
```

Same parameters as v3 above.

---

### Get Mark Price Klines

```
GET /openApi/swap/v1/market/markPriceKlines
```

Kline data based on mark price.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| interval | string | Yes | Kline interval |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Latest Mark Price and Funding Rate

```
GET /openApi/swap/v2/quote/premiumIndex
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair. If omitted, returns all pairs |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| markPrice | string | Current mark price |
| indexPrice | string | Index price |
| lastFundingRate | string | Last settled funding rate |
| nextFundingTime | int64 | Time until next funding settlement (milliseconds) |

---

### Get Funding Rate History

```
GET /openApi/swap/v2/quote/fundingRate
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results (default: 100, max: 1000) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| fundingRate | string | Funding rate |
| fundingTime | int64 | Funding timestamp in milliseconds |

---

### Get Open Interest

```
GET /openApi/swap/v2/quote/openInterest
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| openInterest | string | Open interest quantity |
| symbol | string | Trading pair |
| time | int64 | Matching engine timestamp |

---

### Get Best Order Book Price (Book Ticker)

```
GET /openApi/swap/v2/quote/bookTicker
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair. If omitted, returns all pairs |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| bid_price | float64 | Best bid price |
| bid_qty | float64 | Best bid quantity |
| ask_price | float64 | Best ask price |
| ask_qty | float64 | Best ask quantity |
| time | long | Transaction timestamp in milliseconds |
| lastUpdateId | int64 | Latest trade ID |

---

### Get Latest Price

```
GET /openApi/swap/v2/quote/ticker
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair. If omitted, returns all pairs |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Server Time

```
GET /openApi/swap/v2/server/time
```

No authentication required. Returns the current server timestamp.

---

## Account Endpoints (Private)

### Get Account Balance

```
GET /openApi/swap/v2/user/balance
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| currency | string | No | Currency filter (e.g. USDT) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| asset | string | Asset name |
| balance | string | Total balance |
| equity | string | Net asset value |
| unrealizedProfit | string | Unrealized PnL |
| realisedProfit | string | Realized PnL |
| availableMargin | string | Available margin |
| usedMargin | string | Used margin |
| freezedMargin | string | Frozen margin |

---

### Get Account Balance (v3 - Multi-Asset)

```
GET /openApi/swap/v3/user/balance
```

Same parameters as v2 balance endpoint above.

---

### Get Positions

```
GET /openApi/swap/v2/user/positions
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair. If omitted, returns all positions |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| positionId | string | Position ID |
| positionSide | string | Position side: LONG, SHORT, or BOTH |
| isolated | bool | True if isolated margin mode |
| positionAmt | string | Position quantity |
| availableAmt | string | Available position quantity (for closing) |
| unrealizedProfit | string | Unrealized PnL |
| realisedProfit | string | Realized PnL |
| initialMargin | string | Initial margin |
| avgPrice | string | Average entry price |
| leverage | int | Current leverage |

---

### Get Position History

```
GET /openApi/swap/v1/trade/positionHistory
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| pageIndex | int | No | Page number |
| pageSize | int | No | Results per page |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Income History (Fund Flow)

```
GET /openApi/swap/v2/user/income
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| incomeType | string | No | Fund flow type: `REALIZED_PNL`, `FUNDING_FEE`, `TRADING_FEE`, `INSURANCE_CLEAR`, `TRIAL_FUND`, `ADL`, `SYSTEM_DEDUCTION` |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results (default: 100, max: 1000) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Export Fund Flow (Excel)

```
GET /openApi/swap/v2/user/income/export
```

> **Note:** Response is an Excel file, not JSON.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| incomeType | string | No | Fund flow type |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results (max: 200) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Commission Rate

```
GET /openApi/swap/v2/user/commissionRate
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| makerCommissionRate | string | Maker fee rate |
| takerCommissionRate | string | Taker fee rate |

---

### Get Margin Assets

```
GET /openApi/swap/v1/user/marginAssets
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

## Leverage and Margin Settings (Private)

### Get Leverage

```
GET /openApi/swap/v2/trade/leverage
```

Query the current leverage for a trading pair.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| longLeverage | int64 | Current long position leverage |
| shortLeverage | int64 | Current short position leverage |
| maxLongLeverage | int64 | Maximum long leverage allowed |
| maxShortLeverage | int64 | Maximum short leverage allowed |

**Example Response:**
```json
{
  "code": 0,
  "data": {
    "longLeverage": 50,
    "shortLeverage": 50,
    "maxLongLeverage": 75,
    "maxShortLeverage": 75
  }
}
```

---

### Set Leverage

```
POST /openApi/swap/v2/trade/leverage
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| side | string | Yes | Position side: `LONG` or `SHORT` |
| leverage | int | Yes | New leverage value |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Margin Mode

```
GET /openApi/swap/v2/trade/marginType
```

Query whether the account uses cross margin or isolated margin for a trading pair.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| marginType | string | Margin mode: `CROSSED` (cross) or `ISOLATED` |

---

### Set Margin Mode

```
POST /openApi/swap/v2/trade/marginType
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| marginType | string | Yes | `CROSSED` or `ISOLATED` |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Adjust Isolated Margin

```
POST /openApi/swap/v2/trade/positionMargin
```

Add or remove margin from an isolated position.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| amount | float64 | Yes | Margin amount to adjust |
| type | int | Yes | Direction: `1` = add margin, `2` = remove margin |
| positionSide | string | Yes | Position side: `LONG` or `SHORT` |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| amount | float64 | Adjusted margin amount |
| type | int | Adjustment direction |

---

### Get Position Margin History

```
GET /openApi/swap/v1/positionMargin/history
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| type | int | No | Adjustment type filter (1 = add, 2 = remove) |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get/Set Position Side Mode (Hedge Mode)

```
GET /openApi/swap/v1/positionSide/dual
POST /openApi/swap/v1/positionSide/dual
```

Query or change the position mode (one-way vs hedge).

**Request Parameters (POST):**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| dualSidePosition | string | Yes | `true` = hedge mode, `false` = one-way mode |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

## Trading Endpoints (Private)

### Place Order

```
POST /openApi/swap/v2/trade/order
```

Place a perpetual swap order.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| side | string | Yes | Order side: `BUY` or `SELL` |
| positionSide | string | No | Position side: `BOTH` (one-way mode) or `LONG`/`SHORT` (hedge mode). Defaults to `LONG` |
| type | string | Yes | Order type (see order types below) |
| reduceOnly | string | No | `true` or `false`. Only for one-way mode. Ignored in hedge mode |
| price | float64 | No | Order price (required for LIMIT orders) |
| quantity | float64 | No | Order quantity in base asset (e.g. BTC). Use for coin-quantity orders |
| stopPrice | float64 | No | Trigger price (required for STOP/TAKE_PROFIT/TRIGGER orders) |
| priceRate | float64 | No | Callback rate for `TRAILING_STOP_MARKET` (max: 1) |
| workingType | string | No | Trigger price type: `MARK_PRICE` (default), `CONTRACT_PRICE`, `INDEX_PRICE` |
| timeInForce | string | No | `GTC` (default), `IOC`, `FOK`, `PostOnly` |
| clientOrderID | string | No | Custom order ID (1–40 chars; unique; only for LIMIT/MARKET types) |
| stopLoss | string | No | JSON string for stop-loss when placing order (type: STOP_MARKET/STOP) |
| takeProfit | string | No | JSON string for take-profit when placing order (type: TAKE_PROFIT_MARKET/TAKE_PROFIT) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| orderId | int64 | System order ID |
| side | string | Order side |
| positionSide | string | Position side |
| type | string | Order type |
| clientOrderID | string | Custom order ID |
| workingType | string | Trigger price type |

**Example Response:**
```json
{
  "code": 0,
  "data": {
    "order": {
      "symbol": "BTC-USDT",
      "orderId": 1736012449498123456,
      "side": "BUY",
      "positionSide": "LONG",
      "type": "LIMIT",
      "price": "42000",
      "origQty": "0.001",
      "status": "PENDING"
    }
  }
}
```

---

### Test Place Order (No Execution)

```
POST /openApi/swap/v2/trade/order/test
```

Validates order parameters without actually placing the order. Same request parameters as the live order endpoint.

---

### Query Order

```
GET /openApi/swap/v2/trade/order
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderId | int64 | No | System order ID |
| clientOrderID | string | No | Custom order ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

> Either `orderId` or `clientOrderID` must be provided.

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| orderId | int64 | Order ID |
| symbol | string | Trading pair |
| side | string | Buy/sell direction |
| positionSide | string | Position side |
| type | string | Order type |
| origQty | string | Original order quantity |
| price | string | Order price |
| executedQty | string | Executed quantity |
| avgPrice | string | Average transaction price |
| cumQuote | string | Total transaction amount |
| stopPrice | string | Trigger price |
| status | string | Order status |
| profit | string | Profit and loss |
| commission | string | Trading fee |
| time | int64 | Order creation time (ms) |
| updateTime | int64 | Last update time (ms) |
| clientOrderID | string | Custom order ID |
| workingType | string | Trigger price type |

---

### Query Open Order (Pending)

```
GET /openApi/swap/v2/trade/openOrder
```

Query a single open (pending) order.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderId | int64 | No | System order ID |
| clientOrderID | string | No | Custom order ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

Response fields are the same as Query Order above.

---

### Get All Open Orders

```
GET /openApi/swap/v2/trade/openOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair. If omitted, returns all pairs |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get All Orders (History)

```
GET /openApi/swap/v2/trade/allOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair. If omitted, returns all pairs |
| orderId | int64 | No | Return orders from this ID onwards |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | Yes | Number of results (default: 500, max: 1000) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

Response fields are the same as Query Order.

---

### Get Fill History (Executed Trades)

```
GET /openApi/swap/v2/trade/allFillOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| tradingUnit | string | Yes | Trading unit: `COIN` (asset quantity like BTC/ETH) or `CONT` (contract sheets) |
| startTs | int64 | Yes | Start timestamp in milliseconds |
| endTs | int64 | Yes | End timestamp in milliseconds |
| orderId | int64 | No | Return fills from this order ID onwards |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| volume | string | Fill quantity |
| price | string | Fill price |
| amount | string | Fill value |
| commission | string | Fee charged |
| currency | string | Settlement currency (usually USDT) |
| orderId | string | Order ID |
| filledTime | string | Fill timestamp (format: 2006-01-02T15:04:05.999+0800) |
| liquidatedPrice | string | Estimated liquidation price at time of forced liquidation (only for liquidation orders) |
| liquidatedMarginRatio | string | Margin rate at time of forced liquidation (only for liquidation orders) |
| workingType | string | Trigger price type |

---

### Get Fill History (v2)

```
GET /openApi/swap/v2/trade/fillHistory
```

Alternative fill history endpoint.

---

### Get Liquidation Orders (Force Orders)

```
GET /openApi/swap/v2/trade/forceOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| autoCloseType | string | No | `LIQUIDATION` (force close) or `ADL` (auto-deleveraging) |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results (default: 50, max: 100) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

Response fields are the same as Query Order.

---

### Cancel Order

```
DELETE /openApi/swap/v2/trade/order
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderId | int64 | No | System order ID to cancel |
| clientOrderID | string | No | Custom order ID to cancel |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

Response fields include canceled order details (same structure as Query Order).

---

### Cancel All Open Orders

```
DELETE /openApi/swap/v2/trade/allOpenOrders
```

Cancel all open orders for a trading pair.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| recvWindow | int64 | No | Request validity window in milliseconds |
| timestamp | int64 | Yes | Request timestamp in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| success | array | List of successfully canceled orders |
| failed | array | List of orders that failed to cancel |

---

### Cancel All After Timer

```
POST /openApi/swap/v2/trade/cancelAllAfter
```

Set a timer to auto-cancel all open orders. Acts as a dead-man's switch.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| type | string | Yes | `ACTIVATE` or `DEACTIVATE` |
| timeOut | int64 | Yes | Timeout in milliseconds |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Cancel and Replace Order

```
POST /openApi/swap/v1/trade/cancelReplace
```

Cancel an existing order and place a new one atomically.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| cancelOrderId | int64 | No | Order ID to cancel |
| cancelClientOrderId | string | No | Custom order ID to cancel |
| side | string | Yes | BUY or SELL |
| positionSide | string | No | LONG, SHORT, or BOTH |
| type | string | Yes | Order type |
| quantity | float64 | No | New order quantity |
| price | float64 | No | New order price |
| stopPrice | float64 | No | Trigger price |
| workingType | string | No | Trigger price type |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Batch Cancel and Replace Orders

```
POST /openApi/swap/v1/trade/batchCancelReplace
```

Cancel multiple orders and replace them in a batch.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| batchOrders | string | Yes | URL-encoded JSON array of cancel-replace order objects |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Batch Place Orders

```
POST /openApi/swap/v2/trade/batchOrders
```

Place multiple orders in a single request.

> **Note:** Batch orders are processed concurrently; matching order is not guaranteed.

**How to construct the batch request:**

1. Build the `batchOrders` parameter as a JSON array string:
```json
batchOrders=[{"symbol":"ETH-USDT","type":"MARKET","side":"BUY","positionSide":"LONG","quantity":1},{"symbol":"BTC-USDT","type":"MARKET","side":"BUY","positionSide":"LONG","quantity":0.001}]&timestamp=1692956597902
```

2. Sign the above string (before URL encoding).

3. URL-encode only the **value** of `batchOrders` (not the key, not the timestamp):
```
batchOrders=%5B%7B%22symbol%22%3A%22ETH-USDT%22...%7D%5D&timestamp=1692956597902
```

4. Append the signature and submit.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| batchOrders | string | Yes | URL-encoded JSON array of order objects |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

Each order object in the array uses the same fields as a single order placement.

---

### One-Click Close All Positions

```
POST /openApi/swap/v2/trade/closeAllPositions
```

Close all open positions for a symbol (or all symbols) at market price.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair. If omitted, closes all positions for all pairs |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| success | array of int64 | Order IDs of successful close orders |
| failed | array | Order IDs that failed to close |

---

### Close Position

```
POST /openApi/swap/v1/trade/closePosition
```

Close a specific position by position ID.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| positionId | string | Yes | Position ID to close |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Reverse Position

```
POST /openApi/swap/v1/trade/reverse
```

Reverse an existing position (close and open in opposite direction).

---

### Amend Order

```
PUT /openApi/swap/v1/trade/amend
```

Modify an existing open order.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderId | int64 | No | Order ID to amend |
| clientOrderId | string | No | Custom order ID to amend |
| price | float64 | No | New price |
| quantity | float64 | No | New quantity |
| stopPrice | float64 | No | New trigger price |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Maintenance Margin Ratio

```
GET /openApi/swap/v1/maintMarginRatio
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get/Set Asset Mode

```
GET /openApi/swap/v1/trade/assetMode
POST /openApi/swap/v1/trade/assetMode
```

Query or set whether to use multi-asset margin mode.

---

### Get Multi-Asset Rules

```
GET /openApi/swap/v1/trade/multiAssetsRules
```

---

### Auto Add Margin

```
POST /openApi/swap/v1/trade/autoAddMargin
```

Enable or disable automatic margin addition for isolated positions.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| side | string | Yes | `LONG` or `SHORT` |
| autoAddMargin | string | Yes | `true` or `false` |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get VST Balance

```
GET /openApi/swap/v2/trade/getVst
```

Get Virtual Simulation Token (VST) balance for paper trading.

---

## TWAP Orders

Time-Weighted Average Price (TWAP) orders execute large orders incrementally over time.

### Place TWAP Order

```
POST /openApi/swap/v1/twap/order
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| side | string | Yes | `BUY` or `SELL` |
| positionSide | string | No | `LONG`, `SHORT`, or `BOTH` |
| priceType | string | Yes | Price type: `constant` or `market` |
| priceVariance | string | No | Maximum price deviation from market price |
| triggerPrice | string | No | Trigger price |
| size | int64 | Yes | Total order quantity |
| tradeNum | int64 | Yes | Number of sub-orders to split into |
| duration | int64 | Yes | Total execution duration in seconds |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Cancel TWAP Order

```
DELETE /openApi/swap/v1/twap/cancelOrder
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| mainOrderId | string | Yes | TWAP main order ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get TWAP Order Detail

```
GET /openApi/swap/v1/twap/orderDetail
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| mainOrderId | string | Yes | TWAP main order ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Open TWAP Orders

```
GET /openApi/swap/v1/twap/openOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| pageIndex | int | No | Page number |
| pageSize | int | No | Results per page |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get TWAP Order History

```
GET /openApi/swap/v1/twap/historyOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| pageIndex | int | No | Page number |
| pageSize | int | No | Results per page |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

## Trading Rules

### Get Trading Rules

```
GET /openApi/swap/v1/tradingRules
```

Get trading rule constraints for all or specific symbols.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

## Coin-Margined Swap (CSwap) Endpoints

The `/openApi/cswap/` endpoints follow the same structure as standard perpetual swap endpoints.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/openApi/cswap/v1/market/contracts` | GET | Get coin-margined contract info |
| `/openApi/cswap/v1/market/depth` | GET | Get order book depth |
| `/openApi/cswap/v1/market/klines` | GET | Get kline data |
| `/openApi/cswap/v1/market/openInterest` | GET | Get open interest |
| `/openApi/cswap/v1/market/premiumIndex` | GET | Get mark price and funding rate |
| `/openApi/cswap/v1/market/ticker` | GET | Get ticker |
| `/openApi/cswap/v1/user/balance` | GET | Get balance |
| `/openApi/cswap/v1/user/positions` | GET | Get positions |
| `/openApi/cswap/v1/user/commissionRate` | GET | Get commission rate |
| `/openApi/cswap/v1/trade/order` | POST | Place order |
| `/openApi/cswap/v1/trade/cancelOrder` | DELETE | Cancel order |
| `/openApi/cswap/v1/trade/orderDetail` | GET | Query order |
| `/openApi/cswap/v1/trade/openOrders` | GET | Get open orders |
| `/openApi/cswap/v1/trade/allOpenOrders` | GET | Get all open orders |
| `/openApi/cswap/v1/trade/orderHistory` | GET | Get order history |
| `/openApi/cswap/v1/trade/allFillOrders` | GET | Get fill history |
| `/openApi/cswap/v1/trade/forceOrders` | GET | Get liquidation orders |
| `/openApi/cswap/v1/trade/leverage` | GET/POST | Get/set leverage |
| `/openApi/cswap/v1/trade/marginType` | GET/POST | Get/set margin mode |
| `/openApi/cswap/v1/trade/positionMargin` | POST | Adjust isolated margin |
| `/openApi/cswap/v1/trade/closeAllPositions` | POST | Close all positions |

---

## Reference: Order Types

| Type | Description |
|------|-------------|
| LIMIT | Limit order |
| MARKET | Market order |
| STOP_MARKET | Stop market order (stop-loss at market price) |
| TAKE_PROFIT_MARKET | Take-profit market order |
| STOP | Stop limit order (stop-loss at limit price) |
| TAKE_PROFIT | Take-profit limit order |
| TRIGGER_LIMIT | Planned commission stop-loss limit order |
| TRIGGER_MARKET | Planned commission stop-loss market order |
| TRAILING_STOP_MARKET | Trailing stop market order |

## Reference: Order Status Values

| Status | Description |
|--------|-------------|
| NEW | Order received |
| PENDING | Order pending |
| PARTIALLY_FILLED | Partially executed |
| FILLED | Fully executed |
| CANCELED | Canceled |
| CANCELLED | Canceled (alternative spelling) |
| FAILED | Failed |

## Reference: Working Type (Trigger Price)

| Value | Description |
|-------|-------------|
| MARK_PRICE | Trigger based on mark price (default) |
| CONTRACT_PRICE | Trigger based on last transaction price |
| INDEX_PRICE | Trigger based on index price |

## Reference: Time In Force

| Value | Description |
|-------|-------------|
| GTC | Good Till Cancel |
| IOC | Immediate or Cancel |
| FOK | Fill or Kill |
| PostOnly | Post Only (always maker, rejects if would match immediately) |
