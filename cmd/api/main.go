package main

import (
	"os"

	"birus/api"
	"birus/infrastructure/logger"

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
