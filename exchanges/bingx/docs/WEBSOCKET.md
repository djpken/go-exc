# BingX WebSocket API

## WebSocket Endpoints

| Market | URL |
|--------|-----|
| Perpetual Swap (Public) | `wss://open-api-swap.bingx.com/swap-market` |
| Perpetual Swap (Private) | `wss://open-api-swap.bingx.com/swap-market?listenKey=<listenKey>` |
| Spot (Public) | `wss://open-api-ws.bingx.com/market` |

---

## Protocol

### Data Compression

All messages from the WebSocket server are GZIP-compressed. Clients must decompress each message before processing.

### Message Format

All messages use JSON format.

**Subscribe:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "reqType": "sub",
  "dataType": "<channel>"
}
```

**Unsubscribe:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "reqType": "unsub",
  "dataType": "<channel>"
}
```

**Server Response:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "code": 0,
  "msg": ""
}
```

A `code` of `0` indicates success.

### Unsubscribe

To stop receiving data from a channel, send an unsubscribe message with the same `dataType`:

**Example:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "reqType": "unsub",
  "dataType": "BTC-USDT@kline_1min"
}
```

---

## Authentication (Private Channels)

Private channels require a Listen Key. Obtain one via the REST API:

```
POST /openApi/user/auth/userDataStream
```

Connect to WebSocket with the listen key as a query parameter:
```
wss://open-api-swap.bingx.com/swap-market?listenKey=<your-listen-key>
```

> **Important:** The listen key expires after **1 hour**. Extend it regularly by calling the keep-alive endpoint to avoid subscription interruption.

After connecting with a listen key, **all event types are pushed automatically** — you do not need to subscribe to any specific channel for private data.

---

## Public Channels

### Perpetual Swap Kline/Candlestick Stream

Subscribe to real-time kline data for a perpetual swap trading pair.

**Channel Format:**
```
<symbol>@kline_<interval>
```

**Subscription Example:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "reqType": "sub",
  "dataType": "BTC-USDT@kline_1min"
}
```

**Subscription Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair with hyphen, e.g. BTC-USDT |
| interval | string | Yes | Kline interval (see supported intervals below) |

**Supported Intervals:**

| Interval | Description |
|----------|-------------|
| `1min` | 1 minute |
| `3min` | 3 minutes |
| `5min` | 5 minutes |
| `15min` | 15 minutes |
| `30min` | 30 minutes |
| `1hour` | 1 hour |
| `2hour` | 2 hours |
| `4hour` | 4 hours |
| `6hour` | 6 hours |
| `8hour` | 8 hours |
| `12hour` | 12 hours |
| `1day` | 1 day |
| `3day` | 3 days |
| `1week` | 1 week |
| `1month` | 1 month |

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `dataType` | Subscription data type, e.g. `BTC-USDT@kline_1min` |
| `data.e` | Event type |
| `data.E` | Event time |
| `data.s` | Trading pair |
| `data.K.t` | Kline start time |
| `data.K.T` | Kline end time |
| `data.K.o` | Opening price |
| `data.K.h` | Highest price |
| `data.K.l` | Lowest price |
| `data.K.c` | Closing price (latest) |
| `data.K.q` | Trading volume |
| `data.K.n` | Number of trades |
| `data.K.i` | Kline interval |
| `data.K.s` | Trading pair |

**Example Push:**
```json
{
  "dataType": "BTC-USDT@kline_1min",
  "data": {
    "e": "kline",
    "E": 1702733255486,
    "s": "BTC-USDT",
    "K": {
      "t": 1702733220000,
      "T": 1702733279999,
      "o": "42000.00",
      "h": "42050.00",
      "l": "41990.00",
      "c": "42020.00",
      "q": "12.345",
      "n": 87,
      "i": "1min",
      "s": "BTC-USDT"
    }
  }
}
```

---

### Legacy Kline Stream (market.kline format)

An alternative kline subscription format used by some older endpoints.

**Channel Format:**
```
market.kline.<symbol>.<klineType>
```

**Example:** `market.kline.BTC-USDT.1min`

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `code` | 0 = normal, 1 = error |
| `dataType` | Subscribed data type |
| `data.klineInfosVo` | Kline data array |
| `data.klineInfosVo[].open` | Opening price |
| `data.klineInfosVo[].close` | Closing price |
| `data.klineInfosVo[].high` | High price |
| `data.klineInfosVo[].low` | Low price |
| `data.klineInfosVo[].volume` | Volume |
| `data.klineInfosVo[].time` | Kline timestamp (ms) |
| `data.klineInfosVo[].statDate` | Kline date string |

---

### Market Depth Stream

Subscribe to real-time order book snapshots for a trading pair.

**Channel Format (new):**
```
<symbol>@depth<level>
```
Example: `BTC-USDT@depth50`

**Channel Format (legacy):**
```
market.depth.<symbol>.<step>.<level>
```
Example: `market.depth.BTC-USDT.step0.level5`

Snapshots are pushed once per second.

**Subscription Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |
| level | string | Yes | Depth level: `level5`, `level10`, `level20`, `level50`, `level100` |

**Step Options (price aggregation):**

| Value | Description |
|-------|-------------|
| `step0` | No price merging (raw data) |
| `step1` | Merge at 10x minimum price precision |
| `step2` | Merge at 100x minimum price precision |
| `step3` | Merge at 1,000x minimum price precision |
| `step4` | Merge at 10,000x minimum price precision |
| `step5` | Merge at 100,000x minimum price precision |

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `code` | 0 = normal, 1 = error |
| `dataType` | Subscribed data type |
| `data.asks` | Ask orders: `[[price, quantity], ...]` |
| `data.bids` | Bid orders: `[[price, quantity], ...]` |

---

### Latest Trade Stream

Subscribe to real-time trade data for a trading pair.

**Channel Format:**
```
<symbol>@trade
```

**Subscription Example:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "reqType": "sub",
  "dataType": "BTC-USDT@trade"
}
```

