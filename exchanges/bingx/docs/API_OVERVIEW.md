# BingX API Overview

## Base URLs

| Environment | URL |
|-------------|-----|
| REST API (Primary) | `https://open-api.bingx.com` |
| REST API (Backup) | `https://open-api.bingx.io` |
| WebSocket (Perpetual Swap) | `wss://open-api-swap.bingx.com/swap-market` |
| WebSocket (Spot) | `wss://open-api-ws.bingx.com/market` |

> **Note:** The backup domain `open-api.bingx.io` has a total rate limit of 60 requests/minute and is only available when the primary domain is experiencing issues.

---

## Authentication

BingX REST API authentication uses HMAC SHA256 signatures. Every private endpoint requires:

1. An API key passed via the `X-BX-APIKEY` request header.
2. A `signature` parameter computed from your request parameters and secret key.
3. A `timestamp` parameter (Unix timestamp in milliseconds).

### Request Headers

| Header | Description |
|--------|-------------|
| `X-BX-APIKEY` | Your API key |
| `Content-Type` | `application/json` (for POST requests) |

### Signature Algorithm

The signature is an HMAC SHA256 hex digest of the query string built from your request parameters.

**Step 1: Build the query string**

- For GET requests: Concatenate all parameters without sorting.
- For POST requests: Sort parameters alphabetically (a-z) before concatenating.

Example query string:
```
recvWindow=0&symbol=BTC-USDT&timestamp=1696751141337
```

**Step 2: Sign with your secret key**

```bash
echo -n "recvWindow=0&symbol=BTC-USDT&timestamp=1696751141337" | \
  openssl dgst -sha256 -hmac "<your-secret-key>" -hex
```

**Step 3: Append the signature**

Append `&signature=<computed-signature>` to your request.

### GET Request Example

```bash
curl -H 'X-BX-APIKEY: <your-api-key>' \
  'https://open-api.bingx.com/openApi/swap/v2/user/positions?recvWindow=0&symbol=BTC-USDT&timestamp=1696751141337&signature=<signature>'
```

### POST Request Example

```bash
curl --location 'https://open-api.bingx.com/openApi/subAccount/v1/create' \
  --header 'Content-Type: application/json' \
  --data '{"recvWindow":0,"subAccountString":"abc12345","timestamp":1696751141337,"signature":"<signature>"}'
```

### URL Encoding

Some query string values (e.g. JSON arrays) require URL-encoding of the **value** only (not the key, not the entire string). The `timestamp` value does NOT need URL-encoding.

---

## Common Request Parameters

Every authenticated endpoint accepts these standard parameters:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window in milliseconds (default: 5000) |
| signature | string | Yes | HMAC SHA256 signature |

---

## Response Format

All REST API responses return JSON:

```json
{
  "code": 0,
  "msg": "",
  "data": { ... },
  "timestamp": 1702288510557
}
```

| Field | Description |
|-------|-------------|
| code | Error code. `0` = success; non-zero = error |
| msg | Error message (empty on success) |
| data | Response payload |
| timestamp | Server timestamp |

---

## Rate Limits

| Condition | Behavior |
|-----------|----------|
| Rate limit exceeded | HTTP 429 returned |
| Continued requests after 429 | HTTP 418 (IP auto-banned) |
| Backup domain total limit | 60 requests/minute |

---

## Error Codes

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Bad Request - invalid format |
| 401 | Unauthorized - invalid API key |
| 403 | Forbidden - no access to resource |
| 404 | Not Found |
| 418 | IP auto-banned for continued requests after 429 |
| 429 | Too Many Requests - rate limit exceeded |
| 500 | Internal Server Error |
| 504 | Gateway Timeout - request submitted but response unknown; request may or may not have succeeded |

### Business Error Codes

| Code | Description |
|------|-------------|
| 100001 | Signature verification failed |
| 100004 | Permission denied (API key missing required permission) |
| 100400 | Invalid parameter |
| 100412 | Null signature |
| 100413 | Incorrect API key |
| 100419 | IP not in whitelist |
| 100421 | Null timestamp or timestamp mismatch |
| 100500 | Internal system error |
| 101204 | Insufficient margin |
| 101209 | Maximum position value for this leverage exceeded |
| 101212 | Failed - pending orders exist for this trading pair (cancel them first) |
| 101215 | Post-Only order would match immediately (canceled) |
| 101414 | Maximum leverage for trading pair exceeded |
| 101415 | Trading pair suspended from opening new positions |
| 101460 | Order price below estimated liquidation price of long position |
| 101500 | RPC timeout |
| 101514 | Temporarily suspended from opening positions |
| 109201 | Duplicate order ID submitted within 1 second |
| 80012 | Service unavailable |
| 80013 | Maximum entrusted orders limit reached |
| 80014 | Invalid parameter |
| 80016 | Order does not exist |
| 80017 | Position does not exist |
| 80018 | Order already filled |
| 80019 | Order is being processed (use allOrders to retrieve details) |
| 80020 | Risk forbidden |

