package api

import (
	"context"
	"net/http"

	"github.com/flavioltonon/birus/api/controller"
	"github.com/flavioltonon/birus/application/service"
	"github.com/flavioltonon/birus/infrastructure/engine"
	"github.com/flavioltonon/birus/infrastructure/repository"
	"github.com/flavioltonon/birus/internal/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Server extends *http.Server
type Server struct {
	config *http.Server

	repository repository.Repository
	engine     engine.Engine
}

// NewServer returns a new Server
func NewServer() (*Server, error) {
	e, err := engine.NewGosseract(engine.GosseractOptions{
		TessdataPrefix: viper.GetString("OCR_ENGINE_TESSDATA_PREFIX"),
		Language:       viper.GetString("OCR_ENGINE_LANGUAGE"),
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initialize OCR engine")
	}

	r, err := repository.NewMemoryRepository()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create initialize repository")
	}

	modelsService := service.NewModelsService(r.Models, e)

	ctrl := controller.New(&controller.Usecases{
		Models:              modelsService,
		ImageClassification: service.NewImageClassificationService(modelsService, e),
	})

	return &Server{
		config: &http.Server{
			Handler: ctrl.NewRouter(),
			Addr:    viper.GetString("SERVER_ADDRESS"),
		},
	}, nil
}

// Run starts a Server
func (s *Server) Run() error {
	logger.Log().Info("server listening and serving", zap.String("server_address", s.config.Addr))

	if err := s.config.ListenAndServe(); err != nil {
		return errors.WithMessage(err, "failed to run server")
	}

	return nil
}

// Stop shuts the Server down
func (s *Server) Stop() error {
	ctx := context.Background()

	if err := s.repository.Disconnect(ctx); err != nil {
		return errors.WithMessage(err, "failed to disconnect from repository")
	}

	if err := engine.Stop(s.engine); err != nil {
		return errors.WithMessage(err, "failed to stop OCR engine")
	}

	if err := s.config.Shutdown(ctx); err != nil {
		return errors.WithMessage(err, "failed to shut server down")
	}

	return nil
}
