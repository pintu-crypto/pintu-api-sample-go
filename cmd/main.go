package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pintu-crypto/b2b-order/client"
	"github.com/pintu-crypto/b2b-order/endpoint"
	"github.com/pintu-crypto/b2b-order/liquidation"
)

var addr = flag.String("addr", "", "Pintu websocket address")
var apikey = flag.String("apikey", "", "Pintu api key")
var apisecret = flag.String("apisecret", "", "Pintu api secret")

var serveAddr = flag.String("serve-addr", ":8085", "Liquidate server address")

var interrupt = make(chan os.Signal, 1)

func init() {
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, syscall.SIGTERM)
}

func main() {
	flag.Parse()
	if *addr == "" {
		log.Fatalf("--addr required")
		return
	}
	if *apikey == "" {
		log.Fatalf("--apikey required")
		return
	}
	if *apisecret == "" {
		log.Fatalf("--apisecret required")
		return
	}

	requestsEndpoint, err := endpoint.Serve(*serveAddr)
	if err != nil {
		log.Fatalf("unable to create endpoint: %s", err)
		return
	}

	for attempt := 0; ; attempt++ {
		// check if the user requested shutdown
		select {
		case <-interrupt:
			return
		default:
		}

		// connect to the websocket and serve requests
		if runError := connectAndRun(*addr, *apikey, *apisecret, requestsEndpoint, attempt); runError != nil {
			log.Printf("received error: %s, re-connecting after 5 seconds", runError)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}

func connectAndRun(addr, apikey, apisecret string, requestsEndpoint *endpoint.Endpoint, attempt int) error {
	websocketClient, err := client.Connect(addr, apikey, apisecret)
	if err != nil {
		if attempt == 0 {
			// bail on connection if this is the first attempt
			log.Fatalf(err.Error())
		}
		return err
	}
	defer websocketClient.Close()

	handler, err := liquidation.New(websocketClient.IncomingChannel(),
		websocketClient.OutgoingChannel(),
		requestsEndpoint.RequestsChannel())
	if err != nil {
		log.Fatalf("unable to create liquidation handler %s", err)
	}
	defer handler.Close()

	for {
		// loop until the user requested shutdown or there was a connection error
		select {
		case <-interrupt:
			return nil
		case err = <-websocketClient.ErrorChannel():
			return err
		}
	}
}
