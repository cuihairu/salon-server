package config

import "runtime"

type RedisConfig struct {
	Address  string `mapstructure:"address" json:"address" yaml:"address"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Db       int    `mapstructure:"db" json:"db" yaml:"db"`
	Secret   string `mapstructure:"secret" json:"secret" yaml:"secret"`
	NumConn  int    `mapstructure:"num_conn" json:"num_conn" yaml:"num_conn"`
}

func (c *Config) GetRedisConfig() (*RedisConfig, error) {
	redisConfig := &RedisConfig{}
	err := c.v.UnmarshalKey("redis", redisConfig)
	if err != nil {
		return nil, err
	}
	if redisConfig.NumConn <= 0 {
		redisConfig.NumConn = runtime.NumCPU()
	}

	return redisConfig, nil
}
