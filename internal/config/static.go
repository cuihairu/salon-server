package config

import (
	"fmt"
	"github.com/cuihairu/salon/internal/utils"
)

type StaticConfig struct {
	Domain      string `mapstructure:"domain" yaml:"domain" json:"domain"`
	UrlPath     string `mapstructure:"url_path" yaml:"url_path" json:"url_path"`
	EnableLocal bool   `mapstructure:"enable_local" yaml:"enable_local" json:"enable_local"`
	StaticPath  string `mapstructure:"static_path" yaml:"static_path" json:"static_path"`
}

func (c *Config) GetStaticConfig() (*StaticConfig, error) {
	config := &StaticConfig{}
	err := c.v.UnmarshalKey("static", config)
	if err != nil {
		return nil, err
	}
	if !utils.IsURL(config.Domain) {
		return nil, fmt.Errorf("invalid domain: %s", config.Domain)
	}
	if config.EnableLocal {
		if config.UrlPath == "" {
			return nil, fmt.Errorf("url path cannot be empty")
		}
		if config.StaticPath == "" {
			return nil, fmt.Errorf("static path cannot be empty")
		}
		err = utils.CreateDirIfNotExist(config.StaticPath)
		if err != nil {
			return nil, err
		}
	}
	return config, err
}
