package main

import (
	"context"
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
	ec, err := config.ParseExtraConfig(extra)
	if err != nil {
		return nil, err
	}
	if ec.Name != string(r) {
		return nil, fmt.Errorf("plugin: unknown register %s", ec.Name)
	}

	// return the actual handler wrapping or your custom logic, so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response []byte

		// create a http client, do the request, parse/manipulate the client response and write it back to your responseWriter.
		cfg, err := config.ParseClientConfig(fmt.Sprintf("./plugins/%s", "client.json"))
		if err != nil {
			clogger.Println(err.Error())
			return
		}

		//if !ok {
		//	clogger.Println("No matching endpoint found")
		//	return
		//}

		for _, endpoint := range cfg {
			clogger.Println(req.URL.Path, req.URL.String())
			if endpoint.Endpoint == ec.Endpoint {
				fmt.Println(": ", req.URL.String())
				response = ProcessSteps(req, endpoint.Steps)
			}
		}

		_, _ = w.Write(response)

	}), nil
}

func ProcessSteps(req *http.Request, steps []config.Steps) (resp []byte) {
	sc := len(steps)
	clogger.Println(fmt.Sprintf("Number of Services to call: %d", sc))
	clogger.Println(fmt.Sprintf("Calling %s...", req.URL.String()))

	response, err := ProcessInitialRequest(req, steps[0])
	if err != nil {
		clogger.Println(fmt.Sprintf("Call Error: %s", err.Error()))
		return
	}

	resp, _ = io.ReadAll(response.Body)
	return
	//for _, step := range steps {
	//
	//}

	return
}

func ProcessInitialRequest(initReq *http.Request, step config.Steps) (resp *http.Response, err error) {
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
