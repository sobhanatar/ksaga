package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"os"
)

//ParseExtra parse extra_config key value of plugin in krakend config file
func ParseExtra(extra map[string]interface{}) (ec ExtraConfig, err error) {
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

//ParseClient parse the client.json file
func ParseClient(addr string) (cfg []ClientConfig, err error) {
	f, err := os.ReadFile(addr)
	if err != nil {
		return cfg, errors.New(fmt.Sprintf("Error reading client.json file: %s", err.Error()))
	}

	if err = json.Unmarshal(f, &cfg); err != nil {
		return cfg, errors.New(fmt.Sprintf("Error unmarshaling client.json file: %s", err.Error()))
	}

	return
}

//EndpointIndex get the index of matching endpoint if exists
func EndpointIndex(endpoint string, cfg []ClientConfig) (ex bool, ix int) {
	eps := parsEndpoints(cfg)
	return helpers.InSlice(endpoint, eps)
}

//parsEndpoints get all available endpoints in client config file
func parsEndpoints(cfg []ClientConfig) (eps []string) {
	for _, k := range cfg {
		eps = append(eps, k.Endpoint)
	}

	return
}
