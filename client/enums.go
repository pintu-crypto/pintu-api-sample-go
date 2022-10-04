package client

import (
	"fmt"
	"strings"
)

type SideEnum uint8
type SideValues struct {
	Buy  SideEnum
	Sell SideEnum
}

var Side = SideValues{1, 2}

func SideString(s SideEnum) string {
	switch s {
	case Side.Buy:
		return "Buy"
	case Side.Sell:
		return "Sell"
	default:
		return "Unknown"
	}
}

func ParseSide(str string) (s SideEnum, err error) {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	upper := strings.ToUpper(str)
	switch upper {
	case "BUY":
		s = Side.Buy
	case "SELL":
		s = Side.Sell
	default:
		err = fmt.Errorf("invalid Side %s", str)
	}
	return
}

func (e SideEnum) MarshalJSON() ([]byte, error) {
	return []byte("\"" + SideString(e) + "\""), nil
}

func (e *SideEnum) UnmarshalJSON(b []byte) (err error) {
	*e, err = ParseSide(string(b))
	return
}

type ExecTypeEnum uint8
type ExecTypeValues struct {
	New             ExecTypeEnum
	Trade           ExecTypeEnum
	Canceled        ExecTypeEnum
	Replaced        ExecTypeEnum
	PendingCancel   ExecTypeEnum
	Rejected        ExecTypeEnum
	PendingNew      ExecTypeEnum
	Restated        ExecTypeEnum
	PendingReplace  ExecTypeEnum
	DoneForDay      ExecTypeEnum
	CancelRejected  ExecTypeEnum
	ReplaceRejected ExecTypeEnum
	Expired         ExecTypeEnum
	Stale           ExecTypeEnum
}

var ExecType = ExecTypeValues{0, 2, 4, 5, 6, 8, 10, 13, 14, 15, 37, 38, 12, 16}

func ExecTypeString(s ExecTypeEnum) string {
	switch s {
	case ExecType.New:
		return "New"
	case ExecType.Trade:
		return "Trade"
	case ExecType.Canceled:
		return "Canceled"
	case ExecType.Replaced:
		return "Replaced"
	case ExecType.PendingCancel:
		return "PendingCancel"
	case ExecType.Rejected:
		return "Rejected"
	case ExecType.PendingNew:
		return "PendingNew"
	case ExecType.Restated:
		return "Restated"
	case ExecType.PendingReplace:
		return "PendingReplace"
	case ExecType.DoneForDay:
		return "DoneForDay"
	case ExecType.CancelRejected:
		return "CancelRejected"
	case ExecType.ReplaceRejected:
		return "ReplaceRejected"
	case ExecType.Expired:
		return "Expired"
	case ExecType.Stale:
		return "Stale"
	default:
		return "Unknown"
	}
}

func (e ExecTypeEnum) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ExecTypeString(e) + "\""), nil
}

func (e *ExecTypeEnum) UnmarshalJSON(b []byte) (err error) {
	*e, err = ParseExecType(string(b))
	return
}

func ParseExecType(str string) (s ExecTypeEnum, err error) {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	upper := strings.ToUpper(str)
	switch upper {
	case "NEW":
		s = ExecType.New
	case "TRADE":
		s = ExecType.Trade
	case "CANCELED":
		s = ExecType.Canceled
	case "REPLACED":
		s = ExecType.Replaced
	case "PENDINGCANCEL":
		s = ExecType.PendingCancel
	case "REJECTED":
		s = ExecType.Rejected
	case "PENDINGNEW":
		s = ExecType.PendingNew
	case "RESTATED":
		s = ExecType.Restated
	case "PENDINGREPLACE":
		s = ExecType.PendingReplace
	case "DONEFORDAY":
		s = ExecType.DoneForDay
	case "CANCELREJECTED":
		s = ExecType.CancelRejected
	case "REPLACEREJECTED":
		s = ExecType.ReplaceRejected
	case "EXPIRED":
		s = ExecType.Expired
	case "STALE":
		s = ExecType.Stale
	default:
		err = fmt.Errorf("invalid ExecType %s", str)
	}
	return
}

