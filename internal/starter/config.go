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
		v.AddConfigPath("./starter")
		v.SetConfigType("yaml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading starter file, %s", err)
		return err
	}
	fmt.Fprintf(os.Stdout, "Using starter file:%s \n", viper.ConfigFileUsed())
	return nil
}

func (c *Config) GetDbConfig() (*DatabaseConfig, error) {
	dbType := c.v.GetString("database.type")
	dsn := c.v.GetString("database.dsn")
	dbConf := &DatabaseConfig{
		DbType: dbType,
		DSN:    dsn,
	}
	return dbConf, nil
}

func (c *Config) GetEnv() string {
	return c.env
}

func (c *Config) GetZapConfig() (*zap.Config, error) {
	logConfig := viper.Sub("log")
	zapConfig := zap.NewProductionConfig()
	if logConfig != nil {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, "logs/salon.log")
	} else {
		if err := logConfig.Unmarshal(&zapConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading log starter file, %s", err)
			return nil, err
		}
	}
	return &zapConfig, nil
}
