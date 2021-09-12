package config

type ExtraConfig struct {
	name     string
	endpoint string
}

//SetName set extra_config name
func (ec *ExtraConfig) SetName(name string) {
	ec.name = name
}

//SetEndpoint set extra_config endpoint
func (ec *ExtraConfig) SetEndpoint(endpoint string) {
	ec.endpoint = endpoint
}

//Name get extra config name
func (ec *ExtraConfig) Name() string {
	return ec.name
}

// Endpoint get extra config endpoint
func (ec *ExtraConfig) Endpoint() string {
	return ec.endpoint
}
