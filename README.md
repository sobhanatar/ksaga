## Krakend SAGA Plugin

This package provides the [Saga](https://microservices.io/patterns/data/saga.html) `proxy /client` plugin
for [KrakenD API Gateway](https://krakend.io/).

#### Plugin configuration

`saga_client.json` is the config file for this plugin, and it should be in the same folder as the plugin exists. This
file can contain as many transactions as a system need.

The fields of the configuration files as follows:

- `endpoints`: the array of transactions that can be handled by this plugin. Each endpoint has the following structure:
    - `endpoint`: the name of the endpoint. This parameter and its value also should be the same as `endpoint` name in
      the `extra_config` part of krakenD configuration file. If the plugin cannot find a match, an error is thrown.
        ```
        "extra_config": {
          "github.com/devopsfaith/krakend/transport/http/client/executor": {
             "name": "sagaClient",
             "endpoint": "confirm_payment"
          }
        }
        ```
    - `register`: the message (key) to the message that will be sent back to the user after completing the transaction.
    - `rollback`: the message (key) to the message that will be sent back to the user after the successful rollback.
    - `rollback_failed`: the message (key) to the message that will be sent back to the user after the failure in
      rollback.
    - `steps`: the steps that are required to complete a transaction:
        - `alias`: a name for the backend service.
        - `timeout`: the time in `milliseconds` that HTTP handler will wait for the response
        - `retry_max`: the maximum number of retries
        - `retry_wait_min`: the minimum time the client wait in `milliseconds.`
        - `retry_wait_max`: the maximum time the client wait in `milliseconds.`
        - `statuses`: array of accepted status codes that come back from backend services. It is important to mention
          that these statuses should only include 2xx statuses.
        - `register`: this part contains the required information for calling backend services:
            - `url`: the URL to call the registering endpoint of backend service
            - `method`: the method that should be used to call the endpoint
            - `header`: the additional headers required for this service in `{"kay":"value", "key":"value"}` format.
            - `body`: the boolean indicates that this service requires the body from the previous service.
        - `rollback`:
            - `url`: same as in register
            - `method`: same as in register
            - `header`: same as in register
            - `body`: same as in register

### Building plugin

```
go build -buildmode=plugin -o saga_client.so client/saga_client.go
```

### Run Krakend

```
./krakend run -c krakend-plugin.json -d
```
