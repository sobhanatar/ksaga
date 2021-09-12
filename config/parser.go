package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

//Extra parse extra_config key value of plugin in krakend config file
func Extra(extra map[string]interface{}) (ec ExtraConfig, err error) {
	ec.SetName(extra["name"].(string))
	ec.SetEndpoint(extra["endpoint"].(string))

	if len(ec.Name()) == 0 {
		return ec, errors.New("wrong plugin name")
	}

	if len(ec.Endpoint()) == 0 {
		return ec, errors.New("wrong endpoint setup")
	}

	return
}

//Client parse the client.json file
func Client(addr string) (cfg []ClientConfig, err error) {
	cfgF, err := os.ReadFile(addr)
	if err != nil {
		return cfg, errors.New(fmt.Sprintf("Error reading client.json file: %s", err.Error()))
	}

	if err = json.Unmarshal(cfgF, &cfg); err != nil {
		return cfg, errors.New(fmt.Sprintf("Error unmarshaling client.json file: %s", err.Error()))
	}

	return
}

//Endpoints get all available endpoints in client config file
func Endpoints(cfg []ClientConfig) (eps []string) {
	for _, k := range cfg {
		eps = append(eps, k.Endpoint)
	}

	return
}
