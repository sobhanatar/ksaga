package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/messages"
	"time"
)

const (
	Success = "success"
	Failure = "failure"
)

//ProcessRequests call backend services based on th defined services on go
func ProcessRequests(uTID string, req *http.Request, steps []config.Steps) ([]byte, int, error) {
	var (
		resp []byte
		err  error
	)

	sc := len(steps)
	fmt.Println(fmt.Sprintf(messages.CallNumberOfBackendService, sc))

	for ix, step := range steps {
		fmt.Println(fmt.Sprintf(messages.ClientServiceCall, step.Alias, req.URL.String()))
		req = buildRequest(Success, steps[ix], req, resp)

		resp, err = processRequest(uTID, req, step)
		if err != nil {
			fmt.Println(err.Error())
			return resp, ix, err
		}

		if ix < sc-1 {
			req = buildRequest(Success, steps[ix+1], req, resp)
		}
	}

	return resp, 0, err
}

//ProcessRollbackRequests call backend rollback requests based on th defined services on go
func ProcessRollbackRequests(uTID string, req *http.Request, steps []config.Steps, ix int) (response []byte, err error) {
	if ix == 0 { //if the first transaction failed, the failure callback must be called
		req = buildRequest(Failure, steps[0], req, response)
		fmt.Println(fmt.Sprintf(messages.ClientRollbackError, steps[0].Alias, req.URL.String()))

		_, err = processRequest(uTID, req, steps[0])
		if err != nil {
			/*
			 * Todo: alert, write in kafka, etc
			 */
			fmt.Println(err.Error())
			return
		}

		return
	}

	for step := ix - 1; step >= 0; step-- {
		req = buildRequest(Failure, steps[step], req, nil)
		fmt.Println(fmt.Sprintf(messages.ClientRollbackError, steps[step].Alias, req.URL.String()))

		_, err = processRequest(uTID, req, steps[step])
		if err != nil {
			/*
			 * Todo: alert, write in kafka, etc
			 */
			fmt.Println(err.Error())
			return
		}
	}

	return
}

func processRequest(uTID string, req *http.Request, step config.Steps) (body []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(step.Success.Timeout) * time.Millisecond,
	}

	// Add global transaction id
	req.Header.Set("Universal-Transaction-ID", uTID)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return body, messages.BackendCallError(step.Alias, err.Error())
	}

	ex, _ := helpers.InSlice(resp.StatusCode, step.Statuses)
	if !ex {
		return body, messages.StatusError(step.Alias, resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, messages.CloseBodyError(step.Alias, err.Error())
	}

	return
}

func buildRequest(state string, step config.Steps, req *http.Request, resp []byte) *http.Request {
	var body = new(bytes.Buffer)

	conf := step.Success
	if state == Failure {
		conf = step.Failure
	}

	// If the next service declared need for the body and body is not nil then pass it in
	if conf.Body && len(resp) != 0 {
		body = bytes.NewBuffer(resp)
	}

	request, _ := http.NewRequest(conf.Method, conf.Url, body)
	request.Header = req.Header

	for key, value := range conf.Header {
		request.Header.Set(key, value)
	}

	return request
}
