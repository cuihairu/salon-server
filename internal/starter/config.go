package starter

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

const (
	DEV  = "dev"
	PROD = "prod"
)

type Config struct {
	v   *viper.Viper
	env string
}

func New(v *viper.Viper) (*Config, error) {
	env := os.Getenv("APP_ENV")
	switch env {
	case DEV:
	case PROD:
		env = PROD
	default:
		env = DEV
	}
	c := &Config{
		v:   v,
		env: env,
	}
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

func (c *Config) GetDbConfig() (*DatabaseConfig, error) {
	dbViperConf := c.v.Sub("database")
	if dbViperConf == nil {
		return nil, errors.New("database config not found")
	}
	dbConf := &DatabaseConfig{}
	err := dbViperConf.Unmarshal(dbConf)
	return dbConf, err
}

func (c *Config) GetEnv() string {
	return c.env
}

func (c *Config) GetZapConfig() (*zap.Config, error) {
	logConfig := viper.Sub("log")
	zapConfig := zap.NewProductionConfig()
	if logConfig == nil {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, "logs/salon.log")
	} else {
		if err := logConfig.Unmarshal(&zapConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading log starter file, %s", err)
			return nil, err
		}
	}
	return &zapConfig, nil
}
