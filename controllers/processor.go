package controllers

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/config"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/exceptions"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/logs"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/messages"
	"time"
)

const (
	Register = "register"
	Rollback = "rollback"
)

//ProcessRequests call backend services based on th defined services on go
func ProcessRequests(uTID string, req *http.Request, steps []config.Steps) (int, error) {
	var (
		resp []byte
		err  error
	)

	sc := len(steps)
	logs.Logs(logrus.InfoLevel, fmt.Sprintf(messages.ClientBackendServices, sc))

	for ix, step := range steps {
		logs.Logs(logrus.InfoLevel, fmt.Sprintf(messages.ClientServiceCall, step.Alias, req.URL.String()))
		req = buildRequest(Register, steps[ix], req, resp)

		resp, err = processRequest(uTID, req, step)
		if err != nil {
			logs.Logs(logrus.ErrorLevel, err.Error())
			return ix, err
		}

		if ix < sc-1 {
			req = buildRequest(Register, steps[ix+1], req, resp)
		}
	}

	return 0, err
}

//ProcessRollbackRequests call backend rollback requests based on th defined services on go
func ProcessRollbackRequests(uTID string, req *http.Request, steps []config.Steps, fs int) (err error) {
	for step := fs; step >= 0; step-- {
		req = buildRequest(Rollback, steps[step], req, nil)
		logs.Logs(logrus.InfoLevel, fmt.Sprintf(messages.ClientRollbackError, steps[step].Alias, req.URL.String()))

		_, err = processRequest(uTID, req, steps[step])
		if err != nil {
			logs.Logs(logrus.ErrorLevel, err.Error())
			return
		}
	}

	return
}

func processRequest(uTID string, req *http.Request, step config.Steps) (body []byte, err error) {
	rc := retryablehttp.NewClient()
	rc.RetryMax = step.RetryMax
	rc.RetryWaitMin = time.Duration(step.RetryWaitMin) * time.Millisecond
	rc.RetryWaitMax = time.Duration(step.RetryWaitMax) * time.Millisecond

	client := rc.StandardClient()
	client.Timeout = time.Duration(step.Timeout) * time.Millisecond

	// Add global transaction id
	req.Header.Set("Universal-Transaction-ID", uTID)

	resp, err := client.Do(req)
	if err != nil {
		return body, exceptions.BackendCallError(step.Alias, err.Error())
	}
	defer resp.Body.Close()

	ex, _ := helpers.InSlice(resp.StatusCode, step.Statuses)
	if !ex {
		return body, exceptions.StatusError(step.Alias, resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, exceptions.CloseBodyError(step.Alias, err.Error())
	}

	return
}

func buildRequest(state string, step config.Steps, req *http.Request, resp []byte) *http.Request {
	var body = new(bytes.Buffer)

	conf := step.Register
	if state == Rollback {
		conf = step.Rollback
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
