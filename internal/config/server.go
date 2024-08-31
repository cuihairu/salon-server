package config

type ServerConfig struct {
	Address string `mapstructure:"address" yaml:"address" json:"address"`
}

func (c *Config) GetServerConfig() (*ServerConfig, error) {
	serverConfig := &ServerConfig{}
	err := c.v.UnmarshalKey("server", serverConfig)
	return serverConfig, err
}
