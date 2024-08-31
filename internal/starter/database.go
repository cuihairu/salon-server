package starter

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/model"
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

func createLogger(sqlConfig *config.SQLLogConfig) (logger.Interface, error) {
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

func NewDb(dbConf *config.DatabaseConfig) (*gorm.DB, error) {
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
	return db.AutoMigrate(&model.User{}, &model.Account{}, &model.Member{}, &model.Order{}, &model.Service{}, &model.Admin{})
}
