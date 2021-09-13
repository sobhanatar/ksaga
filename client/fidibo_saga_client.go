package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/controllers"
	"os"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer("fidiboSagaClient")

type registerer string

var cLog *log.Logger

func init() {
	cLog = log.New(os.Stderr, "[KRAKEND][CLIENT] ", log.Ldate|log.Ltime)
	cLog.Println("fidiboSagaClient plugin loaded")
}

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	var ec config.ExtraConfig

	err := ec.ParseExtra(extra)
	if err != nil {
		return nil, err
	}

	if ec.Name() != string(r) {
		return nil, fmt.Errorf("plugin: unknown register %s", ec.Name())
	}

	// return the actual handler wrapping or your custom logic, so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var cfg config.ClientConfigs

		err = cfg.ParseClient(fmt.Sprintf("./plugins/%s", "client.json"))
		if err != nil {
			cLog.Println(err.Error())
			return
		}

		ex, ix := cfg.EndpointIndex(ec.Endpoint())
		if !ex {
			/*
			 * todo: alerting / registering event in sentry, kafka, ...
			 */
			cLog.Println("No matching endpoint found in SAGA client plugin")
			return
		}

		resp, fi, err := controllers.ProcessRequests(req, cfg[ix].Steps)
		if err != nil {
			resp, err = controllers.ProcessRollbackRequests(req, cfg[ix].Steps, fi)
		}
		_, _ = w.Write(resp)

	}), nil
}
