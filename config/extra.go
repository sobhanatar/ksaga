package config

type Configuration struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Timeout int               `json:"timeout"`
	Header  map[string]string `json:"header"`
	Body    bool              `json:"body"`
}

type Steps struct {
	Alias    string        `json:"alias"`
	Statuses []int         `json:"statuses"`
	Register Configuration `json:"register"`
	Rollback Configuration `json:"rollback"`
}

type ClientConfig struct {
	Endpoint       string  `json:"endpoint"`
	Register       string  `json:"register"`
	Rollback       string  `json:"rollback"`
	RollbackFailed string  `json:"rollback_failed"`
	Steps          []Steps `json:"steps"`
}

type ClientConfigs []ClientConfig
