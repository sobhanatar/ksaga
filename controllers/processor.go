package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/exceptions"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"time"
)

//ProcessRequests process the steps based on config file
func ProcessRequests(req *http.Request, steps []config.Steps) ([]byte, int, error) {
	var (
		resp []byte
		err  error
	)

	sc := len(steps)
	fmt.Println(fmt.Sprintf("Number of services to call: %d", sc))

	for ix, step := range steps {
		resp, err = processRequest(req, step)
		if err != nil {
			fmt.Println(err.Error())
			return resp, ix, err
		}

		if ix < sc-1 {
			req = buildRequest("success", steps[ix+1], req, resp)
		}
	}

	return resp, 0, err
}

func ProcessRollbackRequests(req *http.Request, steps []config.Steps, ix int) (response []byte, err error) {
	for step := ix - 1; step >= 0; step-- {
		fmt.Println(fmt.Sprintf(exceptions.ClientRollbackError, steps[step].Alias))
		buildRequest("failure", steps[step], req, nil)
		// todo: tomorrow
	}

	return
}

//processRequest process the first request which is configured in krakend config file
func processRequest(req *http.Request, step config.Steps) (body []byte, err error) {
	fmt.Println(fmt.Sprintf("Calling \"%s\" endpoint...", step.Alias))
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

//buildRequest builds the request based on the step's config
func buildRequest(state string, step config.Steps, req *http.Request, resp []byte) *http.Request {
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
