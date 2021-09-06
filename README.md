## SAGA Plugin

To create plugin run following command.

```
go build -buildmode=plugin -o fidibo_saga_client.so client/fidibo_saga_client.go
go build -buildmode=plugin -o fidibo_saga_handler.so handler/fidibo_saga_handler.go
```

### References
- https://developers.redhat.com/blog/2018/10/01/patterns-for-distributed-transactions-within-a-microservices-architecture#what_is_a_distributed_transaction_
- https://pkg.go.dev/github.com/devopsfaith/krakend/transport/http/server/plugin
- https://pkg.go.dev/github.com/devopsfaith/krakend/transport/http/client/plugin
- https://github.com/devopsfaith/krakend-contrib
- https://echorand.me/posts/getting-started-with-golang-plugins/
- https://github.com/devopsfaith/krakend-ce/issues/284#issuecomment-802822679
- https://github.com/inemtsev/KrakendBasicPlugin/blob/master/headerModPlugin.go