package config

type ServerConfig struct {
	Address string `mapstructure:"address" yaml:"address" json:"address"`
}
