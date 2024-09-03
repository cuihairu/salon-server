package config

import "go.uber.org/zap"

type LoggingConfig struct {
	Level       string   `mapstructure:"level" json:"level" yaml:"level"`
	Path        string   `mapstructure:"path" json:"path" yaml:"path"`
	OutputPaths []string `mapstructure:"outputPaths" json:"outputPaths" yaml:"outputPaths"`
}

func (c *Config) GetZapConfig() (*zap.Config, error) {
	zapConfig := zap.NewProductionConfig()
	if c.v.IsSet("log") {
		loggingConfig := &LoggingConfig{}
		err := c.v.UnmarshalKey("log", loggingConfig)
		if err != nil {
			return nil, err
		}
		atomicLevel, err := zap.ParseAtomicLevel(loggingConfig.Level)
		if err != nil {
			return nil, err
		}
		zapConfig.Level = atomicLevel
		for _, outputPath := range loggingConfig.OutputPaths {
			zapConfig.OutputPaths = append(zapConfig.OutputPaths, outputPath)
		}
	} else {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, "logs/salon.log")
	}
	return &zapConfig, nil
}
