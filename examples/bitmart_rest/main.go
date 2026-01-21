package main

import (
	"context"
	"log"

	"github.com/djpken/go-exc/exchanges/bitmart"
	accountreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/account"
	fundingreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/funding"
	marketreq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/market"
	tradereq "github.com/djpken/go-exc/exchanges/bitmart/requests/rest/trade"
)

func main() {
	// Initialize API credentials
	apiKey := "f42ce865e6123a77b507659452384a2b48165991"
	secretKey := "72cca2bc92ee262f7e5398e3344b15638d43224a6fc5404567fc884f655a4b73"
	memo := ".Vm3djpcl3gj94" // BitMart requires a memo

	ctx := context.Background()

	// Create BitMart client
	client, err := bitmart.NewClient(ctx, apiKey, secretKey, memo)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	log.Println("✅ BitMart client created successfully")

	// ===========================================
	// Market Data API Examples
	// ===========================================

	// Example 1: Get Single Ticker
	log.Println("\n--- Example 1: Get Ticker for BTC_USDT ---")
	ticker, err := client.Rest.Market.GetTicker(marketreq.GetTickerRequest{
		Symbol: "BTC_USDT",
	})
	if err != nil {
		log.Printf("Error getting ticker: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", ticker.Code, ticker.Message)
		log.Printf("Symbol: %s, Last Price: %s, Volume: %s",
			ticker.Data.Symbol, ticker.Data.LastPrice, ticker.Data.BaseVolume)
		log.Printf("24h High: %s, Low: %s, Change: %s%%",
			ticker.Data.HighPrice, ticker.Data.LowPrice, ticker.Data.PercentChange)
	}

	// Example 2: Get All Tickers
	log.Println("\n--- Example 2: Get All Tickers ---")
	tickers, err := client.Rest.Market.GetTickers()
	if err != nil {
		log.Printf("Error getting tickers: %v", err)
	} else {
		log.Printf("Code: %d, Retrieved %d tickers", tickers.Code, len(tickers.Data))
		// Print first 3 tickers as example
		// Format: [symbol, last, v_24h, qv_24h, open_24h, high_24h, low_24h, fluctuation, bid_px, bid_sz, ask_px, ask_sz, ts]
		for i, t := range tickers.Data {
			if i >= 3 {
				break
			}
			if len(t) >= 3 {
				log.Printf("  %v: Last=%v, Volume=%v", t[0], t[1], t[2])
			}
		}
	}

	// Example 3: Get Order Book
	log.Println("\n--- Example 3: Get Order Book for BTC_USDT ---")
	orderBook, err := client.Rest.Market.GetOrderBook(marketreq.GetOrderBookRequest{
		Symbol: "BTC_USDT",
		Depth:  20, // Top 20 levels
	})
	if err != nil {
		log.Printf("Error getting order book: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", orderBook.Code, orderBook.Message)
		log.Printf("Timestamp: %d", orderBook.Data.Timestamp)
		log.Printf("Bids: %d levels, Asks: %d levels",
			len(orderBook.Data.Bids), len(orderBook.Data.Asks))
		if len(orderBook.Data.Bids) > 0 && len(orderBook.Data.Bids[0]) >= 2 {
			log.Printf("Best Bid: %s @ %s", orderBook.Data.Bids[0][1], orderBook.Data.Bids[0][0])
		}
		if len(orderBook.Data.Asks) > 0 && len(orderBook.Data.Asks[0]) >= 2 {
			log.Printf("Best Ask: %s @ %s", orderBook.Data.Asks[0][1], orderBook.Data.Asks[0][0])
		}
	}

	// Example 4: Get Recent Trades
	log.Println("\n--- Example 4: Get Recent Trades for BTC_USDT ---")
	trades, err := client.Rest.Market.GetTrades(marketreq.GetTradesRequest{
		Symbol: "BTC_USDT",
		Limit:  10,
	})
	if err != nil {
		log.Printf("Error getting trades: %v", err)
	} else {
		log.Printf("Code: %d, Retrieved %d trades", trades.Code, len(trades.Data))
		// Format: [symbol, timestamp, price, size, side]
		for i, trade := range trades.Data {
			if i >= 3 {
				break
			}
			if len(trade) >= 5 {
				log.Printf("  Trade: %v @ %v, Side: %v", trade[3], trade[2], trade[4])
			}
		}
	}

	// Example 5: Get Kline/Candlestick Data
	log.Println("\n--- Example 5: Get Kline Data for BTC_USDT ---")
	klines, err := client.Rest.Market.GetKlines(marketreq.GetKlineRequest{
		Symbol: "BTC_USDT",
		Step:   60, // 1 hour
		Limit:  10,
	})
	if err != nil {
		log.Printf("Error getting klines: %v", err)
	} else {
		log.Printf("Code: %d, Retrieved %d klines", klines.Code, len(klines.Data))
		// Format: [timestamp, open, high, low, close, volume, quote_volume]
		for i, k := range klines.Data {
			if i >= 3 {
				break
			}
			if len(k) >= 7 {
				log.Printf("  Kline: O=%v H=%v L=%v C=%v V=%v",
					k[1], k[2], k[3], k[4], k[5])
			}
		}
	}

	// Example 6: Get Trading Symbols
	log.Println("\n--- Example 6: Get Available Trading Symbols ---")
	symbols, err := client.Rest.Market.GetSymbols()
	if err != nil {
		log.Printf("Error getting symbols: %v", err)
	} else {
		log.Printf("Code: %d, Retrieved %d symbols", symbols.Code, len(symbols.Data.Symbols))
		// Print first 5 symbols
		for i, s := range symbols.Data.Symbols {
			if i >= 5 {
				break
			}
			log.Printf("  Symbol: %s", s)
		}
	}

	// ===========================================
	// Account API Examples
	// ===========================================

	// Example 7: Get Account Balance
	log.Println("\n--- Example 7: Get Account Balance ---")
	balance, err := client.Rest.Account.GetBalance(accountreq.GetBalanceRequest{})
	if err != nil {
		log.Printf("Error getting balance: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", balance.Code, balance.Message)
		for i, bal := range balance.Data.Wallet {
			if i >= 5 {
				break
			}
			log.Printf("  %s: Available=%s, Frozen=%s, Total=%s",
				bal.Currency, bal.Available, bal.Frozen, bal.Total)
		}
	}

	// Example 8: Get Wallet Balance (by type)
	log.Println("\n--- Example 8: Get Spot Wallet Balance ---")
	walletBalance, err := client.Rest.Account.GetWalletBalance(accountreq.GetWalletBalanceRequest{
		WalletType: "spot",
	})
	if err != nil {
		log.Printf("Error getting wallet balance: %v", err)
	} else {
		log.Printf("Code: %d, Wallet Type: %s", walletBalance.Code, walletBalance.Data.WalletType)
		for i, bal := range walletBalance.Data.Balances {
			if i >= 3 {
				break
			}
			log.Printf("  %s: Available=%s", bal.Currency, bal.Available)
		}
	}

	// ===========================================
	// Trading API Examples
	// ===========================================

	// Example 9: Place Limit Order (Demo - will fail without valid credentials)
	log.Println("\n--- Example 9: Place Limit Order (Demo) ---")
	orderResp, err := client.Rest.Trade.PlaceOrder(tradereq.PlaceOrderRequest{
		Symbol: "BTC_USDT",
		Side:   "buy",
		Type:   "limit",
		Size:   "0.001",
		Price:  "30000",
	})
	if err != nil {
		log.Printf("Expected error (demo): %v", err)
	} else {
		log.Printf("Order placed! OrderID: %s", orderResp.Data.OrderID)
	}

	// Example 10: Get Order Details
	log.Println("\n--- Example 10: Get Order Details (Demo) ---")
	orderDetail, err := client.Rest.Trade.GetOrder(tradereq.GetOrderRequest{
		OrderID: "12345678", // Example order ID
	})
	if err != nil {
		log.Printf("Expected error (demo): %v", err)
	} else {
		log.Printf("Order: %s, Status: %s, Side: %s",
			orderDetail.Data.OrderID, orderDetail.Data.Status, orderDetail.Data.Side)
	}

	// Example 11: Get Orders List
	log.Println("\n--- Example 11: Get Orders List (Demo) ---")
	orders, err := client.Rest.Trade.GetOrders(tradereq.GetOrdersRequest{
		Symbol: "BTC_USDT",
		Status: "new",
		Limit:  10,
	})
	if err != nil {
		log.Printf("Expected error (demo): %v", err)
	} else {
		log.Printf("Retrieved %d orders", len(orders.Data))
	}

	// Example 12: Cancel Order
	log.Println("\n--- Example 12: Cancel Order (Demo) ---")
	cancelResp, err := client.Rest.Trade.CancelOrder(tradereq.CancelOrderRequest{
		Symbol:  "BTC_USDT",
		OrderID: "12345678",
	})
	if err != nil {
		log.Printf("Expected error (demo): %v", err)
	} else {
		log.Printf("Cancel result: %v", cancelResp.Data.Result)
	}

	// ===========================================
	// Funding API Examples
	// ===========================================

	// Example 13: Get Deposit Address
	log.Println("\n--- Example 13: Get Deposit Address (Demo) ---")
	depositAddr, err := client.Rest.Funding.GetDepositAddress(fundingreq.GetDepositAddressRequest{
		Currency: "USDT",
		Chain:    "TRC20",
	})
	if err != nil {
		log.Printf("Expected error (demo): %v", err)
	} else {
		log.Printf("Deposit Address: %s, Chain: %s",
			depositAddr.Data.Address, depositAddr.Data.Chain)
	}

	// Example 14: Get Deposit History
	log.Println("\n--- Example 14: Get Deposit History (Demo) ---")
	depositHistory, err := client.Rest.Funding.GetDepositHistory(fundingreq.GetDepositHistoryRequest{
		Currency: "USDT",
		Limit:    10,
	})
	if err != nil {
		log.Printf("Expected error (demo): %v", err)
	} else {
		log.Printf("Retrieved %d deposit records", len(depositHistory.Data.Records))
	}

	// Example 15: Get Withdrawal History
	log.Println("\n--- Example 15: Get Withdrawal History (Demo) ---")
	withdrawHistory, err := client.Rest.Funding.GetWithdrawHistory(fundingreq.GetWithdrawHistoryRequest{
		Currency: "USDT",
		Limit:    10,
	})
	if err != nil {
		log.Printf("Expected error (demo): %v", err)
	} else {
		log.Printf("Retrieved %d withdrawal records", len(withdrawHistory.Data.Records))
	}

	log.Println("\n✅ All examples completed!")
	log.Println("\nNote: Examples 9-15 require valid API credentials and will show errors in demo mode.")
	log.Println("Replace YOUR-API-KEY, YOUR-SECRET-KEY, and YOUR-MEMO with real credentials to test.")
}
