package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/controllers"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/messages"
)

type registerer string

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer("sagaClient")

var (
	cfg    config.ClientConfig
	cfgAdr = fmt.Sprintf("./plugins/%s", "saga_client.json")
)

func init() {
	err := cfg.ParseClient(cfgAdr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(fmt.Sprintf(messages.ClientPluginLoad, ClientRegisterer))
}

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	var (
		ec config.ExtraConfig
	)

	err := ec.ParseExtra(extra)
	if err != nil {
		return nil, err
	}

	if ec.Name() != string(r) {
		return nil, fmt.Errorf("plugin: unknown register %s", ec.Name())
	}

	// return the actual handler wrapping or your custom logic, so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			resp []byte
			fi   int
			uTID string
		)

		//todo: the address should go into toml
		ex, ix := cfg.EndpointIndex(ec.Endpoint())
		if !ex {
			/*
			 * todo: alerting / registering event in sentry, kafka, ...
			 */
			fmt.Println(fmt.Sprintf(messages.ClientEndpointNotFoundError, ec.Endpoint()))
			return
		}

		uTID = uuid.New().String()
		fmt.Println(fmt.Sprintf(messages.CallServiceGlobalTransactionID, uTID))

		ep := cfg.Endpoints[ix]
		resp, fi, err = controllers.ProcessRequests(uTID, req, ep.Steps)
		if err != nil {
			resp, err = controllers.ProcessRollbackRequests(uTID, req, ep.Steps, fi)
			if err != nil {
				resp = messages.GenerateMessage(&w, map[string]interface{}{"status": 422, "message": ep.RollbackFailed})
				_, _ = w.Write(resp)
				return
			}
			resp = messages.GenerateMessage(&w, map[string]interface{}{"status": 422, "message": ep.Rollback})
			_, _ = w.Write(resp)
			return
		}

		resp = messages.GenerateMessage(&w, map[string]interface{}{"status": 200, "message": ep.Register})
		_, _ = w.Write(resp)

	}), nil
}