type OrdStatusEnum uint8
type OrdStatusValues struct {
	New             OrdStatusEnum
	PartiallyFilled OrdStatusEnum
	Filled          OrdStatusEnum
	Canceled        OrdStatusEnum
	Replaced        OrdStatusEnum
	PendingCancel   OrdStatusEnum
	Rejected        OrdStatusEnum
	PendingNew      OrdStatusEnum
	PendingReplace  OrdStatusEnum
	DoneForDay      OrdStatusEnum
}

var OrdStatus = OrdStatusValues{0, 1, 2, 4, 5, 6, 8, 10, 14, 15}

func OrdStatusString(s OrdStatusEnum) string {
	switch s {
	case OrdStatus.New:
		return "New"
	case OrdStatus.PartiallyFilled:
		return "PartiallyFilled"
	case OrdStatus.Filled:
		return "Filled"
	case OrdStatus.Canceled:
		return "Canceled"
	case OrdStatus.Replaced:
		return "Replaced"
	case OrdStatus.PendingCancel:
		return "PendingCancel"
	case OrdStatus.Rejected:
		return "Rejected"
	case OrdStatus.PendingNew:
		return "PendingNew"
	case OrdStatus.PendingReplace:
		return "PendingReplace"
	case OrdStatus.DoneForDay:
		return "DoneForDay"
	default:
		return "Unknown"
	}
}

func ParseOrdStatus(str string) (s OrdStatusEnum, err error) {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	upper := strings.ToUpper(str)
	switch upper {
	case "NEW":
		s = OrdStatus.New
	case "PARTIALLYFILLED":
		s = OrdStatus.PartiallyFilled
	case "FILLED":
		s = OrdStatus.Filled
	case "CANCELED":
		s = OrdStatus.Canceled
	case "REPLACED":
		s = OrdStatus.Replaced
	case "PENDINGCANCEL":
		s = OrdStatus.PendingCancel
	case "REJECTED":
		s = OrdStatus.Rejected
	case "PENDINGNEW":
		s = OrdStatus.PendingNew
	case "PENDINGREPLACE":
		s = OrdStatus.PendingReplace
	case "DONEFORDAY":
		s = OrdStatus.DoneForDay
	default:
		err = fmt.Errorf("invalid OrdStatus %s", str)
	}
	return
}

func (e OrdStatusEnum) MarshalJSON() ([]byte, error) {
	return []byte("\"" + OrdStatusString(e) + "\""), nil
}

func (e *OrdStatusEnum) UnmarshalJSON(b []byte) (err error) {
	*e, err = ParseOrdStatus(string(b))
	return
}

type OrdTypeEnum uint8
type OrdTypeValues struct {
	Market OrdTypeEnum
	Limit  OrdTypeEnum
	RFQ    OrdTypeEnum
}

var OrdType = OrdTypeValues{1, 2, 3}

func OrdTypeString(s OrdTypeEnum) string {
	switch s {
	case OrdType.Market:
		return "Market"
	case OrdType.Limit:
		return "Limit"
	case OrdType.RFQ:
		return "RFQ"
	default:
		return "Unknown"
	}
}

func ParseOrdType(str string) (s OrdTypeEnum, err error) {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	upper := strings.ToUpper(str)
	switch upper {
	case "MARKET":
		s = OrdType.Market
	case "LIMIT":
		s = OrdType.Limit
	case "RFQ":
		s = OrdType.RFQ
	default:
		err = fmt.Errorf("invalid OrdType %s", str)
	}
	return
}

func (e OrdTypeEnum) MarshalJSON() ([]byte, error) {
	return []byte("\"" + OrdTypeString(e) + "\""), nil
}

func (e *OrdTypeEnum) UnmarshalJSON(b []byte) (err error) {
	*e, err = ParseOrdType(string(b))
	return
}

type OrdRejReasonEnum uint8
type OrdRejReasonValues struct {
	UnknownSymbol                         OrdRejReasonEnum
	ExchangeClosed                        OrdRejReasonEnum
	OrderExceedsLimit                     OrdRejReasonEnum
	TooLateToEnter                        OrdRejReasonEnum
	UnknownOrder                          OrdRejReasonEnum
	DuplicateOrder                        OrdRejReasonEnum
	DuplicateOfAVerballyCommunicatedOrder OrdRejReasonEnum
	StaleOrder                            OrdRejReasonEnum
	UnknownMarket                         OrdRejReasonEnum
	InternalError                         OrdRejReasonEnum
	BrokerOption                          OrdRejReasonEnum
	RateLimit                             OrdRejReasonEnum
	ForceCancel                           OrdRejReasonEnum
}

