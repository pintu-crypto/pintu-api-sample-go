package endpoint

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/pintu-crypto/b2b-order/client"
)

// RequestsChannel is a channel of http requests.
type RequestsChannel <-chan *Request

// Endpoint is the REST endpoint for the liquidate API.
type Endpoint struct {
	requests chan *Request
	addr     string
}

// Serve returns an http endpoint, which provides the client facing liquidation REST API.
func Serve(addr string) (result *Endpoint, err error) {
	result = &Endpoint{
		addr:     addr,
		requests: make(chan *Request),
	}

	go result.runServe()
	return
}

// RequestsChannel returns a channel that incoming http requests are dispatched to.
func (e *Endpoint) RequestsChannel() RequestsChannel {
	return e.requests
}

func (e *Endpoint) runServe() {
	http.HandleFunc("/liquidate", e.handleClientRequest)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	go func() {
		if err := http.ListenAndServe(e.addr, nil); err != nil {
			log.Fatalf("error listening: %s", err)
		}
	}()

}

func getQueryKeyValue(r *http.Request, key string, required bool) (value string, err error) {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		if required {
			err = fmt.Errorf("missing required parameter '%s'", key)
			return
		}
		return
	}
	value = keys[0]
	return
}

func (e *Endpoint) handleClientRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("received client request %s", r.URL)
	symbol, err := getQueryKeyValue(r, "symbol", true)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currency, err := getQueryKeyValue(r, "currency", false)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sideString, err := getQueryKeyValue(r, "side", true)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	side, err := client.ParseSide(sideString)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	quantityString, err := getQueryKeyValue(r, "quantity", true)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	quantity, err := decimal.NewFromString(quantityString)
	if err != nil {
		err = errors.Wrapf(err, "invalid quantity %s", quantityString)
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// generate a NewOrderSingle structure that will be used to submit a market order
	clOrdID := uuid.New().String()
	newOrderRequest := &Request{
		message: &client.NewOrderSingle{
			Symbol:       symbol,
			Currency:     currency,
			ClOrdID:      clOrdID,
			Side:         side,
			OrderQty:     quantity,
			OrdType:      client.OrdType.Market,
			TimeInForce:  client.TimeInForce.FillOrKill,
			TransactTime: client.MicrosTimestamp(time.Now()),
		},
		response: make(chan string),
	}
	e.requests <- newOrderRequest

	// block on the response channel until we get a response
	response := <-newOrderRequest.response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	log.Printf("sending client response %s", response)
	_, _ = fmt.Fprintf(w, response)
}
