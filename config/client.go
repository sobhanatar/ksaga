package config

type ExtraConfig struct {
	name     string
	endpoint string
}

func (ec *ExtraConfig) SetName(name string) {
	ec.name = name
}

func (ec *ExtraConfig) SetEndpoint(endpoint string) {
	ec.endpoint = endpoint
}

func (ec *ExtraConfig) Name() string {
	return ec.name
}

func (ec *ExtraConfig) Endpoint() string {
	return ec.endpoint
}
