# BingX Spot Trading REST API

## Base URL

```
https://open-api.bingx.com
```

## Symbol Format

Spot trading pairs use a hyphen separator: `BTC-USDT`, `ETH-USDT`, `BNB-USDT`, etc. Always use uppercase letters.

## Authentication

All private endpoints require `X-BX-APIKEY` header and `signature` parameter. See `API_OVERVIEW.md` for authentication details.

---

## Market Data Endpoints (Public)

### Get Trading Symbols

```
GET /openApi/spot/v1/common/symbols
```

Get all supported spot trading pairs and their trading rules.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Filter by specific trading pair (e.g. BTC-USDT) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair symbol |
| minQty | string | Minimum order quantity |
| maxQty | string | Maximum order quantity |
| minNotional | string | Minimum order value (quote asset) |
| maxNotional | string | Maximum order value (quote asset) |
| status | int | Symbol status (1 = online) |
| tickSize | string | Minimum price increment |
| stepSize | string | Minimum quantity increment |

---

### Get All Spot Prices

```
GET /openApi/spot/v1/common/prices
```

Get latest prices for all spot trading pairs.

---

### Get Order Book Depth

```
GET /openApi/spot/v1/market/depth
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| depth | int | No | Number of levels (default: 5, max: 20) |
| type | string | No | Price aggregation step: step0 (no merge), step1–step5 |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| bids | array | Bid orders: `[[price, quantity], ...]` sorted by price descending |
| asks | array | Ask orders: `[[price, quantity], ...]` sorted by price ascending |

**Example Response:**
```json
{
  "code": 0,
  "data": {
    "bids": [["42000.50", "0.5"], ["41999.00", "1.2"]],
    "asks": [["42001.00", "0.8"], ["42002.00", "2.0"]]
  }
}
```

---

### Get Kline/Candlestick Data (v2)

```
GET /openApi/spot/v2/market/kline
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
| volume | string | Base asset volume |
| time | int64 | Kline open time (milliseconds) |

---

### Get Kline Data (v1)

```
GET /openApi/spot/v1/market/kline
```

Same parameters as v2 above.

---

### Get 24hr Ticker Statistics

```
GET /openApi/spot/v1/ticker/24hr
```

Get 24-hour rolling window price statistics.

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
| openPrice | string | Opening price (24h window) |
| highPrice | string | Highest price in 24h |
| lowPrice | string | Lowest price in 24h |
| lastPrice | string | Latest transaction price |
| volume | string | Base asset volume in 24h |
| quoteVolume | string | Quote asset volume in 24h |
| openTime | int64 | Statistics window open time (ms) |
| closeTime | int64 | Statistics window close time (ms) |
| bidPrice | string | Best bid price |
| bidQty | string | Best bid quantity |
| askPrice | string | Best ask price |
| askQty | string | Best ask quantity |

---

### Get Latest Price (v2)

```
GET /openApi/spot/v2/ticker/price
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
| price | string | Latest price |
| time | int64 | Timestamp in milliseconds |

---

### Get Best Order Book Price (Book Ticker)

```
GET /openApi/spot/v1/ticker/bookTicker
```

Get best bid and ask prices for a trading pair.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| bidPrice | string | Best bid price |
| bidQty | string | Best bid quantity |
| askPrice | string | Best ask price |
| askQty | string | Best ask quantity |

---

### Get Recent Trades

```
GET /openApi/spot/v1/market/trades
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| limit | int | No | Number of results (default: 100, max: 500) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

## Account Endpoints (Private)

### Get Account Balance

```
GET /openApi/spot/v1/account/balance
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| asset | string | Asset/currency name |
| free | string | Available balance |
| locked | string | Locked/frozen balance |

**Example Response:**
```json
{
  "code": 0,
  "data": {
    "balances": [
      {"asset": "USDT", "free": "1000.00", "locked": "50.00"},
      {"asset": "BTC", "free": "0.5", "locked": "0.0"}
    ]
  }
}
```

---

### Get Trading Commission Rate

```
GET /openApi/spot/v1/user/commissionRate
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
| makerCommissionRate | string | Maker fee rate (e.g. "0.001" = 0.1%) |
| takerCommissionRate | string | Taker fee rate |

