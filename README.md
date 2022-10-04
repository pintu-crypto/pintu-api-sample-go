# Pintu Websocket Sample

This sample code shows how the Pintu websocket API can be used to build a simple server that sends market orders to
 Pintu to be routed to the exchange and pulls in trades for reporting.

It provides a few pieces of functionality:
1. Sample code that connects and authenticates with the Pintu websocket client and wrappers to encode requests and
 decode responses. See the *client* module.
2. A simple http server that provides an API to send liquidate requests, which are translated to market orders. See
 the *endpoint* module.
3. A simple module that handles liquidation requests and maintains the state of an active request. It also subscribes
 to post trades for reporting purposes. See the *liquidation* module.

## General Order Overview

1. Subscribe to `ExecutionReport` with `StartDate` of the `Timestamp` of the last `ExecutionReport` processed to recover order state.
2. Subscribe to `Trade` with `StartDate` of the `Timestamp` of the last `Trade` processed to recover any missed trades.
3. Maintain a map of (`OrderID`, last received `ExecutionReport`) for open orders.
4. Maintain a set of `ClOrdID` for each pending request.
5. To submit an order, first generate a new `ClOrdID`. Send a `NewOrderSingle` specifying the [params](endpoint/endpoint.go#L109) (`ClOrdID`, `Side`, `Price`, `OrderQty`, `Strategy`, etc). for that request.
6. The response for the request will be received on the `ExecutionReport` stream with the `ClOrdID` being the `ClOrdID` specified on the `NewOrderSingle`. If the order is accepted (`ExecType=PendingNew`), then add this order to the set of open orders.
7. To cancel or modify the open order, generate a new `ClOrdID` and send an `OrderCancelRequest` specifying the the `ClOrdID`, `OrigClOrdID` from the previous successful request, and remaining parameters (see the details below).
8. Fills will be received as `ExecutionReport` with `ExecType=Trade`.

*`OrderCancelRequest` Request Keys
**Key**|**Type**|**Required**|**Description**
:-----:|:-----:|:-----:|:-----:
ClOrdID|string|Y|Unique Client Order ID for this request, usually a UUID.
OrigClOrdID|string|N|Client Order ID of the order to cancel, one of OrderID, OrigClOrdID is required.
OrderID|string|N|Order ID of the order to cancel, one of OrderID, OrigClOrdID is required.
TransactTime|string|Y|An ISO-8601 UTC string of the form 2019-02-13T05:17:32.000000Z.


## Installing

Go 1.19 is required for modules support. To install, checkout from source, then build or run directly.

```shell script
    $ go build ./...
```

## Running

Replace the following variables:
- *ws-address*: pintu websocket address
- *api-key*: your API Key
- *api-secret*: your API Secret

Start the server:

```shell script
    $ go run cmd/main.go --addr <ws-address> --apikey <api-key> --apisecret <api-secret>
```

For example:

```shell script
    $ go run cmd/main.go --addr wss://partner.sandbox.pintu.co.id/ws/v1 --apikey ABCD1234ZXCV --apisecret oin201niasf1920ejalsdknasdnaliw1
```

To request a liquidation of `210 DOGE` to `USDT`, run the following curl command from another window:

```shell script
    $ curl localhost:8085/liquidate?symbol=DOGE-USDT&currency=DOGE&side=Buy&quantity=210
```

If successful, you should see a response like:

```
filled(1 @ 9835)
```

Endpoint parameters:
- `symbol`: a currency pair, like `DOGE-USDT`.
- `currency` : the currency that the quantity is specified in. If not specified, defaults to the base currency for the symbol.
- `side`: `Buy` or `Sell`.
- `quantity`: an quantity of base currency to buy or sell.

## Common Issues

- If you got a response `rejected(Order rejected)`, one of the reasons is the order quantity is less than the minimum size.
