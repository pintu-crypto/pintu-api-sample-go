# Pintu Golang API Websocket Sample

This sample code shows how the Pintu websocket API can be used to build a simple http server that, upon triggering it's `order` endpoint, sends the market order to Pintu's backend and monitors for order execution status.

The sample consists of the following modules: 

1. **client** - implementation of the API communication protocol (primitives and messages, autentication, etc.)
2. **order** - implementation of the order send and order status receive loop handlers on client side
3. **endpoint** - http web-service implementing the 'order' endpoint
4. **cmd** - the application runner (main function)

## General Order Overview

In order to receive updates from Pintu server, the client application has to:

1. Subscribe to `ExecutionReport` channel with `StartDate` parameter. The `StartDate` could be omitted in case of the very first API connection of the client (it will use time.Now() at the backend), or it could be the `Timestamp` field value ofthe of the last `ExecutionReport` message processed from previous connection.
2. Subscribe to `Trade` channel. In case the `StartDate` field is set, the backend will return all the `Trade`'s happened between the date specified and server's time.Now() value.

To send the order to Pintu's backend the following action should be performed:

1. First generate a new `ClOrdID`. The value can use any pattern to be generated. This sample is using the *uuid.New().String()* call
2. Create and send the `NewOrderSingle` message specifying the [params](endpoint/endpoint.go#L109) (`ClOrdID`, `Side`, `Price`, `OrderQty`, `Strategy`, etc).

To send the order cancellation request to Pintu's backend the following actions should be performed:

1. First generate a new `ClOrdID` for the cancellation 
2. Create and send then `OrderCancelRequest` specifying the newly generated `ClOrdID` and `OrigClOrdID`, which is the ClOrdId of the original order created via the `NewOrderSingle` message

Please see the **endpoint** package for the implementation of the order creation flow

Please refer to the Pintu API documentation for the messages and their required fields format

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

To request a order of `210 DOGE` to `USDT`, run the following curl command from another window:

```shell script
    $ curl localhost:8085/order?symbol=DOGE-USDT&currency=DOGE&side=Buy&quantity=210
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