**Subscription Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `dataType` | Subscribed data type |
| `data` | Trade data array |
| `data[].T` | Trade timestamp |
| `data[].s` | Trading pair |
| `data[].m` | Whether buyer is maker |
| `data[].p` | Trade price |
| `data[].q` | Trade quantity |

---

### 24-Hour Ticker Stream

Subscribe to 24-hour rolling price change statistics. Updated every 1000ms.

**Channel Format:**
```
<symbol>@ticker
```

**Subscription Example:**
```json
{
  "id": "975f7385-7f28-4ef1-93af-df01cb9ebb53",
  "reqType": "sub",
  "dataType": "BTC-USDT@ticker"
}
```

**Subscription Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `e` | Event type |
| `E` | Event time |
| `s` | Trading pair |
| `p` | Price change amount |
| `P` | Price change percentage |
| `o` | Opening price |
| `h` | Highest price |
| `l` | Lowest price |
| `c` | Latest transaction price |
| `v` | Trading volume |
| `q` | Trading amount (quote asset) |
| `O` | Statistics start time |
| `C` | Statistics end time |
| `B` | Best bid price |
| `b` | Best bid quantity |
| `A` | Best ask price |
| `a` | Best ask quantity |
| `n` | Number of trades |

---

### Book Ticker Stream

Subscribe to real-time best bid and ask price updates.

**Channel Format:**
```
<symbol>@bookTicker
```

**Subscription Example:**
```json
{
  "id": "24dd0e35-56a4-4f7a-af8a-394c7060909c",
  "reqType": "sub",
  "dataType": "BTC-USDT@bookTicker"
}
```

**Subscription Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| symbol | string | Yes | Trading pair, e.g. BTC-USDT |

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `code` | 0 = normal, 1 = error |
| `dataType` | Subscribed data type, e.g. `BTC-USDT@bookTicker` |
| `data.e` | Event type |
| `data.E` | Event time |
| `data.T` | Transaction time |
| `data.s` | Trading pair |
| `data.u` | Update ID |
| `data.b` | Best bid price |
| `data.B` | Best bid quantity |
| `data.a` | Best ask price |
| `data.A` | Best ask quantity |

