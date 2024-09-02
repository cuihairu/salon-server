package config

import (
	"time"
)

type JwtConfig struct {
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Expire    time.Duration
}

func (c *Config) GetJwtConfig() *JwtConfig {
	secretKey := c.v.GetString("jwt.secret_key")
	expire := c.v.GetDuration("jwt.expire")
	return &JwtConfig{
		SecretKey: secretKey,
		Expire:    expire,
	}
}
