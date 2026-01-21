package main

import (
	"context"
	"log"

	"github.com/djpken/go-exc/exchanges/okex"
	accountreq "github.com/djpken/go-exc/exchanges/okex/requests/rest/account"
	tradereq "github.com/djpken/go-exc/exchanges/okex/requests/rest/trade"
)

func main() {
	// Initialize API credentials
	apiKey := "cf5514a3-4913-4337-be39-d521795e3a13"
	secretKey := "6CE1209CDEEC014A6DE41D9BB583CDF3"
	passphrase := ".Vm3djpcl3gj94"

	ctx := context.Background()

	// Create OKEx client
	// Use okex.DemoServer for testing, okex.NormalServer for production
	client, err := okex.NewClient(ctx, apiKey, secretKey, passphrase, okex.DemoServer)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("✅ OKEx client created successfully")

	// Example 1: Get Account Balance
	log.Println("\n--- Example 1: Get Account Balance ---")
	balanceResp, err := client.Rest.Account.GetBalance(accountreq.GetBalance{
		Ccy: []string{"BTC", "USDT"}, // Optional: specify currencies
	})
	if err != nil {
		log.Printf("Error getting balance: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", balanceResp.Code, balanceResp.Msg)
		for _, bal := range balanceResp.Balances {
			log.Printf("Total Equity: %v", bal.TotalEq)
			for _, detail := range bal.Details {
				log.Printf("  %s: Available=%v, Frozen=%v, Total=%v",
					detail.Ccy, detail.AvailBal, detail.FrozenBal, detail.Eq)
			}
		}
	}

	// Example 2: Get Positions
	log.Println("\n--- Example 2: Get Positions ---")
	positionsResp, err := client.Rest.Account.GetPositions(accountreq.GetPositions{
		InstID: []string{"BTC-USDT-SWAP"}, // Optional: specify instruments
	})
	if err != nil {
		log.Printf("Error getting positions: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", positionsResp.Code, positionsResp.Msg)
		for _, pos := range positionsResp.Positions {
			log.Printf("Position: %s, Side=%s, Quantity=%v, AvgPx=%v, Upl=%v",
				pos.InstID, pos.PosSide, pos.Pos, pos.AvgPx, pos.Upl)
		}
	}

	// Example 3: Place a Limit Order
	log.Println("\n--- Example 3: Place Limit Order ---")
	orderResp, err := client.Rest.Trade.PlaceOrder([]tradereq.PlaceOrder{
		{
			InstID:  "BTC-USDT",
			TdMode:  "cash",    // Trading mode: cash, cross, isolated
			Side:    "buy",     // buy or sell
			OrdType: "limit",   // limit, market, post_only, fok, ioc
			Sz:      0.001,     // Order size
			Px:      30000.0,   // Price
			Tag:     "example", // Optional tag
		},
	})
	if err != nil {
		log.Printf("Error placing order: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", orderResp.Code, orderResp.Msg)
		for _, order := range orderResp.PlaceOrders {
			log.Printf("Order placed: ID=%v, ClientOrderID=%s, Code=%v, Message=%s",
				order.OrdID, order.ClOrdID, order.SCode, order.SMsg)
		}
	}

	// Example 4: Get Order Details
	log.Println("\n--- Example 4: Get Order Details ---")
	orderDetailResp, err := client.Rest.Trade.GetOrderDetail(tradereq.OrderDetails{
		InstID: "BTC-USDT",
		OrdID:  "123456789", // Replace with actual order ID
	})
	if err != nil {
		log.Printf("Error getting order details: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", orderDetailResp.Code, orderDetailResp.Msg)
		for _, order := range orderDetailResp.Orders {
			log.Printf("Order: ID=%s, State=%s, Px=%v, Sz=%v, FilledSz=%v",
				order.OrdID, order.State, order.Px, order.Sz, order.AccFillSz)
		}
	}

	// Example 5: Cancel Order
	log.Println("\n--- Example 5: Cancel Order ---")
	cancelResp, err := client.Rest.Trade.CandleOrder([]tradereq.CancelOrder{
		{
			InstID: "BTC-USDT",
			OrdID:  "123456789", // Replace with actual order ID
		},
	})
	if err != nil {
		log.Printf("Error canceling order: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", cancelResp.Code, cancelResp.Msg)
	}

	// Example 6: Get Order History
	log.Println("\n--- Example 6: Get Order History ---")
	historyResp, err := client.Rest.Trade.GetOrderHistory(tradereq.OrderList{
		InstType: "SPOT",
		InstID:   "BTC-USDT",
		Limit:    10, // Get last 10 orders
	}, false) // false for recent history, true for archived
	if err != nil {
		log.Printf("Error getting order history: %v", err)
	} else {
		log.Printf("Code: %d, Message: %s", historyResp.Code, historyResp.Msg)
		log.Printf("Found %d orders in history", len(historyResp.Orders))
		for _, order := range historyResp.Orders {
			log.Printf("  Order: ID=%s, State=%s, Side=%s, Px=%v, Sz=%v",
				order.OrdID, order.State, order.Side, order.Px, order.Sz)
		}
	}

	log.Println("\n✅ All examples completed!")
}
