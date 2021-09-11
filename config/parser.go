package config

import (
	"encoding/json"
	"errors"
	"os"
)

type ExtraConfiguration struct {
	Message  string `json:"message"`
	Statuses []int  `json:"statuses"`
}

type ExtraConfig struct {
	Name     string
	Endpoint string
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

	return
}

type Configuration struct {
	Url              string   `json:"url"`
	Method           string   `json:"method"`
	Message          string   `json:"message"`
	RequireBody      bool     `json:"require_body"`
	AdditionalParams []string `json:"additional_params"`
	Statuses         []int    `json:"statuses"`
}

type Steps struct {
	Alias   string
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

//GetEndpoints get all available endpoints in client config file
func GetEndpoints(cfg []ClientConfig) (eps []string) {
	for _, k := range cfg {
		eps = append(eps, k.Endpoint)
	}

	return
}
