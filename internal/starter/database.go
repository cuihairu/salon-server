package starter

import (
	"github.com/cuihairu/salon/internal/data"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DatabaseConfig struct {
	DbType      string `yaml:"db_type" json:"db_type"`
	DSN         string `yaml:"dsn" json:"dsn"`
	AutoMigrate bool   `yaml:"auto_migrate" json:"auto_migrate"`
}

func NewDb(dbConf *DatabaseConfig) (*gorm.DB, error) {
	sqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到的错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	var err error
	var db *gorm.DB
	gormConfig := &gorm.Config{
		Logger: sqlLogger,
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
