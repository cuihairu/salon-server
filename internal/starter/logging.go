package starter

import (
	"fmt"
	"github.com/cuihairu/salon/internal/utils"
	"go.uber.org/zap"
	"os"
)

func NewZapLogger(zapConfig *zap.Config) (*zap.Logger, error) {
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
	return zapConfig.Build()
}
