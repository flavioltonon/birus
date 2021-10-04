package logger

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func start() {
	new := zap.NewProduction

	if viper.GetBool("DEVELOPMENT_ENVIRONMENT") {
		new = zap.NewDevelopment
	}

	logger, err := new()
	if err != nil {
		panic(errors.WithMessage(err, "failed to initialize logger"))
	}

	zap.ReplaceGlobals(logger)
}

var _once sync.Once

// Log exposes zap global logger. This function makes the init function in this package to be always called once,
// making it unlikely to end up having any issues by forgetting to import this package silently.
func Log() *zap.Logger {
	_once.Do(start)
	return zap.L()
}
