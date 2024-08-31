package starter

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

const (
	DEV  = "dev"
	PROD = "prod"
	TEST = "test"
)

type Config struct {
	v   *viper.Viper
	env string
}

func New(v *viper.Viper) (*Config, error) {
	env := os.Getenv("APP_ENV")
	switch env {
	case DEV:
	case TEST:
		env = TEST
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

func (c *Config) IsDev() bool {
	return c.env == DEV
}

func (c *Config) IsTest() bool {
	return c.env == TEST
}

func (c *Config) IsProd() bool {
	return c.env == PROD
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
	dbConf := &DatabaseConfig{}
	err := c.v.UnmarshalKey("database", dbConf)
	return dbConf, err
}

func (c *Config) GetEnv() string {
	return c.env
}

func (c *Config) GetZapConfig() (*zap.Config, error) {
	zapConfig := zap.NewProductionConfig()
	if c.v.IsSet("log") {
		level := c.v.GetString("log.level")
		atomicLevel, err := zap.ParseAtomicLevel(level)
		if err != nil {
			return nil, err
		}
		zapConfig.Level = atomicLevel
		outputPaths := c.v.GetStringSlice("log.outputPaths")
		for _, outputPath := range outputPaths {
			zapConfig.OutputPaths = append(zapConfig.OutputPaths, outputPath)
		}
	} else {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, "logs/salon.log")
	}
	return &zapConfig, nil
}

func (c *Config) GetMiniappConfig() (*MiniappConfig, error) {
	config := &MiniappConfig{}
	err := c.v.UnmarshalKey("miniapp", config)
	return config, err
}

func (c *Config) GetServerConfig() (*ServerConfig, error) {
	serverConfig := &ServerConfig{}
	err := c.v.UnmarshalKey("server", serverConfig)
	return serverConfig, err
}
