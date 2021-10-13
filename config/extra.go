package config

type Configuration struct {
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   bool              `json:"body"`
}

type Steps struct {
	Alias        string        `json:"alias"`
	Timeout      int           `json:"timeout"`
	RetryMax     int           `json:"retry_max"`
	RetryWaitMin int           `json:"retry_wait_min"`
	RetryWaitMax int           `json:"retry_wait_max"`
	Statuses     []int         `json:"statuses"`
	Register     Configuration `json:"register"`
	Rollback     Configuration `json:"rollback"`
}

type Endpoints struct {
	Endpoint       string  `json:"endpoint"`
	Register       string  `json:"register"`
	Rollback       string  `json:"rollback"`
	RollbackFailed string  `json:"rollback_failed"`
	Steps          []Steps `json:"steps"`
}

type SagaClientConfig struct {
	LogLevel  string `json:"log_level"`
	Endpoints []Endpoints
}
