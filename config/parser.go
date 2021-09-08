package config

import (
	"encoding/json"
	"errors"
	"os"
)

type ExtraConfig struct {
	Name            string
	Endpoint        string
	SuccessStatuses []int
	FailureStatuses []int
}

//ParseExtraConfig parse the config embedded in extra_config part of krakend config
func ParseExtraConfig(extra map[string]interface{}) (ec ExtraConfig, err error) {
	ec = ExtraConfig{
		Name:     extra["name"].(string),
		Endpoint: extra["endpoint"].(string),
	}

	if len(ec.Name) == 0 {
		return ec, errors.New("wrong plugin name")
	}

	if len(ec.Endpoint) == 0 {
		return ec, errors.New("wrong endpoint setup")
	}

	for _, v := range extra["success_statuses"].([]interface{}) {
		ec.SuccessStatuses = append(ec.SuccessStatuses, int(v.(float64)))
	}

	for _, v := range extra["failure_statuses"].([]interface{}) {
		ec.FailureStatuses = append(ec.FailureStatuses, int(v.(float64)))
	}

	return
}

type Configuration struct {
	Url              string   `json:"url"`
	Method           string   `json:"method"`
	RequireBody      bool     `json:"require_body"`
	AdditionalParams []string `json:"additional_params"`
	Statuses         []int    `json:"statuses"`
}

type Steps struct {
	Success Configuration
	Failure Configuration
}

type ClientConfig struct {
	Endpoint string `json:"endpoint"`
	Steps    []Steps
}

//ParseClientConfig parse the client.json file
func ParseClientConfig(addr string) (cfg []ClientConfig, err error) {
	cfgF, err := os.ReadFile(addr)
	if err != nil {
		return
	}

	if err = json.Unmarshal(cfgF, &cfg); err != nil {
		return
	}

	return
}
