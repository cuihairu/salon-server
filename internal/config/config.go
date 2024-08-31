package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	v   *viper.Viper
	env string
}

func New(v *viper.Viper) (*Config, error) {
	c := &Config{
		v: v,
	}
	c.setEnv(os.Getenv("APP_ENV"))
	return c, c.load(v)
}

func (c *Config) load(v *viper.Viper) error {
	if len(v.ConfigFileUsed()) == 0 {
		v.SetConfigName(c.env)
		v.AddConfigPath("./config")
		v.SetConfigType("yaml")
	}
	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading starter file, %s", err)
		return err
	}
	fmt.Fprintf(os.Stdout, "Using config file:%s \n", v.ConfigFileUsed())
	return nil
}
