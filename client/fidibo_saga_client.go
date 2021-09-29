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
	PName  = "sagaClient"
	UTID   = "utid"
	Msg    = "message"
	Status = "status"
	HKey   = "Content-Type"
	HVal   = "application/json"
	CfgAdr = "./plugins/saga_client_settings.json"
)

type registerer string

var cfg config.SagaClientConfig

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer(PName)

func init() {
	err := cfg.ParseClient(CfgAdr)
	if err != nil {
		logs.LogF(logs.Panic, messages.ClientPluginLoadError, map[string]interface{}{Msg: err.Error()})
		logs.Log(logs.Panic, messages.ClientPluginLoadError)
		fmt.Println(logs.Panic, err.Error())
		return
	}

	logs.Log(logs.Info, fmt.Sprintf(messages.ClientPluginLoad, ClientRegisterer))
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
		return nil, fmt.Errorf(messages.ClientPluginNameError, ec.Name())
	}

	// return the actual handler wrapping or your custom logic, so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			fs   int
			uTID string
		)

		uTID = uuid.New().String()

		ex, ix := cfg.EndpointIndex(ec.Endpoint())
		if !ex {
			m := fmt.Sprintf(messages.ClientEndpointNotFoundError, ec.Endpoint())
			logs.LogF(logs.Error, m, map[string]interface{}{UTID: uTID})
			logs.Log(logs.Error, m)
			return
		}

		logs.Log(logs.Info, fmt.Sprintf(messages.ClientUniversalTransactionID, uTID))

		ep := cfg.Endpoints[ix]
		fs, err = controllers.ProcessRequests(uTID, req, ep.Steps)
		if err != nil {
			err = controllers.ProcessRollbackRequests(uTID, req, ep.Steps, fs-1)
			if err != nil {
				logs.LogF(logs.Panic, ep.RollbackFailed, map[string]interface{}{UTID: uTID, "extra": ep.Steps[fs]})
				generateResponse(&w, map[string]interface{}{Status: http.StatusUnprocessableEntity, Msg: ep.RollbackFailed})
				return
			}

			logs.LogF(logs.Error, ep.Rollback, map[string]interface{}{UTID: uTID, "extra": ep.Steps[fs]})
			generateResponse(&w, map[string]interface{}{Status: http.StatusUnprocessableEntity, Msg: ep.Rollback})
			return
		}

		logs.LogF(logs.Info, ep.Register, map[string]interface{}{UTID: uTID})
		generateResponse(&w, map[string]interface{}{Status: http.StatusOK, Msg: ep.Register})
	}), nil
}

func generateResponse(w *http.ResponseWriter, m map[string]interface{}) {
	(*w).Header().Add(HKey, HVal)
	_, err := (*w).Write(messages.Generate(m))
	if err != nil {
		logs.Log(logs.Error, fmt.Sprintf(messages.ClientResponseWriterError, err.Error()))
	}
}
