package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/messages"
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
		return errors.New(fmt.Sprintf(messages.ClientConfigFileError, err.Error()))
	}

	if err = json.Unmarshal(f, cfg); err != nil {
		return errors.New(fmt.Sprintf(messages.ClientConfigFIleUnmarshalError, err.Error()))
	}

	return
}

//EndpointIndex get the index of matching endpoint if exists
func (cfg *ClientConfigs) EndpointIndex(endpoint string) (ex bool, ix int) {
	eps := cfg.parsEndpoints()
	return helpers.InSlice(endpoint, eps)
}

//parsEndpoints get all available endpoints in client config file
func (cfg *ClientConfigs) parsEndpoints() (eps []string) {
	for _, k := range *cfg {
		eps = append(eps, k.Endpoint)
	}

	return
}
