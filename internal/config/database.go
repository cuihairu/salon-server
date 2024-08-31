package config

type SQLLogConfig struct {
	Enabled                   bool   `mapstructure:"enabled" yaml:"enabled" json:"enabled"`
	Filename                  string `mapstructure:"filename" yaml:"filename" json:"filename"`
	LogLevel                  string `mapstructure:"log_level" yaml:"log_level" json:"log_level"`
	MaxSize                   int    `mapstructure:"max_size" yaml:"max_size" json:"max_size"`
	MaxBackups                int    `mapstructure:"max_backups" yaml:"max_backups" json:"max_backups"`
	MaxAge                    int    `mapstructure:"max_age" yaml:"max_age" json:"max_age"`
	Compress                  bool   `mapstructure:"compress" yaml:"compress" json:"compress"`
	SlowThreshold             int    `mapstructure:"slow_threshold" yaml:"slow_threshold" json:"slow_threshold"`
	IgnoreRecordNotFoundError bool   `mapstructure:"ignore_record_not_found_error" yaml:"ignore_record_not_found_error" json:"ignore_record_not_found_error"`
	Colorful                  bool   `mapstructure:"colorful" yaml:"colorful" json:"colorful"`
}

type DatabaseConfig struct {
	DbType      string       `mapstructure:"db_type" yaml:"db_type" json:"db_type"`
	DSN         string       `mapstructure:"dsn" yaml:"dsn" json:"dsn"`
	AutoMigrate bool         `mapstructure:"auto_migrate" yaml:"auto_migrate" json:"auto_migrate"`
	Log         SQLLogConfig `mapstructure:"log" yaml:"log" json:"log"`
}
