package api

import (
	"context"
	"net/http"

	"birus/api/config"
	"birus/api/controller"
	"birus/application/service"
	"birus/infrastructure/engine"
	"birus/infrastructure/logger"
	"birus/infrastructure/repository"
	"birus/infrastructure/repository/mongodb"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Server extends *http.Server
type Server struct {
	core   *http.Server
	config *config.Config

	repository repository.Repository
	engine     engine.Engine
}

// NewServer returns a new Server
func NewServer() (*Server, error) {
	config, err := config.FromFile("config.yaml")
	if err != nil {
		return nil, err
	}

	if !config.Server.DevelopmentEnvironment {
		gin.SetMode(gin.ReleaseMode)
	}

	e, err := engine.NewGosseract(engine.GosseractOptions{
		TessdataPrefix: config.OCR.TessdataPrefix,
		Language:       config.OCR.Language,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initialize OCR engine")
	}

	r, err := mongodb.NewRepository(&mongodb.Options{
		DatabaseName: config.Database.Name,
		ClientOptions: []*options.ClientOptions{
			options.Client().ApplyURI(config.Database.URI),
		},
	})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create initialize repository")
	}

	ctx := context.Background()

	if err := r.Connect(ctx); err != nil {
		return nil, errors.WithMessage(err, "failed to establish connection with the repository")
	}

	ctrl := controller.New(&controller.Usecases{
		ImageClassification: service.NewImageClassificationService(
			service.NewOpticalCharacterRecognitionService(e),
			service.NewTextProcessingService(),
			r.ClassifierRepository,
		),
		OpticalCharacterRecognition: service.NewOpticalCharacterRecognitionService(e),
	})

	return &Server{
		core: &http.Server{
			Handler: ctrl.NewRouter(),
			Addr:    config.Server.Address,
		},
		config:     config,
		repository: r,
		engine:     e,
	}, nil
}

// Run starts a Server
func (s *Server) Run() error {
	logger.Log().Info("server listening and serving", zap.String("server_address", s.core.Addr))

	if err := s.core.ListenAndServe(); err != nil {
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

	if err := s.core.Shutdown(ctx); err != nil {
		return errors.WithMessage(err, "failed to shut server down")
	}

	return nil
}
