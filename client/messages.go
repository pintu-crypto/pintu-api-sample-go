package client

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// Request contains the base set of fields for a websocket request.
type request struct {
	Id        int64           `json:"reqid"`
	Type      string          `json:"type"`
	Timestamp MicrosTimestamp `json:"ts"`
}

// StreamParameters contains the data for a subscription.
type StreamParameters struct {
	Name      string           `json:"name"`
	StartDate *MicrosTimestamp `json:"StartDate,omitempty"`
	EndDate   *MicrosTimestamp `json:"EndDate,omitempty"`
}

// SubscribeRequest contains the data for a subscription.
type subscribeRequest struct {
	request
	Streams []StreamParameters `json:"streams,omitempty"`
}

// NewSubscribeRequest returns a subscribe request for the given stream names.
func NewSubscribeRequest(now time.Time, requestID int64,
	streams ...StreamParameters) (result *subscribeRequest) {
	result = &subscribeRequest{
		request: request{
			Id:        requestID,
			Type:      "subscribe",
			Timestamp: MicrosTimestamp(now),
		},
		Streams: streams,
	}
	return
}

// Hello is the message that's sent by the server on connection.
type Hello struct {
	Type      string          `json:"type"`
	Timestamp MicrosTimestamp `json:"ts"`
	SessionID string          `json:"session_id"`
}

// Error is sent as part of a response on an error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// Response is a generic response containing data from the server. The Data field should be decoded
// according to the value of Type.
type Response struct {
	ReqID     int64             `json:"reqid"`
	Type      string            `json:"type"`
	Seq       int64             `json:"seq"`
	Timestamp MicrosTimestamp   `json:"ts"`
	Error     *Error            `json:"error"`
	Data      []json.RawMessage `json:"data"`
}

// ExecutionReport is the result of a subscription to execution reports and is an update on a parent
// order. It is returned in the Data field on a response.
type ExecutionReport struct {
	Timestamp       MicrosTimestamp
	User            string
	Symbol          string
	OrderID         string
	ClOrdID         string
	OrigClOrdID     string
	SubmitTime      MicrosTimestamp
	ExecID          string
	Side            SideEnum
	TransactTime    MicrosTimestamp
	ExecType        ExecTypeEnum
	OrdStatus       OrdStatusEnum
	OrderQty        decimal.Decimal
	OrdType         OrdTypeEnum
	Price           decimal.Decimal
	Currency        string
	LeavesQty       decimal.Decimal
	CumQty          decimal.Decimal
	AvgPx           decimal.Decimal
	TimeInForce     TimeInForceEnum
	LastMarket      string
	LastPx          decimal.Decimal
	LastQty         decimal.Decimal
	LastAmt         decimal.Decimal
	LastFee         decimal.Decimal
	LastFeeCurrency string
	CumAmt          decimal.Decimal
	CumFee          decimal.Decimal
	FeeCurrency     string
	OrdRejReason    OrdRejReasonEnum
	LastExecID      string
	CxlRejReason    CxlRejReasonEnum
	AmountCurrency  string
	SessionID       string
	CancelSessionID string
	SubAccount      string
	Group           string
	Text            string
}

// Trade is the result of a subscription to trades and is a trade on a market.
// It is returned in the Data field on a response.
type Trade struct {
	Timestamp      MicrosTimestamp
	User           string
	Symbol         string
	OrderID        string
	TradeID        string
	Side           SideEnum
	TransactTime   MicrosTimestamp
	Price          decimal.Decimal
	Quantity       decimal.Decimal
	Currency       string
	Market         string
	Amount         decimal.Decimal
	Fee            decimal.Decimal
	FeeCurrency    string
	MarketTradeID  string
	TradeStatus    string
	AggressorSide  SideEnum
	AmountCurrency string
	DealtCurrency  string
	SubAccount     string
	Group          string
}

// NewOrderSingle is a request to submit an order. It should be sent as the Data
// field on a request.
type NewOrderSingle struct {
	Symbol          string
	Currency        string `json:",omitempty"`
	ClOrdID         string
	Side            SideEnum
	OrderQty        decimal.Decimal
	OrdType         OrdTypeEnum
	Price           *decimal.Decimal `json:",omitempty"`
	TimeInForce     TimeInForceEnum
	TransactTime    MicrosTimestamp
	CancelSessionID string `json:",omitempty"`
}

// NewOrderSingleRequest is a request message for a new order.
type newOrderSingleRequest struct {
	request
	Data []NewOrderSingle `json:"data"`
}

// NewNewOrderSingleRequest returns a new order single request with the given params.
func NewNewOrderSingleRequest(now time.Time, requestID int64, message *NewOrderSingle) (result *newOrderSingleRequest) {
	result = &newOrderSingleRequest{
		request: request{
			Id:        requestID,
			Type:      "NewOrderSingle",
			Timestamp: MicrosTimestamp(now),
		},
		Data: []NewOrderSingle{
			*message,
		},
	}
	return
}
