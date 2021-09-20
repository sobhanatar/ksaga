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

//ParseClient parse the saga_client.json file
func (cfg *ClientConfig) ParseClient(addr string) (err error) {
	f, err := os.ReadFile(addr)
	if err != nil {
		return errors.New(fmt.Sprintf(messages.ClientConfigFileError, err.Error()))
	}

	if err = json.Unmarshal(f, cfg); err != nil {
		return errors.New(fmt.Sprintf(messages.ClientConfigFIleUnmarshalError, err.Error()))
	}

	if errs := cfg.validate(); len(errs) > 0 {
		errL := strings.Join(errs, "\n")
		return errors.New(errL)
	}

	return
}

func (cfg *ClientConfig) validate() (err []string) {
	statuses := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusAlreadyReported,
		http.StatusIMUsed,
		http.StatusMultiStatus,
		http.StatusPaymentRequired,
		http.StatusNonAuthoritativeInfo,
	}
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete}

	for k, v := range (*cfg).Endpoints {
		if len(strings.Trim(v.Endpoint, " ")) == 0 {
			err = append(err, fmt.Sprintf(messages.ClientConfigEndpointEmptyError, k+1))
		}

		if len(strings.Trim(v.Register, " ")) == 0 ||
			len(strings.Trim(v.Rollback, " ")) == 0 ||
			len(strings.Trim(v.RollbackFailed, " ")) == 0 {
			err = append(err, fmt.Sprintf(messages.ClientConfigMessagesEmptyError, k+1))
		}

		for ks, vs := range v.Steps {
			if len(strings.Trim(vs.Alias, " ")) == 0 {
				err = append(err, fmt.Sprintf(messages.ClientConfigAliasEmptyError, k+1, ks+1))
			}

			if vs.Timeout <= 0 {
				err = append(err, fmt.Sprintf(messages.ClientConfigTimeoutError, k+1, ks+1))
			}

			if vs.RetryMax < 0 {
				err = append(err, fmt.Sprintf(messages.ClientConfigRetryError, k+1, ks+1))
			}

			if vs.RetryWaitMax <= 0 || vs.RetryWaitMin <= 0 {
				err = append(err, fmt.Sprintf(messages.ClientConfigRetryWaitError, k+1, ks+1))
			}

			if len(strings.Trim(vs.Register.Url, " ")) == 0 {
				err = append(err, fmt.Sprintf(messages.ClientConfigUrlEmptyError, k+1, ks+1))
			}

			if len(strings.Trim(vs.Rollback.Url, " ")) == 0 {
				err = append(err, fmt.Sprintf(messages.ClientConfigUrlEmptyError, k+1, ks+1))
			}

			for _, st := range vs.Statuses {
				if ex, _ := helpers.InSlice(st, statuses); !ex {
					err = append(err, fmt.Sprintf(messages.ClientConfigStatusError, k+1, ks+1))
				}
			}

			if ex, _ := helpers.InSlice(vs.Register.Method, methods); !ex {
				err = append(err, fmt.Sprintf(messages.ClientConfigMethodError, k+1, ks+1))
			}

			if ex, _ := helpers.InSlice(vs.Rollback.Method, methods); !ex {
				err = append(err, fmt.Sprintf(messages.ClientConfigMethodError, k+1, ks+1))
			}
		}
	}

	return
}

//EndpointIndex get the index of matching endpoint if exists
func (cfg *ClientConfig) EndpointIndex(endpoint string) (ex bool, ix int) {
	eps := cfg.parsEndpoints()
	return helpers.InSlice(endpoint, eps)
}

//parsEndpoints get all available endpoints in client config file
func (cfg *ClientConfig) parsEndpoints() (eps []string) {
	for _, k := range (*cfg).Endpoints {
		eps = append(eps, k.Endpoint)
	}

	return
}
