## Krakend SAGA Plugin

This package provides the [Saga](https://microservices.io/patterns/data/saga.html) `proxy /client` plugin
for [KrakenD API Gateway](https://krakend.io/).

#### Plugin configuration

`saga_client.json` is the config file for this plugin, and it should be in the same folder as the plugin exists. This
file can contain as many transactions as a system need.

The fields of the configuration files as follows:

- `log_level`: the level of debug application will log in file
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
        - `statuses`: array of accepted status codes coming back from backend service. It's important to mention that,
          no matter the service is moving forward or rolling back, these statuses will be used to move forward to
          next/previous service. It's recommended to use `2xx` statuses.
        - `register`: this part contains the required information for calling backend services:
            - `url`: the URL to call the registering endpoint of backend service
            - `method`: the method that should be used to call the endpoint.
            - `header`: the additional headers required for this service in `{"kay":"value", "key":"value"}` format.
            - `body`: the boolean indicates that this service requires the body from the previous service.
        - `rollback`:
            - `url`: same as in register
            - `method`: same as in register
            - `header`: same as in register
            - `body`: same as in register

#### Plugin Logging

This package use [logrus](https://github.com/sirupsen/logrus) for logging rollbacks and rollback failures to file. The
file name follows the pattern of `saga-client-plugin-{date}.log,` and is in JSON format so that it can be consumed by
services like `logstash.`

#### Plugin Request Consistency

As we live in the real world, nothing is guaranteed. So there is always the possibility of things going wrong, and
calling backend services is not an exception. When we send requests to a series of backend services via the saga plugin,
if any of those services do not respond, we need to roll back the transaction and make other services execute their
rollback procedure. For consistency, before sending a request to any backend services, a unique id is generated and
placed in a header named `Universal-Transaction-ID.`

This value will be in every request's header, and all the backend services read it and store it as a reference to the
transaction, so in case of rollback, they can find the data related to it.

[uuid](https://github.com/google/uuid) package is used for generating unique ids.

### Building plugin

Compile the plugin with `go build -buildmode=plugin -o yourplugin.so,` and then reference them in the KrakenD
configuration file. For instance:

```
//backend part of endpoints
"backend": [
        {
          "method": "POST",
          "encoding": "json",
          "host": [
            "http://localhost:8080"
          ],
          "url_pattern": "/api/krakend/payment/register",
          "extra_config": {
            "github.com/devopsfaith/krakend/transport/http/client/executor": {
              "name": "sagaClient",
              "endpoint": "confirm_payment"
            }
          }
        }
      ]
      //rest of the config
```

### Tests

### Run Krakend

To use the SAGA plugin with KrakenD, check the krakend-plugin.json, which is a blueprint for injecting a client plugin.
