package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func init() {
	fmt.Println("Plugin: sampleClientPlugin loaded...")
}

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClients interface
var ClientRegisterer = registerer("sampleClientPlugin")

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	fmt.Println("Plugin: RegisterClients Called...")
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
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
		fmt.Println("Plugin: sampleClientPlugin message")
		client := &http.Client{}
		resp, _ := client.Do(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		_, _ = w.Write(bodyBytes)
	}), nil
}
