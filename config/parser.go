package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/helpers"
	"newgit.fidibo.com/fidiborearc/krakend/plugins/saga/messages"
	"os"
	"strings"
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

	if err = cfg.validateClient(); err != nil {
		return errors.New(err.Error())
	}

	return
}

func (cfg *ClientConfigs) validateClient() (err error) {
	for k, v := range *cfg {
		if len(strings.Trim(v.Endpoint, " ")) == 0 {
			return errors.New(fmt.Sprintf(messages.ClientConfigEndpointEmptyError, k))
		}

		if len(strings.Trim(v.Register, " ")) == 0 ||
			len(strings.Trim(v.Rollback, " ")) == 0 ||
			len(strings.Trim(v.RollbackFailed, " ")) == 0 {
			return errors.New(fmt.Sprintf(messages.ClientConfigMessagesEmptyError, k))
		}

		for ks, vs := range v.Steps {
			// todo: problem solving area: how can I check only 2xx numbers available here

			if len(strings.Trim(vs.Alias, " ")) == 0 {
				return errors.New(fmt.Sprintf(messages.ClientConfigAliasEmptyError, k, ks))
			}

			if len(strings.Trim(vs.Register.Url, " ")) == 0 {
				return errors.New(fmt.Sprintf(messages.ClientConfigUrlEmptyError, k, ks))
			}

			if len(strings.Trim(vs.Rollback.Url, " ")) == 0 {
				return errors.New(fmt.Sprintf(messages.ClientConfigUrlEmptyError, k, ks))
			}

			if ex, _ := helpers.InSlice(vs.Register.Method, []string{http.MethodGet, http.MethodPost, http.MethodPatch,
				http.MethodPut, http.MethodDelete}); !ex {
				return errors.New(fmt.Sprintf(messages.ClientConfigMethodError, k, ks))
			}

			if ex, _ := helpers.InSlice(vs.Rollback.Method, []string{http.MethodGet, http.MethodPost, http.MethodPatch,
				http.MethodPut, http.MethodDelete}); !ex {
				return errors.New(fmt.Sprintf(messages.ClientConfigMethodError, k, ks))
			}
		}
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
