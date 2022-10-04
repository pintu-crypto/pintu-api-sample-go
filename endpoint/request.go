package endpoint

import "github.com/pintu-crypto/b2b-order/client"

// Request is an incoming client request to liquidate, for example. It has a message that represents the incoming
// request, and a channel to respond to the request.
type Request struct {
	message  *client.NewOrderSingle
	response chan string
}

// Message returns the client request data.
func (r *Request) Message() *client.NewOrderSingle {
	return r.message
}

// Respond should be called to send back a message with the outcome of this request. Must be called once the
// request has been resolved.
func (r *Request) Respond(message string) {
	r.response <- message
}