---

## Trading Endpoints (Private)

### Place Order

```
POST /openApi/spot/v1/trade/order
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| side | string | Yes | Order direction: `BUY` or `SELL` |
| type | string | Yes | Order type: `MARKET`, `LIMIT`, or `LIMIT_MAKER` |
| quantity | float | No | Base asset quantity. For MARKET BUY, use `quoteOrderQty` instead |
| quoteOrderQty | float | No | Quote asset amount for MARKET BUY orders |
| price | float | No | Limit price (required for LIMIT orders) |
| newClientOrderId | string | No | Custom order ID (1–40 characters; must be unique) |
| timeInForce | string | No | `GTC` (default), `IOC`, `FOK`, `POC` |
| stopPrice | float | No | Trigger price for stop orders |
| recvWindow | int64 | No | Request validity window in milliseconds |
| timestamp | int64 | Yes | Request timestamp in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| orderId | int64 | System order ID |
| transactTime | int64 | Transaction timestamp in milliseconds |
| price | string | Order price |
| origQty | string | Original quantity |
| executedQty | string | Executed quantity |
| cummulativeQuoteQty | string | Total quote asset transacted |
| status | string | Order status |
| type | string | Order type |
| side | string | Order side |

---

### Test Place Order (No Execution)

```
POST /openApi/spot/v1/trade/order/test
```

Same parameters as the regular place order endpoint. Validates parameters without creating a real order.

---

### Cancel and Replace Order

```
POST /openApi/spot/v1/trade/order/cancelReplace
```

Cancels an existing order and places a new one in a single operation.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| cancelOrderId | int64 | No | Order ID to cancel |
| cancelClientOrderId | string | No | Custom order ID to cancel |
| side | string | Yes | BUY or SELL |
| type | string | Yes | Order type |
| quantity | float | No | New order quantity |
| price | float | No | New order price |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Query Order

```
GET /openApi/spot/v1/trade/query
```

Get details of a single order.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderId | int64 | No | System order ID |
| clientOrderID | string | No | Custom order ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

> Either `orderId` or `clientOrderID` must be provided.

---

### Get Open Orders

```
GET /openApi/spot/v1/trade/openOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| orderId | int64 | No | Start from this order ID |
| pageIndex | int | No | Page number (starts at 1) |
| pageSize | int | No | Results per page (max: 100) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Order History

```
GET /openApi/spot/v1/trade/historyOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | No | Trading pair |
| orderId | int64 | No | Return orders from this ID onwards |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| pageIndex | int | No | Page number |
| pageSize | int | No | Results per page |
| status | string | No | Filter by order status |
| type | string | No | Filter by order type |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Cancel Order

```
POST /openApi/spot/v1/trade/cancel
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderId | int64 | No | System order ID to cancel |
| clientOrderID | string | No | Custom order ID to cancel |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| orderId | int64 | Canceled order ID |
| price | string | Order price |
| origQty | string | Original quantity |
| executedQty | string | Quantity already executed |
| cummulativeQuoteQty | string | Total quote asset transacted |
| status | string | Order status after cancellation |
| type | string | Order type |
| side | string | Order side |

---

### Cancel All Open Orders for Symbol

```
POST /openApi/spot/v1/trade/cancelOpenOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Cancel a Batch of Orders

```
POST /openApi/spot/v1/trade/cancelOrders
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderIds | string | No | Comma-separated list of order IDs |
| clientOrderIds | string | No | Comma-separated list of custom order IDs |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Cancel All After Timer

```
POST /openApi/spot/v1/trade/cancelAllAfter
```

