package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var HandlerRegisterer = registerer("fidiboSagaHandler")

type Message struct {
	Phone          int64   `json:"phone"`
	OrderId        string  `json:"order_id"`
	Amount         float32 `json:"amount"`
	Description    string  `json:"description"`
	MobileOperator string  `json:"mobile_operator"`
	Language       string  `json:"language"`
	TestMode       int8    `json:"testmode"`
}

type registerer string

func init() {
	fmt.Println("Handler: fidiboSagaHandler plugin loaded!!!")
}

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Handler: fidiboSagaHandler message")
		fmt.Println("Handler: " + req.URL.String())

		//response, _ := io.ReadAll(req.Body)
		//fmt.Println("Handler: " + string(response))

		_, _ = w.Write([]byte("HAPAL"))
		handler.ServeHTTP(w, req)
	}), nil
}
