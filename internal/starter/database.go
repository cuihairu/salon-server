package starter

import (
	"github.com/cuihairu/salon/internal/data"
	"github.com/cuihairu/salon/internal/utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

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

func createLogger(sqlConfig *SQLLogConfig) (logger.Interface, error) {
	if sqlConfig == nil || !sqlConfig.Enabled {
		return nil, nil
	}
	var err error
	logWirter := log.New(os.Stdout, "\r\n", log.LstdFlags)
	filename := sqlConfig.Filename
	if len(filename) == 0 {
		filename, err = utils.CreateFileIfNotExistInCurPath("logs", "sql.log")
		if err != nil {
			return nil, err
		}
	}
	maxSize := sqlConfig.MaxSize
	if maxSize <= 0 {
		maxSize = 400 * 1024
	}
	maxBackups := sqlConfig.MaxBackups
	if maxBackups <= 0 {
		maxBackups = 3
	}
	maxAge := sqlConfig.MaxAge
	if maxAge <= 0 {
		maxAge = 30
	}
	logWirter.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		MaxSize:    maxSize,
		Compress:   sqlConfig.Compress,
	})
	slowThreshold := sqlConfig.SlowThreshold
	if slowThreshold < 1 {
		slowThreshold = 2
	}
	logLevel := logger.Info
	switch sqlConfig.LogLevel {
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	default:
	}
	sqlLogger := logger.New(
		logWirter, // io writer
		logger.Config{
			SlowThreshold:             time.Duration(slowThreshold) * time.Second, // 慢 SQL 阈值
			LogLevel:                  logLevel,                                   // Log level
			IgnoreRecordNotFoundError: sqlConfig.IgnoreRecordNotFoundError,        // 忽略记录未找到的错误
			Colorful:                  sqlConfig.Colorful,                         // 禁用彩色打印
		},
	)
	return sqlLogger, nil
}

func NewDb(dbConf *DatabaseConfig) (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	gormConfig := &gorm.Config{}
	sqlLogger, err := createLogger(&dbConf.Log)
	if err != nil {
		return nil, err
	}
	if sqlLogger != nil {
		gormConfig.Logger = sqlLogger
	}
	switch dbConf.DbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbConf.DSN), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbConf.DSN), &gorm.Config{})
	default:
		db, err = gorm.Open(postgres.Open(dbConf.DSN), gormConfig)
	}
	if err != nil {
		return nil, err
	}
	if dbConf.AutoMigrate {
		err = AutoMigrate(db)
	}
	return db, err
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&data.User{}, &data.Account{}, &data.Member{}, &data.Order{}, &data.Service{}, &data.Admin{})
}
