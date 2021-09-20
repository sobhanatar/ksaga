package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/controllers"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/logs"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/messages"
)

const (
	HKey   = "Content-Type"
	HVal   = "application/json"
	CfgAdr = "./plugins/saga_client.json"
)

type registerer string

var cfg config.ClientConfig

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer("sagaClient")

func init() {
	err := cfg.ParseClient(CfgAdr)
	if err != nil {
		logs.Log(logs.ERROR, messages.ClientPluginLoadError)
		fmt.Println(logs.ERROR, err.Error())
		return
	}

	logs.Log(logs.INFO, fmt.Sprintf(messages.ClientPluginLoad, ClientRegisterer))
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
			fi   int
			uTID string
		)

		//todo: the address should go into toml
		ex, ix := cfg.EndpointIndex(ec.Endpoint())
		if !ex {
			/*
			 * todo: alerting / registering event in sentry, kafka, ...
			 */
			logs.Log(logs.ERROR, fmt.Sprintf(messages.ClientEndpointNotFoundError, ec.Endpoint()))
			return
		}

		uTID = uuid.New().String()
		logs.Log(logs.INFO, fmt.Sprintf(messages.CallServiceGlobalTransactionID, uTID))

		ep := cfg.Endpoints[ix]
		fi, err = controllers.ProcessRequests(uTID, req, ep.Steps)
		if err != nil {
			err = controllers.ProcessRollbackRequests(uTID, req, ep.Steps, fi)
			if err != nil {
				w.Header().Add(HKey, HVal)
				_, err = w.Write(messages.Generate(map[string]interface{}{"status": 422, "message": ep.RollbackFailed}))
				if err != nil {
					logs.Log(logs.ERROR, fmt.Sprintf(messages.ClientResponseWriterError, err.Error()))
				}
				return
			}
			w.Header().Add(HKey, HVal)
			_, err = w.Write(messages.Generate(map[string]interface{}{"status": 422, "message": ep.Rollback}))
			if err != nil {
				logs.Log(logs.ERROR, fmt.Sprintf(messages.ClientResponseWriterError, err.Error()))
			}
			return
		}

		w.Header().Add(HKey, HVal)
		_, err = w.Write(messages.Generate(map[string]interface{}{"status": 200, "message": ep.Register}))
		if err != nil {
			logs.Log(logs.ERROR, fmt.Sprintf(messages.ClientResponseWriterError, err.Error()))
		}

	}), nil
}
