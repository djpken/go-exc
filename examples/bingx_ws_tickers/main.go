package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	exc "github.com/djpken/go-exc"
)

const (
	API_KEY    = "" // public data does not require credentials
	SECRET_KEY = ""
)

var SYMBOLS = []string{"BTC-USDT", "ETH-USDT"}

func main() {
	ctx := context.Background()

	client, err := exc.NewExchange(ctx, exc.BingX, exc.Config{
		APIKey:    API_KEY,
		SecretKey: SECRET_KEY,
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	defer client.Close()

	tickerCh := make(chan *exc.TickerUpdate, 100)

	log.Printf("Subscribing to tickers: %v on %s\n", SYMBOLS, client.Name())
	if err := client.SubscribeTickers(tickerCh, SYMBOLS...); err != nil {
		log.Fatalf("SubscribeTickers: %v", err)
	}

	// Allow subscription to establish
	time.Sleep(2 * time.Second)
	log.Println("Listening for ticker updates (Ctrl+C to stop)...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case update := <-tickerCh:
			fmt.Printf("\n=== %s ===\n", update.Symbol)
			fmt.Printf("Last:     %s\n", update.LastPrice)
			fmt.Printf("Bid/Ask:  %s / %s\n", update.BidPrice, update.AskPrice)
			fmt.Printf("24h H/L:  %s / %s\n", update.High24h, update.Low24h)
			fmt.Printf("Volume:   %s\n", update.Volume24h)
			fmt.Printf("Time:     %s\n", time.Time(update.Timestamp).Format("15:04:05"))

		case <-sigCh:
			log.Println("Shutting down...")
			if err := client.UnsubscribeTickers(SYMBOLS...); err != nil {
				log.Printf("UnsubscribeTickers: %v", err)
			}
			return
		}
	}
}