**Example Push:**
```json
{
  "code": 0,
  "dataType": "BTC-USDT@bookTicker",
  "data": {
    "e": "bookTicker",
    "u": 123456789,
    "E": 1702733255486,
    "T": 1702733255480,
    "s": "BTC-USDT",
    "b": "42000.00",
    "B": "1.5",
    "a": "42001.00",
    "A": "0.8"
  }
}
```

---

## Private Channels

Private channels are pushed automatically after connecting with a valid `listenKey`. No subscription message is needed.

### Account Balance Update (`ACCOUNT_UPDATE`)

Pushed when account balance changes due to deposits, withdrawals, order activity, funding fees, etc.

> **Note:** No subscription required. Connect with a listenKey and all account events are pushed automatically.

**Event Trigger Reasons (`m` field):**

| Reason | Description |
|--------|-------------|
| DEPOSIT | Deposit |
| WITHDRAW | Withdrawal |
| ORDER | Order activity |
| FUNDING_FEE | Funding fee settlement |
| WITHDRAW_REJECT | Withdrawal rejected |
| ADJUSTMENT | Manual adjustment |
| INSURANCE_CLEAR | Insurance fund clearance |
| ADMIN_DEPOSIT | Admin deposit |
| ADMIN_WITHDRAW | Admin withdrawal |
| MARGIN_TRANSFER | Margin transfer |
| MARGIN_TYPE_CHANGE | Margin type changed |
| ASSET_TRANSFER | Asset transfer |
| OPTIONS_PREMIUM_FEE | Options premium fee |
| OPTIONS_SETTLE_PROFIT | Options settlement profit |
| AUTO_EXCHANGE | Auto exchange |

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `e` | Event type: `ACCOUNT_UPDATE` |
| `E` | Event time |
| `T` | Matching engine time |
| `a` | Account data |
| `a.m` | Event trigger reason |
| `a.B` | Balance array |
| `a.B[].a` | Asset name |
| `a.B[].wb` | Wallet balance |
| `a.B[].cw` | Cross wallet balance (excluding isolated margin) |
| `a.B[].bc` | Balance change (excluding PnL and trading fees) |

**Subscription Example (legacy format):**
```json
{
  "id": "gdfg2311-d0f6-4a70-8d5a-043e4c741b40",
  "dataType": "ACCOUNT_UPDATE"
}
```

---

### Order Update (`ORDER_TRADE_UPDATE`)

Pushed when a new order is created, an order is filled, or an order status changes.

> **Note:** No subscription required. Connect with a listenKey and all order events are pushed automatically.

**Execution Types (`x` field):**

| Value | Description |
|-------|-------------|
| NEW | New order |
| CANCELED | Order canceled |
| CALCULATED | ADL or liquidation |
| EXPIRED | Order expired |
| TRADE | Trade executed |

**Order Status (`X` field):**

| Value | Description |
|-------|-------------|
| NEW | New |
| PARTIALLY_FILLED | Partially filled |
| FILLED | Fully filled |
| CANCELED | Canceled |
| EXPIRED | Expired |

**Push Data Response:**

| Field | Description |
|-------|-------------|
| `e` | Event type: `ORDER_TRADE_UPDATE` |
| `E` | Event time |
| `o` | Order object |
| `o.s` | Trading pair |
| `o.c` | Client order ID |
| `o.i` | Order ID |
| `o.S` | Order side: `BUY` or `SELL` |
| `o.o` | Order type: `MARKET`, `LIMIT`, `STOP`, `TAKE_PROFIT`, `LIQUIDATION` |
| `o.q` | Order quantity |
| `o.p` | Order price |
| `o.sp` | Trigger price (stop price) |
| `o.ap` | Average fill price |
| `o.x` | Execution type (this event) |
| `o.X` | Current order status |
| `o.l` | Last filled quantity |
| `o.z` | Total filled quantity |
| `o.L` | Last filled price |
| `o.n` | Fee amount |
| `o.N` | Fee asset |
| `o.T` | Trade time |
| `o.t` | Trade ID |
| `o.b` | Bid notional |
| `o.a` | Ask notional |
| `o.m` | Is maker |
| `o.R` | Is reduce-only |
| `o.wt` | Working type (trigger price type) |
| `o.ot` | Original order type |
| `o.ps` | Position side: `LONG`, `SHORT`, or `BOTH` |
| `o.rp` | Realized PnL of this trade |

