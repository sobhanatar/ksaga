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
	Alias   string
	Success Configuration
	Failure Configuration
}

type ClientConfig struct {
	Endpoint string `json:"endpoint"`
	Steps    []Steps
}
