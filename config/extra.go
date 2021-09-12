package config

type Configuration struct {
	Url      string            `json:"url"`
	Method   string            `json:"method"`
	Timeout  int               `json:"timeout"`
	Header   map[string]string `json:"headers"`
	Message  string            `json:"message"`
	Body     bool              `json:"body"`
	Statuses []int             `json:"statuses"`
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

type ClientConfigs []ClientConfig
