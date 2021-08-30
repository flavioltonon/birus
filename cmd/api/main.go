package main

import (
	"os"

	"github.com/flavioltonon/birus/api"
	"github.com/flavioltonon/birus/internal/logger"

	"go.uber.org/zap"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		logger.Log().Error("failed to create new API server", zap.Error(err))
		os.Exit(1)
	}

	if err := server.Run(); err != nil {
		logger.Log().Error("failed to run server", zap.Error(err))
		os.Exit(1)
	}

	defer server.Stop()
}
