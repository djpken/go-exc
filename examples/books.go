package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/djpken/go-exc/exchanges/okex"
	"github.com/djpken/go-exc/exchanges/okex/events/private"
	"github.com/djpken/go-exc/exchanges/okex/ws"
)

func main() {

	// Start the pprof server
	go func() {
		log.Println("Starting pprof server on localhost:6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("could not start pprof server: %v", err)
		}
	}()

	apiKey := "cf5514a3-4913-4337-be39-d521795e3a13"
	secretKey := "6CE1209CDEEC014A6DE41D9BB583CDF3"
	passphrase := ".Vm3djpcl3gj94"
	ctx := context.Background()
	client, err := okex.NewClient(ctx, apiKey, secretKey, passphrase, okex.DemoServer)
	if err != nil {
		log.Fatalln(err)
	}

	// Create channels for system messages and errors
	systemMsgChan := make(chan *ws.SystemMessage, 10)
	systemErrChan := make(chan *ws.SystemError, 10)
	client.Ws.SetSystemChannels(systemMsgChan, systemErrChan)

	// Monitor system messages and errors in a separate goroutine
	go func() {
		for {
			select {
			case msg := <-systemMsgChan:
				log.Printf("[%s] %s (private=%v) at %s\n",
					msg.Type, msg.Message, msg.Private, msg.Timestamp.Format("15:04:05"))
			case err := <-systemErrChan:
				log.Printf("[%s ERROR] %v (private=%v) at %s\n",
					err.Type, err.Error, err.Private, err.Timestamp.Format("15:04:05"))
			}
		}
	}()

	positions := make(chan *private.BalanceAndPosition)
	err = client.Ws.Private.BalanceAndPosition(positions)
	if err != nil {
		panic(err)
	}

	// Listen for updates
	for update := range positions {
		log.Printf("Received order book update: %+v\n", update)
		insId, _ := update.Arg.Get("instId")
		log.Printf("Instrument ID: %s\n", insId)
	}
}
