# BitMart REST API Examples

This example demonstrates how to use the BitMart REST API client.

## Features Demonstrated

### Market Data API (Public)
1. **Get Single Ticker** - Real-time price for a specific symbol
2. **Get All Tickers** - All symbols' real-time prices
3. **Get Order Book** - Market depth data
4. **Get Recent Trades** - Latest executed trades
5. **Get Kline Data** - Candlestick/OHLCV data
6. **Get Trading Symbols** - Available trading pairs

### Account API (Private)
7. **Get Account Balance** - Spot wallet balances
8. **Get Wallet Balance** - Wallet balances by type (spot/margin/futures)

### Trading API (Private)
9. **Place Order** - Create limit/market orders
10. **Get Order Details** - Query specific order information
11. **Get Orders List** - Query multiple orders with filters
12. **Cancel Order** - Cancel an existing order

### Funding API (Private)
13. **Get Deposit Address** - Get deposit address for a currency
14. **Get Deposit History** - Query deposit records
15. **Get Withdrawal History** - Query withdrawal records

## Prerequisites

- Go 1.25 or higher
- BitMart API credentials (API Key, Secret Key, Memo)

## Setup

1. **Get API Credentials from BitMart**
   - Log in to BitMart: https://www.bitmart.com
   - Go to API Management
   - Create a new API key
   - Save your API Key, Secret Key, and Memo

2. **Update Credentials in main.go**
   ```go
   apiKey := "YOUR-API-KEY"
   secretKey := "YOUR-SECRET-KEY"
   memo := "YOUR-MEMO"
   ```

## Running the Example

```bash
# From the project root
go run examples/bitmart_rest/main.go

# Or build and run
go build -o bitmart_example examples/bitmart_rest/main.go
./bitmart_example
```

## Expected Output

### Public API Examples (Examples 1-6)
These will work without API credentials (read-only market data):
- ✅ Get Ticker
- ✅ Get All Tickers
- ✅ Get Order Book
- ✅ Get Recent Trades
- ✅ Get Kline Data
- ✅ Get Trading Symbols

### Private API Examples (Examples 7-15)
These require valid API credentials:
- ⚠️ Will show authentication errors in demo mode
- ✅ Will work with valid credentials

## Example Output

```
✅ BitMart client created successfully

--- Example 1: Get Ticker for BTC_USDT ---
Code: 1000, Message: OK
Symbol: BTC_USDT, Last Price: 43250.50, Volume: 1234.56
24h High: 43500.00, Low: 42800.00, Change: 1.25%

--- Example 2: Get All Tickers ---
Code: 1000, Retrieved 500 tickers
  BTC_USDT: Last=43250.50, Volume=1234.56
  ETH_USDT: Last=2250.30, Volume=5678.90
  BNB_USDT: Last=315.20, Volume=890.12

... (more examples)

✅ All examples completed!

Note: Examples 9-15 require valid API credentials and will show errors in demo mode.
Replace YOUR-API-KEY, YOUR-SECRET-KEY, and YOUR-MEMO with real credentials to test.
```

## API Rate Limits

BitMart API has rate limits:
- Public endpoints: ~10 requests/second
- Private endpoints: ~10 requests/second

The example includes delays to avoid hitting rate limits.

## Error Handling

All examples include proper error handling:
- Network errors
- API errors (returned in response)
- Authentication errors (invalid credentials)

## Security Notes

⚠️ **NEVER commit API credentials to version control!**

- Use environment variables for production
- Use `.env` files (add to `.gitignore`)
- Consider using secret management tools

## Learn More

- [BitMart API Documentation](https://developer-pro.bitmart.com/)
- [BitMart Go Client Documentation](../../exchanges/bitmart/README.md)
- [Project README](../../README.md)

## Troubleshooting

### "Authentication failed" error
- Check your API Key, Secret Key, and Memo are correct
- Ensure API key has appropriate permissions
- Verify your IP is whitelisted (if IP whitelist is enabled)

### "Invalid signature" error
- Check system time is synchronized
- Verify Secret Key is correct
- Ensure request parameters are properly formatted

### Network timeout
- Check your internet connection
- BitMart API servers may be temporarily unavailable
- Increase timeout in client configuration if needed
