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

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer("sagaClient")

type registerer string

func init() {
	fmt.Println(messages.ClientPluginLoad)
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
			cfg  config.ClientConfigs
			resp []byte
			fi   int
			uTID string
		)

		//todo: the address should go into toml
		err = cfg.ParseClient(fmt.Sprintf("./plugins/%s", "client.json"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

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

		resp, fi, err = controllers.ProcessRequests(uTID, req, cfg[ix].Steps)
		if err != nil {
			resp, err = controllers.ProcessRollbackRequests(uTID, req, cfg[ix].Steps, fi)
			if err != nil {
				resp = messages.GenerateMessage(&w, map[string]string{"message": cfg[ix].RollbackFailed})
				_, _ = w.Write(resp)
				return
			}
			resp = messages.GenerateMessage(&w, map[string]string{"message": cfg[ix].Rollback})
			_, _ = w.Write(resp)
			return
		}

		resp = messages.GenerateMessage(&w, map[string]string{"message": cfg[ix].Register})
		_, _ = w.Write(resp)

	}), nil
}