var OrdRejReason = OrdRejReasonValues{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

func OrdRejReasonString(s OrdRejReasonEnum) string {
	switch s {
	case OrdRejReason.UnknownSymbol:
		return "UnknownSymbol"
	case OrdRejReason.ExchangeClosed:
		return "ExchangeClosed"
	case OrdRejReason.OrderExceedsLimit:
		return "OrderExceedsLimit"
	case OrdRejReason.TooLateToEnter:
		return "TooLateToEnter"
	case OrdRejReason.UnknownOrder:
		return "UnknownOrder"
	case OrdRejReason.DuplicateOrder:
		return "DuplicateOrder"
	case OrdRejReason.DuplicateOfAVerballyCommunicatedOrder:
		return "DuplicateOfAVerballyCommunicatedOrder"
	case OrdRejReason.StaleOrder:
		return "StaleOrder"
	case OrdRejReason.UnknownMarket:
		return "UnknownMarket"
	case OrdRejReason.InternalError:
		return "InternalError"
	case OrdRejReason.BrokerOption:
		return "BrokerOption"
	case OrdRejReason.RateLimit:
		return "RateLimit"
	case OrdRejReason.ForceCancel:
		return "ForceCancel"
	default:
		return "Unknown"
	}
}

func ParseOrdRejReason(str string) (s OrdRejReasonEnum, err error) {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	upper := strings.ToUpper(str)
	switch upper {
	case "UNKNOWNSYMBOL":
		s = OrdRejReason.UnknownSymbol
	case "EXCHANGECLOSED":
		s = OrdRejReason.ExchangeClosed
	case "ORDEREXCEEDSLIMIT":
		s = OrdRejReason.OrderExceedsLimit
	case "TOOLATETOENTER":
		s = OrdRejReason.TooLateToEnter
	case "UNKNOWNORDER":
		s = OrdRejReason.UnknownOrder
	case "DUPLICATEORDER":
		s = OrdRejReason.DuplicateOrder
	case "DUPLICATEOFAVERBALLYCOMMUNICATEDORDER":
		s = OrdRejReason.DuplicateOfAVerballyCommunicatedOrder
	case "STALEORDER":
		s = OrdRejReason.StaleOrder
	case "UNKNOWNMARKET":
		s = OrdRejReason.UnknownMarket
	case "INTERNALERROR":
		s = OrdRejReason.InternalError
	case "BROKEROPTION":
		s = OrdRejReason.BrokerOption
	case "RATELIMIT":
		s = OrdRejReason.RateLimit
	case "FORCECANCEL":
		s = OrdRejReason.ForceCancel
	default:
		err = fmt.Errorf("invalid OrdRejReason %s", str)
	}
	return
}

func (e OrdRejReasonEnum) MarshalJSON() ([]byte, error) {
	return []byte("\"" + OrdRejReasonString(e) + "\""), nil
}

func (e *OrdRejReasonEnum) UnmarshalJSON(b []byte) (err error) {
	*e, err = ParseOrdRejReason(string(b))
	return
}

type CxlRejReasonEnum uint8
type CxlRejReasonValues struct {
	UnknownOrder                                      CxlRejReasonEnum
	Broker                                            CxlRejReasonEnum
	OrderAlreadyInPendingCancelOrPendingReplaceStatus CxlRejReasonEnum
	UnableToProcessOrderMassCancelRequest             CxlRejReasonEnum
	OrigOrdModTime                                    CxlRejReasonEnum
	DuplicateClOrdID                                  CxlRejReasonEnum
	TooLateToCancel                                   CxlRejReasonEnum
	StaleRequest                                      CxlRejReasonEnum
	RateLimit                                         CxlRejReasonEnum
	Other                                             CxlRejReasonEnum
}

var CxlRejReason = CxlRejReasonValues{1, 2, 3, 4, 5, 6, 7, 8, 9, 99}

func CxlRejReasonString(s CxlRejReasonEnum) string {
	switch s {
	case CxlRejReason.UnknownOrder:
		return "UnknownOrder"
	case CxlRejReason.Broker:
		return "Broker"
	case CxlRejReason.OrderAlreadyInPendingCancelOrPendingReplaceStatus:
		return "OrderAlreadyInPendingCancelOrPendingReplaceStatus"
	case CxlRejReason.UnableToProcessOrderMassCancelRequest:
		return "UnableToProcessOrderMassCancelRequest"
	case CxlRejReason.OrigOrdModTime:
		return "OrigOrdModTime"
	case CxlRejReason.DuplicateClOrdID:
		return "DuplicateClOrdID"
	case CxlRejReason.TooLateToCancel:
		return "TooLateToCancel"
	case CxlRejReason.StaleRequest:
		return "StaleRequest"
	case CxlRejReason.RateLimit:
		return "RateLimit"
	case CxlRejReason.Other:
		return "Other"
	default:
		return "Unknown"
	}
}

func ParseCxlRejReason(str string) (s CxlRejReasonEnum, err error) {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	upper := strings.ToUpper(str)
	switch upper {
	case "UNKNOWNORDER":
		s = CxlRejReason.UnknownOrder
	case "BROKER":
		s = CxlRejReason.Broker
	case "ORDERALREADYINPENDINGCANCELORPENDINGREPLACESTATUS":
		s = CxlRejReason.OrderAlreadyInPendingCancelOrPendingReplaceStatus
	case "UNABLETOPROCESSORDERMASSCANCELREQUEST":
		s = CxlRejReason.UnableToProcessOrderMassCancelRequest
	case "ORIGORDMODTIME":
		s = CxlRejReason.OrigOrdModTime
	case "DUPLICATECLORDID":
		s = CxlRejReason.DuplicateClOrdID
	case "TOOLATETOCANCEL":
		s = CxlRejReason.TooLateToCancel
	case "STALEREQUEST":
		s = CxlRejReason.StaleRequest
	case "RATELIMIT":
		s = CxlRejReason.RateLimit
	case "OTHER":
		s = CxlRejReason.Other
	default:
		err = fmt.Errorf("invalid CxlRejReason %s", str)
	}
	return
}

func (e CxlRejReasonEnum) MarshalJSON() ([]byte, error) {
	return []byte("\"" + CxlRejReasonString(e) + "\""), nil
}

func (e *CxlRejReasonEnum) UnmarshalJSON(b []byte) (err error) {
	*e, err = ParseCxlRejReason(string(b))
	return
}

type TimeInForceEnum uint8
type TimeInForceValues struct {
	GoodTillCancel TimeInForceEnum
	Day            TimeInForceEnum
	FillAndKill    TimeInForceEnum
	FillOrKill     TimeInForceEnum
}

var TimeInForce = TimeInForceValues{0, 1, 3, 4}

func TimeInForceString(s TimeInForceEnum) string {
	switch s {
	case TimeInForce.GoodTillCancel:
		return "GoodTillCancel"
	case TimeInForce.Day:
		return "Day"
	case TimeInForce.FillAndKill:
		return "FillAndKill"
	case TimeInForce.FillOrKill:
		return "FillOrKill"
	default:
		return "Unknown"
	}
}

func ParseTimeInForce(str string) (s TimeInForceEnum, err error) {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	upper := strings.ToUpper(str)
	switch upper {
	case "GOODTILLCANCEL":
		s = TimeInForce.GoodTillCancel
	case "DAY":
		s = TimeInForce.Day
	case "FILLANDKILL":
		s = TimeInForce.FillAndKill
	case "FILLORKILL":
		s = TimeInForce.FillOrKill
	default:
		err = fmt.Errorf("invalid TimeInForce %s", str)
	}
	return
}

func (e TimeInForceEnum) MarshalJSON() ([]byte, error) {
	return []byte("\"" + TimeInForceString(e) + "\""), nil
}

func (e *TimeInForceEnum) UnmarshalJSON(b []byte) (err error) {
	*e, err = ParseTimeInForce(string(b))
	return
}
