package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer("fidiboSagaClient")

type registerer string

func init() {
	fmt.Println("Client: fidiboSagaClient plugin loaded!!!")
}

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	fmt.Println("Client: ", extra)
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("plugin: unknown register %s", name)
	}

	// return the actual handler wrapping or your custom logic, so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// create a http client, do the request, parse/manipulate the client response and write it back to your responseWriter.
		fmt.Println("Client: fidiboSagaClient message")
		fmt.Println("Client: " + req.URL.String())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Client Error: ", err.Error())
			return
		}
		fmt.Println("Client: After Request: " + resp.Header.Get("Content-Type"))
		response, _ := io.ReadAll(resp.Body)
		_, _ = w.Write(response)
	}), nil
}
