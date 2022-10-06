package order

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/pintu-crypto/b2b-order/client"
	"github.com/pintu-crypto/b2b-order/endpoint"
)

// Handler is the main order state machine.
type Handler struct {
	incoming client.IncomingChannel
	outgoing client.OutgoingChannel
	requests endpoint.RequestsChannel

	requestID        int64
	pendingResponses map[string]*endpoint.Request
	pendingRequests  map[int64]*endpoint.Request

	sessionID string

	closeC    chan interface{}
	closeWait sync.WaitGroup
}

// New initializes a order handler, services incoming client order requests,
// forwards those requests to the API, receives order and trade updates.
func New(incoming client.IncomingChannel,
	outgoing client.OutgoingChannel,
	requests endpoint.RequestsChannel) (res *Handler, err error) {
	res = &Handler{
		incoming:         incoming,
		outgoing:         outgoing,
		requests:         requests,
		pendingResponses: make(map[string]*endpoint.Request),
		pendingRequests:  make(map[int64]*endpoint.Request),
		closeC:           make(chan interface{}),
	}
	go res.runLoop()
	return
}

// Close stops the handler.
func (h *Handler) Close() {
	close(h.closeC)
	h.closeWait.Wait()
}

// runLoop is run forever to handle incoming events.
func (h *Handler) runLoop() {
	h.closeWait.Add(1)
	if err := h.handleInit(); err != nil {
		log.Printf("error during init " + err.Error())
	}
	if err := h.handleSubscribe(); err != nil {
		log.Printf("error during subscribe " + err.Error())
	}
	if err := h.handleRunning(); err != nil {
		log.Printf("error during run " + err.Error())
	}
	h.closeWait.Done()
}

// handleInit waits for the hello message.
func (h *Handler) handleInit() (err error) {
	msg := <-h.incoming
	hello := client.Hello{}
	if err = json.Unmarshal(msg, &hello); err != nil {
		err = errors.Wrap(err, "unable to decode hello message")
		return
	}
	// record the sessionID to use when placing orders, we want to
	// cancel all our orders if we get disconnected
	h.sessionID = hello.SessionID
	return
}

// handleSubscribe subscribes to execution reports and post trades.
func (h *Handler) handleSubscribe() (err error) {
	h.requestID++
	// subscribe to ExecutionReport. This will return any open orders and any future order updates.
	// subscribe to Trade, and recover any trades for the last 15 minutes.
	// Usually, this would checkpoint a persist store like a database.
	tradesStartDate := client.MicrosTimestamp(time.Now().Add(-15 * time.Minute))
	err = h.sendJSON(client.NewSubscribeRequest(time.Now(), h.requestID,
		client.StreamParameters{
			Name: "ExecutionReport",
		},
		client.StreamParameters{
			Name:      "Trade",
			StartDate: &tradesStartDate,
		}))
	if err != nil {
		err = errors.Wrap(err, "failed to send ExecutionReport subscribe")
		return
	}
	return
}

// handleRunning is the main handler that processes the next event,
// either a order request or a response from the websocket server.
func (h *Handler) handleRunning() (err error) {
	for {
		select {
		case data := <-h.incoming:
			log.Printf("received message %s\n", string(data))
			response := &client.Response{}
			err = json.Unmarshal(data, response)
			if err != nil {
				err = errors.Wrap(err, "unable to decode response")
				return
			}
			err = h.handleResponse(response)
			if err != nil {
				err = errors.Wrap(err, "error handling response")
				return
			}
		case request := <-h.requests:
			err = h.handleRequest(request)
			if err != nil {
				err = errors.Wrap(err, "error sending request")
				return
			}
		case <-h.closeC:
			return
		}
	}
}

// sendJSON sends the given request to the websocket server.
func (h *Handler) sendJSON(request interface{}) (err error) {
	data, err := json.Marshal(request)
	if err != nil {
		err = errors.Wrap(err, "unable to encode request")
		return
	}
	log.Printf("sending message %s\n", string(data))
	h.outgoing <- data
	return
}

// handleRequest processes a order request.
func (h *Handler) handleRequest(request *endpoint.Request) (err error) {
	h.requestID++

	newOrder := request.Message()
	// add the current sessionID to the request to ensure that it's cancelled if we're disconnected
	newOrder.CancelSessionID = h.sessionID
	h.pendingRequests[h.requestID] = request
	message := client.NewNewOrderSingleRequest(time.Now(), h.requestID, newOrder)
	err = h.sendJSON(message)
	if err != nil {
		return
	}
	h.pendingResponses[request.Message().ClOrdID] = request
	return
}

// handleResponse processes a response from the websocket server.
func (h *Handler) handleResponse(response *client.Response) (err error) {
	// check for any errors
	if response.Error != nil {
		return h.handleError(response.ReqID, *response.Error)
	}

	// then decode and process the response by type
	switch response.Type {
	case "ExecutionReport":
		for _, data := range response.Data {
			executionReport := &client.ExecutionReport{}
			err = json.Unmarshal(data, executionReport)
			if err != nil {
				err = errors.Wrap(err, "unable to decode execution report")
				return
			}
			log.Printf("received execution report %s\n", string(data))
			err = h.handleExecutionReport(response.ReqID, executionReport)
			if err != nil {
				err = errors.Wrap(err, "error handling execution report")
				return
			}
		}
	case "Trade":
		for _, data := range response.Data {
			trade := &client.Trade{}
			err = json.Unmarshal(data, trade)
			if err != nil {
				err = errors.Wrap(err, "unable to decode execution report")
				return
			}
			log.Printf("received trade %s\n", string(data))
			err = h.handleTrade(trade)
			if err != nil {
				err = errors.Wrap(err, "error handling execution report")
				return
			}
		}
	default:
		log.Printf("unhandled response %s\n", response.Data)
	}
	return
}

// handleError handles an error from the websocket.
func (h *Handler) handleError(requestID int64, e client.Error) (err error) {
	log.Printf("received error: " + e.Message)
	if request, ok := h.pendingRequests[requestID]; ok {
		request.Respond(fmt.Sprintf("rejected(%s)", e.Message))
		delete(h.pendingRequests, requestID)
		delete(h.pendingResponses, request.Message().ClOrdID)
	}
	return
}

// handleExecutionReport handles an execution report from the websocket server.
func (h *Handler) handleExecutionReport(requestID int64, report *client.ExecutionReport) (err error) {
	if request, ok := h.pendingResponses[report.ClOrdID]; ok {
		switch report.OrdStatus {
		case client.OrdStatus.DoneForDay, client.OrdStatus.Filled:
			request.Respond(fmt.Sprintf("filled(%s @ %s)", report.CumQty, report.AvgPx))
			delete(h.pendingRequests, requestID)
			delete(h.pendingResponses, report.ClOrdID)
		case client.OrdStatus.Rejected,
			client.OrdStatus.Canceled:
			if report.CumQty.IsZero() {
				request.Respond(fmt.Sprintf("rejected(%s)", report.Text))
				delete(h.pendingRequests, requestID)
				delete(h.pendingResponses, report.ClOrdID)
			}
			// otherwise, wait for the done for day
		}
	}
	return
}

// handleExecutionReport handles a post trade from the websocket server for reporting purposes.
func (h *Handler) handleTrade(trade *client.Trade) (err error) {
	// process the trade data (for example store it into DB)
	return
}
