package main

import (
	"fmt"
	"github.com/cuihairu/salon/internal/starter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

var (
	cfgFile string
	version = "0.1.0"

	zapLogger *zap.Logger
	db        *gorm.DB
)

var rootCmd = &cobra.Command{
	Use:     "salon",
	Short:   "salon is a mini-program server",
	Long:    `salon is a mini-program server`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		v := viper.New()
		v.SetConfigFile(cfgFile)
		app, err := starter.NewApp(v)
		if err != nil {
			return err
		}
		return app.Run()
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "starter", "c", "", "starter file (default is based on APP_ENV)")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
