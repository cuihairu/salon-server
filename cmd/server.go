package main

import (
	"fmt"
	"github.com/cuihairu/salon/internal"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	cfgFile   string
	version   = "0.1.0"
	DEV       = "dev"
	PROD      = "prod"
	zapLogger *zap.Logger
	db        *gorm.DB
)

func loadConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		env := os.Getenv("APP_ENV")
		switch env {
		case DEV:
		case PROD:
			env = PROD
		default:
			env = DEV
		}
		viper.SetConfigName(env)
		viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file, %s", err)
	}
	fmt.Fprintf(os.Stdout, "Using config file:%s \n", viper.ConfigFileUsed())
}

func initLogger() {
	var zapConfig zap.Config
	logConfig := viper.Sub("log")
	if logConfig != nil {
		zapConfig = zap.NewProductionConfig()
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, "logs/salon.log")
	} else {
		if err := logConfig.Unmarshal(&zapConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading log config file, %s", err)
			os.Exit(1)
		}
	}
	err := utils.CreateDirIfNotExist(zapConfig.OutputPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log dir, %s", err)
		os.Exit(1)
	}
	err = utils.CreateDirIfNotExist(zapConfig.ErrorOutputPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating error log dir, %s", err)
		os.Exit(1)
	}
	zapLogger, err = zapConfig.Build()
	zapLogger.Sugar()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger, %s", err)
	}
}

func initDatabase() error {
	sqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到的错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	dbType := viper.GetString("database.type")
	dsn := viper.GetString("database.dsn")
	var err error
	gormConfig := &gorm.Config{
		Logger: sqlLogger,
	}
	switch dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	}
	if err != nil {
		return err
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:     "salon",
	Short:   "salon is a mini-program server",
	Long:    `salon is a mini-program server`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()
		initLogger()
		initDatabase()
		app := internal.NewApp(db, zapLogger)
		app.Start()
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is based on APP_ENV)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
