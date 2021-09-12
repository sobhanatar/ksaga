package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"os"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer("fidiboSagaClient")

type registerer string

var clogger *log.Logger

func init() {
	clogger = log.New(os.Stderr, "[KRAKEND][CLIENT] ", log.Ldate|log.Ltime)
	clogger.Println("fidiboSagaClient plugin loaded")
}

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	ec, err := config.ParseExtra(extra)
	if err != nil {
		return nil, err
	}
	if ec.Name() != string(r) {
		return nil, fmt.Errorf("plugin: unknown register %s", ec.Name())
	}

	// return the actual handler wrapping or your custom logic, so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			response []byte
			cfg      []config.ClientConfig
		)

		cfg, err = config.ParseClient(fmt.Sprintf("./plugins/%s", "client.json"))
		if err != nil {
			clogger.Println(err.Error())
			return
		}

		ex, ix := config.EndpointIndex(ec.Endpoint(), cfg)
		if !ex {
			//todo: alert somehow
			clogger.Println("No matching endpoint found in SAGA client plugin")
			resp, _ := json.Marshal(map[string]string{"message": "No matching endpoint found"})
			w.Header().Add("Content-Type", "application/json")
			_, _ = w.Write(resp)
			return
		}

		response = ProcessSteps(req, cfg[ix].Steps)
		_, _ = w.Write(response)

	}), nil
}

//ProcessSteps process the steps based on config file
func ProcessSteps(req *http.Request, steps []config.Steps) (resp []byte) {
	sc := len(steps)
	clogger.Println(fmt.Sprintf("Number of services to call: %d", sc))
	response, err := ProcessInitialRequest(req, steps[0])
	clogger.Println(fmt.Sprintf("Response status from %s: %d", steps[0].Alias, response.StatusCode))
	resp, _ = io.ReadAll(response.Body)
	if err != nil {
		clogger.Println(fmt.Sprintf("Call Error: %s", err.Error()))
		return
	}

	return
	//for _, step := range steps {
	//
	//}

	return
}

//ProcessInitialRequest process the first request which is configured in krakend config file
func ProcessInitialRequest(initReq *http.Request, step config.Steps) (resp *http.Response, err error) {
	clogger.Println(fmt.Sprintf("Calling %s...", step.Alias))
	client := &http.Client{}
	resp, err = client.Do(initReq)
	if err != nil {
		clogger.Println("Call has failed. No compensation has called for initiating call")
		return
	}

	ex, _ := helpers.InSlice(resp.StatusCode, step.Failure.Statuses)
	if ex {
		return resp, errors.New(fmt.Sprintf("Call Failed with status %d. No compensation has called for initiating call", resp.StatusCode))
	}

	return
}
