package client

import (
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 5 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 1 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 32768
)

// IncomingChannel is used to receive a message from the websocket.
type IncomingChannel <-chan []byte

// OutgoingChannel is used to send a message to the websocket.
type OutgoingChannel chan<- []byte

// ErrorChannel is used to receive connection errors.
type ErrorChannel <-chan error

type client struct {
	conn               *websocket.Conn
	incoming, outgoing chan []byte
	errorC             chan error

	closeC         chan interface{}
	closeRequested int32
}

// Connect connects to the Pintu websocket API on the given address, which should be a full
// websocket address such as wss://partner.pintu.co.id/ws/v1. It dispatches incoming messages
// to the incoming channel. To send a message, use the outgoing channel.
func Connect(addr string, apikey string, apisecret string) (result *client, err error) {
	var conn *websocket.Conn

	uri, err := url.Parse(addr)
	if err != nil {
		err = errors.Wrapf(err, "invalid url %s", addr)
		return
	}

	// connect
	ts := time.Now()
	hostAndPort := uri.Hostname()
	if uri.Port() != "" {
		hostAndPort += ":" + uri.Port()
	}
	signature := sign(apisecret, "GET", ts, hostAndPort, uri.Path)

	dialer := &websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	header := http.Header{
		"ApiKey":       []string{apikey},
		"ApiSign":      []string{signature},
		"ApiTimestamp": []string{MicrosTimestamp(ts).String()},
	}

	log.Printf("connecting to %s", addr)
	if conn, _, err = dialer.Dial(addr, header); err != nil {
		err = errors.Wrapf(err, "unable to connect to %s", addr)
		return
	}

	result = &client{
		incoming: make(chan []byte, 1000),
		outgoing: make(chan []byte, 1000),
		errorC:   make(chan error),
		closeC:   make(chan interface{}),
		conn:     conn,
	}
	go result.writePump()
	go result.readPump()
	log.Printf("successfully connected to %s", addr)
	return
}

// IncomingChannel returns the channel to receive websocket messages from the server.
func (client *client) IncomingChannel() IncomingChannel {
	return client.incoming
}

// OutgoingChannel returns the channel to send websocket messages to the server.
func (client *client) OutgoingChannel() OutgoingChannel {
	return client.outgoing
}

// ErrorChannel returns a channel for any connection errors. Errors must be processed by the user.
func (client *client) ErrorChannel() ErrorChannel {
	return client.errorC
}

// Close closes the websocket connection.
func (client *client) Close() {
	atomic.StoreInt32(&client.closeRequested, 1)
	_ = client.conn.Close()
	<-client.closeC
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (client *client) readPump() {
	defer func() {
		_ = client.conn.Close()
	}()
	client.conn.SetReadLimit(maxMessageSize)
	if err := client.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		client.onError(err)
		return
	}
	client.conn.SetPongHandler(func(string) error {
		return client.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			client.onError(err)
			return
		}
		client.incoming <- message
	}
}

func (client *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.outgoing:
			if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				client.onError(err)
				return
			}
			if !ok {
				// The hub closed the channel.
				if err := client.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					client.onError(err)
					return
				}
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err = w.Write(message); err != nil {
				client.onError(err)
				return
			}
			if err = w.Close(); err != nil {
				client.onError(err)
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *client) onError(err error) {
	// if there's an error, the connection is closed and can't be re-used
	close(client.closeC)
	if atomic.LoadInt32(&client.closeRequested) != 0 {
		// close was requested, so unblock the caller and don't forward an error
		return
	}
	log.Printf("error: %v", err)
	client.errorC <- err
}
