# BitMart WebSocket API Examples

This directory contains comprehensive examples demonstrating the BitMart WebSocket API implementation.

## Prerequisites

- Go 1.25 or higher
- BitMart API credentials (for private channel examples)

## Setup

1. **For Public Channels** (Examples 1-5):
   ```bash
   # No credentials needed
   go run main.go
   ```

2. **For Private Channels** (Examples 6-9):
   ```bash
   # Set your API credentials as environment variables
   export BITMART_API_KEY="your_api_key"
   export BITMART_SECRET_KEY="your_secret_key"
   export BITMART_MEMO="your_memo"

   go run main.go
   ```

## Examples Overview

### Public Channels (No Authentication Required)

#### Example 1: Ticker Subscription
Subscribe to real-time ticker updates for a trading pair.
- Channel: `spot/ticker:{symbol}`
- Data: Last price, 24h volume, high, low, etc.

#### Example 2: Depth/Order Book Subscription
Subscribe to order book depth updates.
- Channel: `spot/depth{depth}:{symbol}`
- Depth levels: 5, 20, or 50
- Data: Bids and asks with price/size

#### Example 3: Trade Subscription
Subscribe to real-time trade executions.
- Channel: `spot/trade:{symbol}`
- Data: Price, size, side, timestamp

#### Example 4: Kline Subscription
Subscribe to candlestick/kline data.
- Channel: `spot/kline{step}:{symbol}`
- Steps: 1m, 3m, 5m, 15m, 30m, 45m, 1H, 2H, 3H, 4H, 1D, 1W, 1M
- Data: OHLCV (Open, High, Low, Close, Volume)

#### Example 5: Multiple Public Channels
Subscribe to multiple public channels simultaneously.
- Demonstrates concurrent channel handling
- Shows how to manage multiple data streams

### Private Channels (Authentication Required)

#### Example 6: Private Order Updates
Subscribe to real-time order status updates.
- Channel: `spot/user/order`
- Data: Order ID, status, filled amount, etc.
- Use case: Track order lifecycle

#### Example 7: Private Balance Updates
Subscribe to account balance changes.
- Channel: `spot/user/balance`
- Data: Currency, available, frozen, total
- Use case: Real-time balance monitoring

#### Example 8: Private Trade Updates
Subscribe to your trade executions.
- Channel: `spot/user/trade`
- Data: Trade ID, price, size, fee
- Use case: Track filled orders

#### Example 9: All Private Channels
Subscribe to all private channels simultaneously.
- Demonstrates full private channel integration
- Unified handling of orders, balances, and trades

## Running Examples

```bash
# Run the examples
go run main.go

# Select an example by number when prompted
Select example (1-9): 1
```

## WebSocket Features

### Connection Management
- Automatic connection establishment
- Ping/pong heartbeat (every 20 seconds)
- Graceful disconnection
- Context-based cancellation

### Public Channels
```go
// Create WebSocket client
cfg := &ws.BitMartConfig{
    WSBaseURL: types.ProductionWSServer,
}
client, err := ws.NewClientWs(ctx, cfg)

// Connect
if err := client.Connect(); err != nil {
    log.Fatal(err)
}
defer client.Close()

// Subscribe to ticker
if err := client.Public.SubscribeTicker("BTC_USDT"); err != nil {
    log.Fatal(err)
}

// Receive updates
tickerCh := client.Public.GetTickerChan()
for ticker := range tickerCh {
    fmt.Printf("Price: %s\n", ticker.LastPrice)
}
```

### Private Channels
```go
// Create client with credentials
cfg := &ws.BitMartConfig{
    APIKey:    "your_api_key",
    SecretKey: "your_secret_key",
    Memo:      "your_memo",
    WSBaseURL: types.ProductionWSServer,
}
client, err := ws.NewClientWs(ctx, cfg)

// Connect and authenticate
client.Connect()
client.Login()

// Subscribe to private channels
client.Private.SubscribeOrder()
client.Private.SubscribeBalance()
client.Private.SubscribeTrade()

// Receive updates
orderCh := client.Private.GetOrderChan()
for order := range orderCh {
    fmt.Printf("Order %s: %s\n", order.OrderID, order.Status)
}
```

## Channel Buffer

All channels have a buffer size of 100 messages. If the consumer is too slow and the buffer fills up, new messages will be dropped to prevent blocking.

## Error Handling

The examples include basic error handling. In production:
- Add reconnection logic
- Implement exponential backoff
- Log all errors properly
- Monitor connection status
- Handle subscription failures

## Security Best Practices

1. **Never hardcode credentials**
   - Use environment variables
   - Use secure credential management systems
   - Rotate API keys regularly

2. **API Key Permissions**
   - Use read-only keys for monitoring
   - Use trading keys only when necessary
   - Set IP whitelist in BitMart settings

3. **Connection Security**
   - Always use WSS (secure WebSocket)
   - Verify server certificates
   - Monitor for suspicious activity

## Troubleshooting

### Connection Issues
```
Error: Failed to connect
```
**Solution**: Check network connectivity and firewall settings

### Authentication Failures
```
Error: Failed to login
```
**Solution**: Verify API credentials and memo are correct

### No Data Received
```
Subscribed but no updates
```
**Solution**:
- Check if symbol is correct (e.g., BTC_USDT not BTC/USDT)
- Verify market is active
- Ensure you're authenticated for private channels

### Channel Full
```
Channel buffer full, dropping messages
```
**Solution**: Process messages faster or increase buffer size

## API Reference

### Public Methods

#### ClientWs
- `Connect() error` - Establish WebSocket connection
- `Close() error` - Close connection
- `Login() error` - Authenticate (for private channels)
- `IsConnected() bool` - Check connection status
- `IsAuthenticated() bool` - Check authentication status

#### Public Channels
- `SubscribeTicker(symbol string, ch ...chan *TickerEvent) error`
- `SubscribeDepth(symbol string, depth int, ch ...chan *DepthEvent) error`
- `SubscribeTrade(symbol string, ch ...chan *TradeEvent) error`
- `SubscribeKline(symbol string, step string, ch ...chan *KlineEvent) error`
- `UnsubscribeTicker(symbol string) error`
- `UnsubscribeDepth(symbol string, depth int) error`
- `UnsubscribeTrade(symbol string) error`
- `UnsubscribeKline(symbol string, step string) error`

#### Private Channels
- `SubscribeOrder(ch ...chan *OrderEvent) error`
- `SubscribeBalance(ch ...chan *BalanceEvent) error`
- `SubscribeTrade(ch ...chan *TradeEvent) error`
- `UnsubscribeOrder() error`
- `UnsubscribeBalance() error`
- `UnsubscribeTrade() error`

## Additional Resources

- [BitMart WebSocket API Documentation](https://developer-pro.bitmart.com/en/spot/#websocket-market-data)
- [Go-exc Project Documentation](../../README.md)
- [BitMart REST API Examples](../bitmart_rest/README.md)

## License

Same as the parent project.