---

## API Key Permissions

When creating an API key, set the appropriate permissions:

| Permission Value | Description |
|-----------------|-------------|
| 1 | Read - read-only access to account data |
| 2 | Trade - place and cancel orders |
| 4 | Withdraw - initiate withdrawals |
| 5 | Withdraw (sub-account) |

> **Security:** Never share your API key or secret key. If accidentally leaked, delete it immediately and create a new one. Store credentials securely and never commit them to source control.

---

## Symbol Formats

| Market | Format | Example |
|--------|--------|---------|
| Spot | `BASE-QUOTE` | `BTC-USDT` |
| Perpetual Swap | `BASE-QUOTE` | `BTC-USDT` |
| Coin-Margined Swap | `BASE-QUOTE` | `BTC-USDT` |

> **Important:** Always use uppercase letters for trading pair symbols.

---

## Listen Key (WebSocket Authentication)

Private WebSocket channels require a Listen Key for authentication. The listen key is valid for **1 hour** and must be renewed regularly.

### Generate Listen Key

```
POST /openApi/user/auth/userDataStream
```

No request body parameters required (only standard `timestamp` and `signature`).

**Response:**
```json
{
  "listenKey": "a8ea75681542e66f1a50a1616dd06ed77dab61baa0c296bca03a9b13ee5f2dd7"
}
```

### Extend Listen Key (Keep Alive)

```
PUT /openApi/user/auth/userDataStream
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| listenKey | string | Yes | The listen key to extend |

### Close Listen Key

```
DELETE /openApi/user/auth/userDataStream
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| listenKey | string | Yes | The listen key to close |

**Response Status:**
- HTTP 200 - Success
- HTTP 204 - No request parameters
- HTTP 404 - Listen key not found

### Connect to Private WebSocket

```
wss://open-api-swap.bingx.com/swap-market?listenKey=<your-listen-key>
```

> **Important:** Renew the listen key regularly (e.g., every 30 minutes) to keep subscriptions active.

---

## WebSocket Protocol

### Data Compression

All WebSocket server responses are compressed using GZIP. Clients must decompress responses before processing.

### Message Format

**Subscribe:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "reqType": "sub",
  "dataType": "BTC-USDT@kline_1min"
}
```

**Unsubscribe:**
```json
{
  "id": "e745cd6d-d0f6-4a70-8d5a-043e4c741b40",
  "reqType": "unsub",
  "dataType": "BTC-USDT@kline_1min"
}
```

---

## Account Types

BingX supports both spot and perpetual swap (futures) trading through separate API namespaces:

| Namespace | Description |
|-----------|-------------|
| `/openApi/spot/` | Spot trading endpoints |
| `/openApi/swap/` | Perpetual swap (USDT-margined futures) endpoints |
| `/openApi/cswap/` | Coin-margined perpetual swap endpoints |
| `/openApi/contract/` | Standard contract endpoints |
| `/openApi/wallets/` | Wallet and transfer endpoints |
| `/openApi/subAccount/` | Sub-account management endpoints |
| `/openApi/account/` | Account-level endpoints |

---

## Historical Data API

| Endpoint | Description |
|----------|-------------|
| `GET /openApi/market/his/v1/kline` | Historical kline data |
| `GET /openApi/market/his/v1/trade` | Historical trade data |

---

## Fund Transfer

Universal asset transfer between account types:

```
POST /openApi/api/v3/post/asset/transfer
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| asset | string | Yes | Asset/currency name (e.g., USDT) |
| amount | float | Yes | Transfer amount |
| type | string | Yes | Transfer type (e.g., FUND_SFUTURES: Fund to Standard Futures) |
| timestamp | int64 | Yes | Request timestamp in milliseconds |
| recvWindow | int64 | No | Request validity window |

**Query Transfer Records:**
```
GET /openApi/api/v3/get/asset/transfer
```

---

## Deposit and Withdrawal

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/openApi/wallets/v1/capital/deposit/address` | GET | Get deposit address |
| `/openApi/wallets/v1/capital/withdraw/apply` | POST | Submit withdrawal |
| `/openApi/api/v3/capital/deposit/hisrec` | GET | Deposit history (multi-network) |
| `/openApi/api/v3/capital/withdraw/history` | GET | Withdrawal history (multi-network) |
| `/openApi/wallets/v1/capital/config/getall` | GET | Get all supported coins config |
