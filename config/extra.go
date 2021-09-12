package config

type Configuration struct {
	Url              string   `json:"url"`
	Method           string   `json:"method"`
	Message          string   `json:"message"`
	RequireBody      bool     `json:"require_body"`
	AdditionalParams []string `json:"additional_params"`
	Statuses         []int    `json:"statuses"`
}

type Steps struct {
	Alias   string        `json:"alias"`
	Success Configuration `json:"success"`
	Failure Configuration `json:"failure"`
}

type ClientConfig struct {
	Endpoint string  `json:"endpoint"`
	Steps    []Steps `json:"steps"`
}

func (cc *ClientConfig) ClientConfig() *ClientConfig {
	return cc
}
