package main

import (
	"context"
	"log"
	"time"

	"github.com/djpken/go-exc/exchanges/okex"
	"github.com/djpken/go-exc/exchanges/okex/events"
	"github.com/djpken/go-exc/exchanges/okex/events/private"
	"github.com/djpken/go-exc/exchanges/okex/events/public"
	ws_private_requests "github.com/djpken/go-exc/exchanges/okex/requests/ws/private"
	ws_public_requests "github.com/djpken/go-exc/exchanges/okex/requests/ws/public"
	"github.com/djpken/go-exc/exchanges/okex/constants"
	"github.com/djpken/go-exc/exchanges/okex/ws"
)

func main() {
	// Initialize API credentials
	apiKey := "cf5514a3-4913-4337-be39-d521795e3a13"
	secretKey := "6CE1209CDEEC014A6DE41D9BB583CDF3"
	passphrase := ".Vm3djpcl3gj94"

	ctx := context.Background()

	// Create OKEx client
	client, err := okex.NewClient(ctx, apiKey, secretKey, passphrase, okex.DemoServer)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("✅ OKEx WebSocket client created successfully")

	// Create channels for system messages and errors
	systemMsgChan := make(chan *ws.SystemMessage, 10)
	systemErrChan := make(chan *ws.SystemError, 10)
	client.Ws.SetSystemChannels(systemMsgChan, systemErrChan)

	// Monitor system messages and errors
	go func() {
		for {
			select {
			case msg := <-systemMsgChan:
				log.Printf("[SYSTEM] %s: %s (private=%v) at %s\n",
					msg.Type, msg.Message, msg.Private, msg.Timestamp.Format("15:04:05"))
			case err := <-systemErrChan:
				log.Printf("[ERROR] %s: %v (private=%v) at %s\n",
					err.Type, err.Error, err.Private, err.Timestamp.Format("15:04:05"))
			}
		}
	}()

	// Set up event channels (optional, for detailed control)
	errChan := make(chan *events.Error, 10)
	subChan := make(chan *events.Subscribe, 10)
	uSubChan := make(chan *events.Unsubscribe, 10)
	loginChan := make(chan *events.Login, 10)
	successChan := make(chan *events.Success, 10)
	client.Ws.SetChannels(errChan, subChan, uSubChan, loginChan, successChan)

	// Monitor events
	go func() {
		for {
			select {
			case login := <-loginChan:
				log.Printf("[LOGIN] Code: %s, Message: %s", login.Code, login.Msg)
			case sub := <-subChan:
				channel, _ := sub.Arg.Get("channel")
				log.Printf("[SUBSCRIBED] Channel: %s", channel)
			case unsub := <-uSubChan:
				channel, _ := unsub.Arg.Get("channel")
				log.Printf("[UNSUBSCRIBED] Channel: %s", channel)
			case wsErr := <-errChan:
				log.Printf("[WS ERROR] Code: %v, Message: %s, Op: %s",
					wsErr.Code, wsErr.Msg, wsErr.Op)
			}
		}
	}()

	// Example 1: Subscribe to Public Order Book
	log.Println("\n--- Example 1: Subscribe to Order Book ---")
	orderBookChan := make(chan *public.OrderBook, 10)
	err = client.Ws.Public.OrderBook([]ws_public_requests.OrderBook{{
		InstID:  "BTC-USDT",
		Channel: "books",
	}}, orderBookChan)
	if err != nil {
		log.Fatalf("Failed to subscribe to order book: %v", err)
	}

	// Listen to order book updates for a while
	go func() {
		count := 0
		for update := range orderBookChan {
			count++
			if count <= 5 { // Print first 5 updates
				channel, _ := update.Arg.Get("channel")
				instID, _ := update.Arg.Get("instId")
				log.Printf("[ORDER BOOK] Channel: %s, InstID: %s", channel, instID)
				for _, book := range update.Books {
					log.Printf("  Asks: %d, Bids: %d, Timestamp: %v",
						len(book.Asks), len(book.Bids), book.TS)
					if len(book.Asks) > 0 {
						log.Printf("    Best Ask: Price=%f, Qty=%f", book.Asks[0].DepthPrice, book.Asks[0].Size)
					}
					if len(book.Bids) > 0 {
						log.Printf("    Best Bid: Price=%f, Qty=%f", book.Bids[0].DepthPrice, book.Bids[0].Size)
					}
				}
			}
		}
	}()

	// Example 2: Subscribe to Ticker
	log.Println("\n--- Example 2: Subscribe to Ticker ---")
	tickerChan := make(chan *public.Tickers, 10)
	err = client.Ws.Public.Tickers(ws_public_requests.Tickers{
		InstID: "BTC-USDT",
	}, tickerChan)
	if err != nil {
		log.Printf("Failed to subscribe to ticker: %v", err)
	}

	// Listen to ticker updates
	go func() {
		count := 0
		for update := range tickerChan {
			count++
			if count <= 5 { // Print first 5 updates
				for _, ticker := range update.Tickers {
					log.Printf("[TICKER] %s: Last=%v, High24h=%v, Low24h=%v, Vol24h=%v, Time=%v",
						ticker.InstID, ticker.Last, ticker.High24h, ticker.Low24h,
						ticker.Vol24h, ticker.TS)
				}
			}
		}
	}()

	// Example 3: Subscribe to Private Account Channel (requires authentication)
	log.Println("\n--- Example 3: Subscribe to Account Updates ---")
	accountChan := make(chan *private.Account, 10)
	err = client.Ws.Private.Account(ws_private_requests.Account{
		Ccy: "",
	}, accountChan)
	if err != nil {
		log.Printf("Failed to subscribe to account: %v", err)
	}

	// Listen to account updates
	go func() {
		for update := range accountChan {
			log.Printf("[ACCOUNT] Update received, EventType: %s", update.EventType)
			for _, balance := range update.Balances {
				log.Printf("  Total Equity: %v, Update Time: %v", balance.TotalEq, balance.UTime)
				for _, detail := range balance.Details {
					log.Printf("    Currency: %s, Available: %v, Frozen: %v, Equity: %v",
						detail.Ccy, detail.AvailBal, detail.FrozenBal, detail.Eq)
				}
			}
		}
	}()

	// Example 4: Subscribe to Position Updates
	log.Println("\n--- Example 4: Subscribe to Position Updates ---")
	positionChan := make(chan *private.Position, 10)
	err = client.Ws.Private.Position(ws_private_requests.Position{
		InstType: constants.SwapInstrument,
	}, positionChan)
	if err != nil {
		log.Printf("Failed to subscribe to positions: %v", err)
	}

	// Listen to position updates
	go func() {
		for update := range positionChan {
			log.Printf("[POSITION] Update received, EventType: %s", update.EventType)
			for _, pos := range update.Positions {
				log.Printf("  Instrument: %s, Side: %s, Qty: %v, AvgPx: %v, Upl: %v",
					pos.InstID, pos.PosSide, pos.Pos, pos.AvgPx, pos.Upl)
			}
		}
	}()

	// Example 5: Subscribe to Balance and Position Updates (combined)
	log.Println("\n--- Example 5: Subscribe to Balance and Position ---")
	balPosChan := make(chan *private.BalanceAndPosition, 10)
	err = client.Ws.Private.BalanceAndPosition(balPosChan)
	if err != nil {
		log.Printf("Failed to subscribe to balance and position: %v", err)
	}

	// Listen to combined updates
	go func() {
		for update := range balPosChan {
			log.Printf("[BALANCE & POSITION] Update received")
			for _, balPos := range update.BalanceAndPositions {
				log.Printf("  EventType: %s, Time: %v", balPos.EventType, balPos.PTime)

				// Balance data
				for _, bal := range balPos.BalData {
					log.Printf("    Balance - Currency: %s, CashBal: %v", bal.Ccy, bal.CashBal)
				}

				// Position data
				for _, pos := range balPos.PosData {
					log.Printf("    Position - Instrument: %s, Qty: %v, AvgPx: %v",
						pos.InstId, pos.Pos, pos.AvgPx)
				}
			}
		}
	}()

	// Keep the program running to receive updates
	log.Println("\n📡 WebSocket connections established. Listening for updates...")
	log.Println("Press Ctrl+C to exit")

	// Run for 30 seconds to see updates
	time.Sleep(30 * time.Second)

	log.Println("\n✅ WebSocket example completed!")
}