---

### Position and Balance Update (Combined Push)

When connected with a listen key, the server may also push combined position and balance updates via the `ACCOUNT_CONFIG_UPDATE` and related events.

---

## WebSocket Channel Summary

### Perpetual Swap Public Channels

| Channel | Description | Update Frequency |
|---------|-------------|-----------------|
| `<symbol>@kline_<interval>` | Kline/candlestick data | Per new kline |
| `<symbol>@depth<level>` | Order book depth snapshot | 1 second |
| `<symbol>@trade` | Latest trades | Real-time |
| `<symbol>@ticker` | 24h rolling ticker stats | 1 second |
| `<symbol>@bookTicker` | Best bid/ask prices | Real-time |

### Perpetual Swap Private Channels (Auto-pushed with listenKey)

| Event Type | Description |
|------------|-------------|
| `ACCOUNT_UPDATE` | Balance changes |
| `ORDER_TRADE_UPDATE` | Order status/fill updates |
| `ACCOUNT_CONFIG_UPDATE` | Account configuration changes |

### Legacy Channel Formats (Older API)

| Channel | Description |
|---------|-------------|
| `market.kline.<symbol>.<interval>` | Kline data |
| `market.depth.<symbol>.<step>.<level>` | Order book depth |
| `ACCOUNT_UPDATE` | Account balance push |

---

## Kline Interval Reference

When subscribing to klines via the `@kline_<interval>` format:

| Interval Value | Description |
|---------------|-------------|
| `1min` | 1 minute |
| `3min` | 3 minutes |
| `5min` | 5 minutes |
| `15min` | 15 minutes |
| `30min` | 30 minutes |
| `1hour` | 1 hour |
| `2hour` | 2 hours |
| `4hour` | 4 hours |
| `6hour` | 6 hours |
| `8hour` | 8 hours |
| `12hour` | 12 hours |
| `1day` | 1 day |
| `3day` | 3 days |
| `1week` | 1 week |
| `1month` | 1 month |

---

## Implementation Notes

1. **Heartbeat:** Implement ping/pong to keep the connection alive. The server will close idle connections.

2. **Reconnection:** Implement automatic reconnection with exponential backoff on connection failure.

3. **Message ID:** Use unique UUIDs for message IDs to correlate responses with requests.

4. **Decompression:** All messages must be GZIP-decompressed before parsing as JSON.

5. **Listen Key Renewal:** Renew the listen key every 30 minutes using the REST endpoint to prevent expiration.

6. **Symbol Format:** Trading pair symbols must use a hyphen separator (e.g., `BTC-USDT`). Always use uppercase letters.

7. **Error Handling:** Check the `code` field in responses; `0` means success, `1` means error.

---

## Go Implementation Example

```go
package main

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "io"
    "log"

    "github.com/gorilla/websocket"
)

func decompressGzip(data []byte) ([]byte, error) {
    r, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    defer r.Close()
    return io.ReadAll(r)
}

func subscribeKline(conn *websocket.Conn, symbol, interval string) error {
    msg := map[string]string{
        "id":       "unique-uuid-here",
        "reqType":  "sub",
        "dataType": symbol + "@kline_" + interval,
    }
    data, _ := json.Marshal(msg)
    return conn.WriteMessage(websocket.TextMessage, data)
}

func main() {
    conn, _, err := websocket.DefaultDialer.Dial(
        "wss://open-api-swap.bingx.com/swap-market",
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    subscribeKline(conn, "BTC-USDT", "1min")

    for {
        _, raw, err := conn.ReadMessage()
        if err != nil {
            log.Fatal(err)
        }
        decompressed, err := decompressGzip(raw)
        if err != nil {
            log.Printf("decompress error: %v", err)
            continue
        }
        log.Printf("received: %s", decompressed)
    }
}
```