Schedule a cancellation of all open orders after a specified timeout. Useful as a dead-man's switch.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| type | string | Yes | `ACTIVATE` to set timer, `DEACTIVATE` to cancel it |
| timeOut | int64 | Yes | Timeout in milliseconds after which all orders are canceled |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Batch Place Orders

```
POST /openApi/spot/v1/trade/batchOrders
```

Place multiple orders in a single request.

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| data | string | Yes | URL-encoded JSON array of order objects |
| sync | bool | No | Process orders synchronously (default: false) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

Each order object in `data` uses the same fields as a single order placement.

---

### Get Trade History (My Trades)

```
GET /openApi/spot/v1/trade/myTrades
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderId | int64 | No | Filter trades for this order |
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| fromId | int64 | No | Trade ID to start from |
| limit | int | No | Number of results (default: 500, max: 1000) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| symbol | string | Trading pair |
| id | int64 | Trade ID |
| orderId | int64 | Order ID |
| price | string | Trade price |
| qty | string | Trade quantity (base asset) |
| quoteQty | string | Trade amount (quote asset) |
| commission | string | Commission charged |
| commissionAsset | string | Asset used for commission |
| time | int64 | Trade timestamp in milliseconds |
| isBuyer | bool | True if buyer side |
| isMaker | bool | True if maker |

---

## OCO (One-Cancels-the-Other) Orders

An OCO order pairs a limit order with a stop order. When one fills, the other is automatically canceled.

### Place OCO Order

```
POST /openApi/spot/v1/oco/order
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| listClientOrderId | string | No | Custom OCO order list ID |
| side | string | Yes | `BUY` or `SELL` |
| quantity | float | Yes | Order quantity |
| limitClientOrderId | string | No | Custom ID for the limit leg |
| price | float | Yes | Limit leg price |
| stopClientOrderId | string | No | Custom ID for the stop leg |
| stopPrice | float | Yes | Stop leg trigger price |
| stopLimitPrice | float | No | Stop limit price (if stop leg is STOP_LIMIT) |
| stopLimitTimeInForce | string | No | `GTC`, `FOK`, or `IOC` |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Cancel OCO Order

```
POST /openApi/spot/v1/oco/cancel
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair |
| orderListId | int64 | No | OCO order list ID |
| listClientOrderId | string | No | Custom OCO order list ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get OCO Order

```
GET /openApi/spot/v1/oco/orderList
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| orderListId | int64 | No | OCO order list ID |
| listClientOrderId | string | No | Custom OCO list ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get Open OCO Orders

```
GET /openApi/spot/v1/oco/openOrderList
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

### Get OCO Order History

```
GET /openApi/spot/v1/oco/historyOrderList
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| startTime | int64 | No | Start timestamp in milliseconds |
| endTime | int64 | No | End timestamp in milliseconds |
| limit | int | No | Number of results |
| fromId | int64 | No | Start from this order list ID |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds |

---

## Reference: Order Status Values

| Status | Description |
|--------|-------------|
| NEW | Order received but not yet processed |
| PENDING | Order pending processing |
| PARTIALLY_FILLED | Order partially executed |
| FILLED | Order fully executed |
| CANCELED | Order canceled |
| FAILED | Order failed |

## Reference: Order Types

| Type | Description |
|------|-------------|
| LIMIT | Limit order - execute at specified price or better |
| MARKET | Market order - execute at best available price immediately |
| LIMIT_MAKER | Limit order that is always a maker (rejected if it would match immediately) |

## Reference: Time In Force (TIF)

| Value | Description |
|-------|-------------|
| GTC | Good Till Cancel - remains active until filled or manually canceled |
| IOC | Immediate or Cancel - fills as much as possible immediately, cancels remainder |
| FOK | Fill or Kill - must fill completely immediately or the entire order is canceled |
| POC | Post Only Cancel - ensures the order is always a maker |
