package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	Extra  = "extra"
	Status = "status"
	HKey   = "Content-Type"
	HVal   = "application/json"
	CfgAdr = "./plugins/saga_client_settings.json"
)

type registerer string

var (
	cfg  config.SagaClientConfig
	cLog *logrus.Logger
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer(PName)

func init() {
	lvl, _ := logrus.ParseLevel(cfg.LogLevel)
	cLog = logs.GetInstance(lvl)

	err := cfg.ParseClient(CfgAdr)
	if err != nil {
		cLog.WithFields(map[string]interface{}{Msg: err.Error()}).Panic(messages.ClientPluginLoadError)
		logs.Logs(logrus.PanicLevel, fmt.Sprintf(messages.ClientPluginLoadError, err.Error()))
		return
	}

	logs.Logs(logrus.InfoLevel, fmt.Sprintf(messages.ClientPluginLoad, ClientRegisterer))
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

		cLog.Panic(logs.GenerateLog(map[string]interface{}{Msg: "my message", UTID: uTID, Extra: "my extra", "severity": logrus.ErrorLevel}))
		fmt.Println("error created")
		ex, ix := cfg.EndpointIndex(ec.Endpoint())
		if !ex {
			m := fmt.Sprintf(messages.ClientEndpointNotFoundError, ec.Endpoint())
			cLog.WithFields(map[string]interface{}{UTID: uTID}).Error(m)
			logs.Logs(logrus.ErrorLevel, m)
			return
		}

		logs.Logs(logrus.InfoLevel, fmt.Sprintf(messages.ClientUniversalTransactionID, uTID))

		ep := cfg.Endpoints[ix]
		fs, err = controllers.ProcessRequests(uTID, req, ep.Steps)
		if err != nil {
			err = controllers.ProcessRollbackRequests(uTID, req, ep.Steps, fs-1)
			if err != nil {
				cLog.WithFields(map[string]interface{}{UTID: uTID, Extra: ep.Steps[fs].Alias}).Panic(ep.RollbackFailed)
				generateResponse(&w, map[string]interface{}{Status: http.StatusUnprocessableEntity, Msg: ep.RollbackFailed})
				return
			}

			cLog.WithFields(map[string]interface{}{UTID: uTID, Extra: ep.Steps[fs].Alias}).Error(ep.Rollback)
			generateResponse(&w, map[string]interface{}{Status: http.StatusUnprocessableEntity, Msg: ep.Rollback})
			return
		}

		cLog.WithFields(map[string]interface{}{UTID: uTID, Extra: ep.Endpoint}).Info(ep.Register)
		generateResponse(&w, map[string]interface{}{Status: http.StatusOK, Msg: ep.Register})
	}), nil
}

func generateResponse(w *http.ResponseWriter, m map[string]interface{}) {
	(*w).Header().Add(HKey, HVal)
	_, err := (*w).Write(messages.Generate(m))
	if err != nil {
		logs.Logs(logrus.ErrorLevel, fmt.Sprintf(messages.ClientResponseWriterError, err.Error()))
	}
}
