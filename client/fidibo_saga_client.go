package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/exceptions"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"os"
	"time"
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
		var (
			resp []byte
			cfg  config.ClientConfigs
			fi   int
		)

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
			//resp, _ := json.Marshal(map[string]string{"message": "No matching endpoint found"})
			//w.Header().Add("Content-Type", "application/json")
			//_, _ = w.Write(resp)
			return
		}

		resp, fi, err = ProcessRequests(req, cfg[ix].Steps)
		if err != nil {
			resp, err = ProcessRollbackRequests(req, cfg[ix].Steps, fi)
		}
		_, _ = w.Write(resp)

	}), nil
}

//ProcessRequests process the steps based on config file
func ProcessRequests(req *http.Request, steps []config.Steps) ([]byte, int, error) {
	var (
		resp []byte
		err  error
	)

	sc := len(steps)
	cLog.Println(fmt.Sprintf("Number of services to call: %d", sc))

	for ix, step := range steps {
		resp, err = ProcessRequest(req, step)
		if err != nil {
			cLog.Println(err.Error())
			return resp, ix, err
		}

		if ix < sc-1 {
			req = BuildRequest("success", steps[ix+1], req, resp)
		}
	}

	return resp, 0, err
}

func ProcessRollbackRequests(req *http.Request, steps []config.Steps, ix int) (response []byte, err error) {
	for step := ix - 1; step >= 0; step-- {
		cLog.Println(fmt.Sprintf(exceptions.ClientRollbackError, steps[step].Alias))
		BuildRequest("failure", steps[step], req, nil)
		// todo: tomorrow
	}

	return
}

//ProcessRequest process the first request which is configured in krakend config file
func ProcessRequest(req *http.Request, step config.Steps) (body []byte, err error) {
	cLog.Println(fmt.Sprintf("Calling \"%s\" endpoint...", step.Alias))
	client := &http.Client{
		Timeout: time.Duration(step.Success.Timeout) * time.Millisecond,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(exceptions.ClientBackendCallError, step.Alias, err.Error()))
	}

	ex, _ := helpers.InSlice(resp.StatusCode, step.Failure.Statuses)
	if ex {
		return nil, errors.New(fmt.Sprintf(exceptions.ClientStatusCodeError, step.Alias, resp.StatusCode, err.Error()))
	}

	body, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf(exceptions.ClientReadBodyError, err.Error()))
	}

	return
}

//BuildRequest builds the request based on the step's config
func BuildRequest(state string, step config.Steps, req *http.Request, resp []byte) *http.Request {
	var body *bytes.Buffer

	conf := step.Success
	if state == "failure" {
		conf = step.Failure
	}

	// If the next service declared need for the body pass it in
	if conf.Body {
		body = bytes.NewBuffer(resp)
	}

	request, _ := http.NewRequest(conf.Method, conf.Url, body)
	request.Header = req.Header

	for key, value := range conf.Header {
		request.Header.Add(key, value)
	}

	fmt.Println(request.Header)
	return request
}
