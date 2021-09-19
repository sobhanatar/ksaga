## Krakend SAGA Plugin

This package provides the [Saga](https://microservices.io/patterns/data/saga.html) `proxy /client` plugin
for [KrakenD API Gateway](https://krakend.io/).

#### Plugin config

The file named `client.json` is the config file for this plugin, so it should be with this file when it's deployed to
KrakenD plugins folder next to the Saga plugin.

This file can contain as many transactions as you want, and each transaction have the following properties:

- `endpoint`: the name of the transaction. This parameter and it's value also should be in the configuration of the
  KrakenD configuration file. If the plugin can't find a match an error will be thrown.

```
  "extra_config": {
    "github.com/devopsfaith/krakend/transport/http/client/executor": {
      "name": "sagaClient",
      "endpoint": "confirm_payment"
    }
  }
```

- `register`: the message (key) to the message that will be sent back to user after the successful completion of
  transaction.
- `rollback`: the message (key) to the message that will be sent back to user after the successful rollback.
- `rollback_failed`: the message (key) to the message that will be sent back to user after the failure in rollback.
- `steps`: the steps that is required to complete a transaction:
    - `alias`: a name for the backend service. This can be any name use in logging the process
    - `timeout`: the time in `milliseconds` that http handler will wait for response
    -  `retry_max`: the maximum number of retries
    -  `retry_wait_min`: the minimum time the client wait in `milliseconds`
    -  `retry_wait_max`: the maximum time the client wait in `milliseconds`
    - `statuses`: array of accepted status codes that comeback from backend services. It's important to mention that
      these statuses should just include 2xx statuses.
    - `register`: this part contains the required information for calling backend services:
        - `url`: the url to call the register endpoint of backend service
        - `method`: the method that should be used to call the endpoint
        - `header`: the additional headers required for this service in `{"kay":"value", "key":"value"}` format.
        - `body`: the boolean indicates that this service requires the body from previous service.
    - `rollback`:
        - `url`: same as in register
        - `method`: same as in register
        - `header`: same as in register
        - `body`: same as in register

## Building plugin

```
go build -buildmode=plugin -o saga_client.so client/saga_client.go
```

## Run Krakend

```
./krakend run -c krakend-plugin.json -d
```
### References

- https://developers.redhat.com/blog/2018/10/01/patterns-for-distributed-transactions-within-a-microservices-architecture#what_is_a_distributed_transaction_
- https://pkg.go.dev/github.com/devopsfaith/krakend/transport/http/server/plugin
- https://pkg.go.dev/github.com/devopsfaith/krakend/transport/http/client/plugin
- https://github.com/devopsfaith/krakend-contrib
- https://echorand.me/posts/getting-started-with-golang-plugins/
- https://github.com/devopsfaith/krakend-ce/issues/284#issuecomment-802822679
- https://github.com/inemtsev/KrakendBasicPlugin/blob/master/headerModPlugin.go