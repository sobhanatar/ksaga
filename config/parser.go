package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"os"
)

//ParseExtra parse extra_config key value of plugin in krakend config file
func (ec *ExtraConfig) ParseExtra(extra map[string]interface{}) (err error) {
	ec.SetName(extra["name"].(string))
	ec.SetEndpoint(extra["endpoint"].(string))

	if len(ec.Name()) == 0 {
		return errors.New("wrong plugin name")
	}

	if len(ec.Endpoint()) == 0 {
		return errors.New("wrong endpoint setup")
	}

	return
}

//ParseClient parse the client.json file
func (cfg *ClientConfigs) ParseClient(addr string) (err error) {
	f, err := os.ReadFile(addr)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading client.json file: %s", err.Error()))
	}

	if err = json.Unmarshal(f, &cfg); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshaling client.json file: %s", err.Error()))
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
